package dbHelpProvider

import (
	"context"
	"errors"
	"github.com/gauravcoco/bms/db"
	"github.com/gauravcoco/bms/dbModels"
	"github.com/gauravcoco/bms/models"
	"github.com/gauravcoco/bms/providers"
	"github.com/gauravcoco/bms/sql"
	"github.com/gauravcoco/bms/utils"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"golang.org/x/sync/errgroup"
)

type DBHelper struct {
	DB          *sqlx.DB
	KeyProvider providers.KeyProvider
}

func NewDBHelper(db *sqlx.DB, kp providers.KeyProvider) providers.DBHelpProvider {
	return &DBHelper{
		DB:          db,
		KeyProvider: kp,
	}
}

func (helper *DBHelper) CreateUser(_ context.Context, user *dbModels.User) (int, error) {
	if user == nil {
		return 0, nil
	}
	var userID int

	args := []interface{}{
		user.Name,
		user.Email,
		user.Password,
	}
	err := helper.DB.Get(&userID, sql.CreateNewUserSQL, args...)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (helper *DBHelper) GetUserByID(_ context.Context, id int) (*dbModels.User, error) {
	var user dbModels.User
	err := helper.DB.Get(&user, sql.GetUserInfoByIdSQL, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (helper *DBHelper) GetAllMovies(ctx context.Context, req models.MovieSearchRequest) ([]dbModels.Movie, error) {
	movies := make([]dbModels.Movie, 0)
	SQL := `SELECT 
				movies.id,
				movies.title AS name,
				movies.release_date,
				movies.duration_in_minutes,
				movie_dimensions.dimensions,
				movie_languages.languages
			FROM movies
			JOIN LATERAL (
					SELECT count(mts.id) > 0 AS is_available
					FROM movie_theater_schedule mts 
					JOIN theater ON theater.id = mts.theater_id
					WHERE mts.movie_id = movies.id AND 
						  mts.status <> 'closed' AND
						  mts.archived_at IS NULL AND 
						  (theater.city_id = ? OR ?)
			) movie_available ON TRUE
			JOIN LATERAL (
					SELECT ARRAY_AGG(md.dimension) AS dimensions
					FROM movie_dimension md 
					WHERE md.movie_id = movies.id
			) movie_dimensions ON TRUE
			JOIN LATERAL (
					SELECT ARRAY_AGG(ml.language) AS languages
					FROM movie_language ml 
					WHERE ml.movie_id = movies.id
			) movie_languages ON TRUE
			WHERE 
				release_date >= now() AND 
				archived_at IS NULL AND 
				movie_available.is_available
				`

	args := make([]interface{}, 0)

	args = append(args, req.CityID, !req.CityID.Valid)

	if req.Name != "" {
		SQL += ` AND movies.title ILIKE '%%' || ? || '%%'`
		args = append(args, req.Name)
	}

	if len(req.Dimensions) > 0 {
		SQL += " AND movie_dimensions.dimensions&&(?)"
		args = append(args, req.Dimensions)
	}

	if len(req.Languages) > 0 {
		SQL += " AND movie_languages.languages&&(?)"
		args = append(args, req.Languages)
	}

	SQL = db.QuestionToDollar(SQL)

	err := helper.DB.Select(&movies, SQL, args...)
	if err != nil {
		return nil, err
	}

	movieIDs := make([]int64, 0)

	for _, movie := range movies {
		movieIDs = append(movieIDs, int64(movie.ID))
	}

	movieImages, err := helper.GetMovieBannerImages(ctx, movieIDs)
	if err != nil {
		return nil, err
	}

	for i := range movies {
		movies[i].Banners = movieImages[movies[i].ID]
	}

	return movies, nil
}

func (helper *DBHelper) GetMovieBannerImages(_ context.Context, ids pq.Int64Array) (map[int][]dbModels.Image, error) {
	var movieImages []struct {
		MovieID int `db:"movie_id"`
		dbModels.Image
	}
	SQL := `SELECT
				movie_banners.movie_id,
				images.id,
				images.bucket,
				images.path
			FROM movie_banners
			JOIN images ON images.id = movie_banners.banner_id
			WHERE 
				movie_banners.movie_id = ANY($1)
				`

	err := helper.DB.Select(&movieImages, SQL, ids)
	if err != nil {
		return nil, err
	}

	movieImageMap := make(map[int][]dbModels.Image)

	for _, mi := range movieImages {
		movieImageMap[mi.MovieID] = append(movieImageMap[mi.MovieID], mi.Image)
	}
	return movieImageMap, nil
}

func (helper *DBHelper) GetMovieDetails(ctx context.Context, id int) (*dbModels.MovieDetails, error) {
	var movieDetails dbModels.MovieDetails

	err := helper.DB.Get(&movieDetails, sql.MovieDetailsSQL, id)
	if err != nil {
		return nil, err
	}

	movieIDs := make([]int64, 1)
	movieIDs[0] = int64(movieDetails.ID)

	egp := new(errgroup.Group)

	egp.Go(func() error {
		movieImages, err := helper.GetMovieBannerImages(ctx, movieIDs)
		if err != nil {
			return err
		}
		movieDetails.Banners = movieImages[movieDetails.ID]
		return nil
	})

	egp.Go(func() error {
		castImages, err := helper.GetMovieCastImages(ctx, movieDetails.ID)
		if err != nil {
			return err
		}

		for i := range movieDetails.Cast {
			movieDetails.Cast[i].Image = castImages[movieDetails.Cast[i].ID]
		}
		return nil
	})

	err = egp.Wait()
	if err != nil {
		return nil, err
	}

	return &movieDetails, nil
}

func (helper *DBHelper) GetMovieShowDetails(_ context.Context, id int) (*dbModels.MovieShowDetails, error) {
	var movieShowsDetails dbModels.MovieShowDetails

	egp := new(errgroup.Group)

	egp.Go(func() error {
		err := helper.DB.Get(&movieShowsDetails, sql.MovieDetailsSQL, id)
		if err != nil {
			return err
		}
		return nil
	})

	egp.Go(func() error {
		err := helper.DB.Select(&movieShowsDetails.ShowDetails, sql.MovieShowDetailsSQL, id)
		if err != nil {
			return err
		}
		return nil
	})

	err := egp.Wait()
	if err != nil {
		return nil, err
	}

	return &movieShowsDetails, nil

}

func (helper *DBHelper) BookMovieTicket(ctx context.Context, request models.BookingRequest) (dbModels.Booking, error) {
	var bookingID int
	err := db.WithTransaction(helper.DB, func(tx *sqlx.Tx) error {

		SQL := `SELECT 
					theater_id,
					movie_id,
					hall_id
				FROM show_timings
				WHERE id = $1 AND 
					  archived_at IS NULL AND
					  show_start_time > now()
				 `

		var movieID, theaterID, hallID int
		err := tx.QueryRowx(SQL, request.ShowID).Scan(&theaterID, &movieID, &hallID)
		if err != nil {
			return err
		}

		SQL = `INSERT INTO bookings (show_id, hall_id, theater_id, movie_id, user_id, booking_pretty_id, status)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
				RETURNING id`

		args := []interface{}{
			request.ShowID,
			hallID,
			theaterID,
			movieID,
			request.UserID,
			helper.KeyProvider.GenerateUniqueKey(),
			dbModels.BookingStatusUnConfirmed,
		}

		err = tx.Get(&bookingID, SQL, args...)
		if err != nil {
			return err
		}

		SQL = `SELECT count(bookings_seats.id) > 0 AS is_any_seat_already_occupied
				FROM bookings_seats
				JOIN movie_hall_seating mhs on bookings_seats.seat_id = mhs.id
				JOIN bookings b on bookings_seats.booking_id = b.id
				WHERE b.show_id = $1 AND 
					  b.status = $2 AND
					  mhs.id = ANY ($3) AND
					  mhs.status = $4
				 `

		var isAnySeatOccupied bool

		args = []interface{}{
			request.ShowID,
			dbModels.BookingStatusConfirmed,
			pq.Int64Array(utils.ToInt64Slice(request.SeatIDs)),
			dbModels.SeatStatusFunctional,
		}

		err = tx.Get(&isAnySeatOccupied, SQL, args...)
		if err != nil {
			return err
		}

		if isAnySeatOccupied {
			return errors.New("someone else booked this seat already. Please tr again")
		}

		SQL = `INSERT INTO bookings_seats (booking_id, seat_id)
				SELECT 
					$1, id 
				FROM 
					movie_hall_seating WHERE id = ANY($2)
			`

		args = []interface{}{
			bookingID,
			pq.Int64Array(utils.ToInt64Slice(request.SeatIDs)),
		}

		_, err = tx.Exec(SQL, args...)
		if err != nil {
			return err
		}

		SQL = `UPDATE bookings 
			   SET 
				  status = $1
			   WHERE 
				  id = $2
			`

		args = []interface{}{
			dbModels.BookingStatusConfirmed,
			bookingID,
		}
		_, err = tx.Exec(SQL, args...)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return dbModels.Booking{}, err
	}

	booking, err := helper.GetBookingByID(ctx, bookingID)
	if err != nil {
		return dbModels.Booking{}, err
	}

	return booking, nil
}

func (helper *DBHelper) GetMovieCastImages(_ context.Context, id int) (map[int]dbModels.Image, error) {
	var castImages []struct {
		CastID int `db:"cast_id"`
		dbModels.Image
	}
	SQL := `SELECT
				movie_cast.id AS cast_id,
				images.id,
				images.bucket,
				images.path
			FROM movie_cast
			JOIN images ON images.id = movie_cast.image_id
			WHERE 
				movie_cast.movie_id = $1
				`

	err := helper.DB.Select(&castImages, SQL, id)
	if err != nil {
		return nil, err
	}

	castImageMap := make(map[int]dbModels.Image)

	for _, ci := range castImages {
		castImageMap[ci.CastID] = ci.Image
	}
	return castImageMap, nil
}

func (helper *DBHelper) GetShowSeats(_ context.Context, showID int) (dbModels.ShowSeats, error) {
	var seatDetails dbModels.ShowSeats
	SQL := `SELECT
				movie_halls.id AS hall_id,
				movie_halls.total_rows AS total_rows,
				movie_halls.total_columns AS total_columns,
				movie_halls.total_seats AS total_seats
			FROM show_timings
			JOIN movie_halls ON movie_halls.id = show_timings.hall_id
			WHERE
				show_timings.id = $1 AND
				show_timings.archived_at IS NULL;`

	err := helper.DB.Get(&seatDetails, SQL, showID)
	if err != nil {
		return seatDetails, err
	}

	SQL = `SELECT
				movie_hall_seating.id,
				movie_hall_seating.seat_code,
				movie_hall_seating.status,
				movie_hall_seating.row_number AS row_pos,
				movie_hall_seating.column_number column_pos,
				movie_hall_seating.is_recliner,
				(b.id IS NOT NULL) AS is_booked
			FROM movie_hall_seating
			JOIN show_timings st ON st.hall_id = movie_hall_seating.hall_id
			LEFT JOIN bookings_seats bs on movie_hall_seating.id = bs.seat_id
			LEFT JOIN bookings b on bs.booking_id = b.id AND b.status = $2 AND b.show_id = $3
			WHERE movie_hall_seating.hall_id = $1 AND
				  st.id = $3;`

	args := []interface{}{
		seatDetails.HallID,
		dbModels.BookingStatusConfirmed,
		showID,
	}

	err = helper.DB.Select(&seatDetails.SeatDetails, SQL, args...)
	if err != nil {
		return seatDetails, err
	}

	return seatDetails, err
}

func (helper *DBHelper) GetAllUserBookings(ctx context.Context, userID int) ([]dbModels.Booking, error) {
	SQL := `SELECT bookings.id                AS id,
				   bookings.booking_pretty_id AS booking_id,
				   bookings.status            AS booking_status,
				   m.id                       AS movie_id,
				   m.title                    AS movie_name,
				   movie_halls.id             AS hall_id,
				   movie_halls.name           AS hall_name,
				   st.id                      AS show_id,
				   st.show_start_time         AS show_start_time,
				   st.show_end_time           AS show_end_time,
				   t.id                       AS theater_id,
				   t.name                     AS theater_name
			FROM bookings
					 JOIN movies m on bookings.movie_id = m.id
					 JOIN movie_halls ON bookings.hall_id = movie_halls.id
					 JOIN show_timings st on bookings.show_id = st.id
					 JOIN theater t on movie_halls.theater_id = t.id
			WHERE user_id = $1;
			`

	bookings := make([]dbModels.Booking, 0)

	err := helper.DB.Select(&bookings, SQL, userID)
	if err != nil {
		return nil, err
	}

	if len(bookings) <= 0 {
		return bookings, nil
	}

	var bookingIDs = make(pq.Int64Array, len(bookings))

	for i := range bookings {
		bookingIDs[i] = int64(bookings[i].ID)
	}

	seatsMap, err := helper.GetBookingsSeats(ctx, bookingIDs)
	if err != nil {
		return nil, err
	}

	for i := range bookings {
		bookings[i].Seats = seatsMap[bookings[i].ID]
	}

	return bookings, nil

}

func (helper *DBHelper) GetBookingsSeats(_ context.Context, bookingIDs pq.Int64Array) (map[int][]dbModels.BookingSeats, error) {
	SQL := `SELECT movie_hall_seating.id,
				   movie_hall_seating.seat_code,
				   movie_hall_seating.is_recliner,
				   bookings_seats.booking_id
			FROM bookings_seats
					 JOIN movie_hall_seating ON bookings_seats.seat_id = movie_hall_seating.id
					 WHERE bookings_seats.booking_id = ANY($1);
			`

	seats := make([]dbModels.BookingSeats, 0)

	err := helper.DB.Select(&seats, SQL, bookingIDs)
	if err != nil {
		return nil, err
	}

	seatsMap := make(map[int][]dbModels.BookingSeats, 0)

	for i := range seats {
		seatsMap[seats[i].BookingID] = append(seatsMap[seats[i].BookingID], seats[i])
	}

	return seatsMap, nil
}

func (helper *DBHelper) GetBookingByID(ctx context.Context, bookingID int) (dbModels.Booking, error) {
	SQL := `SELECT bookings.id                AS id,
				   bookings.booking_pretty_id AS booking_id,
				   bookings.status            AS booking_status,
				   m.id                       AS movie_id,
				   m.title                    AS movie_name,
				   movie_halls.id             AS hall_id,
				   movie_halls.name           AS hall_name,
				   st.id                      AS show_id,
				   st.show_start_time         AS show_start_time,
				   st.show_end_time           AS show_end_time,
				   t.id                       AS theater_id,
				   t.name                     AS theater_name
			FROM bookings
					 JOIN movies m on bookings.movie_id = m.id
					 JOIN movie_halls ON bookings.hall_id = movie_halls.id
					 JOIN show_timings st on bookings.show_id = st.id
					 JOIN theater t on t.id = st.theater_id
			WHERE bookings.id = $1;
			`

	var booking dbModels.Booking
	err := helper.DB.Get(&booking, SQL, bookingID)
	if err != nil {
		return dbModels.Booking{}, err
	}

	seatsMap, err := helper.GetBookingsSeats(ctx, []int64{int64(booking.ID)})
	if err != nil {
		return dbModels.Booking{}, err
	}

	booking.Seats = seatsMap[booking.ID]

	return booking, nil

}
