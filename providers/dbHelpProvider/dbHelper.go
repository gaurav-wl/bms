package dbHelpProvider

import (
	"context"
	"github.com/gauravcoco/bms/db"
	"github.com/gauravcoco/bms/dbModels"
	"github.com/gauravcoco/bms/models"
	"github.com/gauravcoco/bms/providers"
	"github.com/gauravcoco/bms/sql"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"golang.org/x/sync/errgroup"
)

type DBHelper struct {
	DB *sqlx.DB
}

func NewDBHelper(db *sqlx.DB) providers.DBHelpProvider {
	return &DBHelper{
		DB: db,
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
				movies.name,
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
					FROM movie_languages ml 
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
		SQL = " AND movies.name ilike %?%"
		args = append(args, req.Name)
	}

	if len(req.Dimensions) > 0 {
		SQL += " AND movie_dimensions.dimensions = ANY(?)"
		args = append(args, req.Dimensions)
	}

	if len(req.Languages) > 0 {
		SQL += " AND movie_languages.languages = ANY(?)"
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

func (helper *DBHelper) GetMovieBannerImages(ctx context.Context, ids pq.Int64Array) (map[int][]dbModels.Image, error) {
	var movieImages []struct {
		MovieID int `db:"movie_id"`
		dbModels.Image
	}
	SQL := `SELECT
				movie_banners.movie_id
				images.id,
				images.bucket,
				images.path,
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

func (helper *DBHelper) GetMovieShowDetails(ctx context.Context, id int) (*dbModels.MovieShowDetails, error) {
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

func (helper *DBHelper) BookMovieTicket(ctx context.Context, request models.BookingRequest) (models.BookingDetails, error) {
	panic("not implemented")
}

func (helper *DBHelper) GetMovieCastImages(ctx context.Context, id int) (map[int]dbModels.Image, error) {
	var castImages []struct {
		CastID int `db:"cast_id"`
		dbModels.Image
	}
	SQL := `SELECT
				movie_cast.id AS cast_id
				images.id,
				images.bucket,
				images.path,
			FROM movie_cast
			JOIN images ON images.id = movie_cast.banner_id
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

func (helper *DBHelper) GetShowSeats(ctx context.Context, showID int) ([]dbModels.ShowSeats, error) {
	//seatDetails := make([]dbModels.ShowSeats, 0)
	//SQL := `SELECT
	//			movie_halls.id AS hall_id
	//			movie_halls.total_rows AS total_rows
	//			movie_halls.total_columns AS total_columns
	//			movie_halls.total_seats AS total_seats
	//		FROM show_timings
	//		JOIN movie_halls ON movie_halls.id = show_timings.hall_id
	//		WHERE
	//			show_timings.id = $1 AND
	//			show_timings.archived_at IS NULL`
	//
	//err := helper.DB.Get(&seatDetails, sql.MovieDetailsSQL, id)
	//if err != nil {
	//	return nil, err
	//}
	//
	//return nil
	return nil, nil
}
