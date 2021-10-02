package converter

import (
	"github.com/gauravcoco/bms/dbModels"
	"github.com/gauravcoco/bms/models"
	"github.com/gauravcoco/bms/providers"
	"github.com/gauravcoco/bms/utils"
	"time"
)

type converter struct{}

func NewConverter() providers.Converter {
	return &converter{}
}

func (c *converter) ToUser(user *dbModels.User) *models.User {
	if user == nil {
		return nil
	}
	webUser := &models.User{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}

	return webUser
}

func (c *converter) ToMovies(movies []dbModels.Movie) []models.Movie {
	jsonMovies := make([]models.Movie, len(movies))
	for i := range movies {
		jsonMovies[i] = c.ToMovie(movies[i])
	}
	return jsonMovies
}

func (c *converter) ToMovie(movie dbModels.Movie) models.Movie {
	jsonMovie := models.Movie{
		ID:                movie.ID,
		Name:              movie.Name,
		Languages:         movie.Languages,
		Dimensions:        movie.Dimensions,
		ReleaseDate:       movie.ReleaseDate,
		DurationInMinutes: movie.DurationInMinutes,
	}
	for i := range movie.Banners {
		jsonMovie.Banners = append(jsonMovie.Banners, utils.GetImageURL(movie.Banners[i].Bucket, movie.Banners[i].Path))
	}
	return jsonMovie
}

func (c *converter) ToMovieDetails(details dbModels.MovieDetails) models.MovieDetails {
	jsonMovieDetails := models.MovieDetails{
		ID:                details.ID,
		Name:              details.Name,
		Dimensions:        details.Dimensions,
		Language:          details.Language,
		ReleaseDate:       details.ReleaseDate,
		DurationInMinutes: details.DurationInMinutes,
		Cast:              c.ToMovieCasts(details.Cast),
	}

	for i := range details.Banners {
		jsonMovieDetails.Banners = append(jsonMovieDetails.Banners, utils.GetImageURL(details.Banners[i].Bucket, details.Banners[i].Path))
	}
	return jsonMovieDetails
}

func (c *converter) ToMovieCasts(casts []dbModels.Cast) []models.Cast {
	jsonCast := make([]models.Cast, len(casts))
	for i := range casts {
		jsonCast[i] = c.ToMovieCast(casts[i])
	}
	return jsonCast
}

func (c *converter) ToMovieCast(cast dbModels.Cast) models.Cast {
	return models.Cast{
		Name:  cast.Name,
		Image: utils.GetImageURL(cast.Image.Bucket, cast.Image.Path),
	}
}

func (c *converter) ToMovieShowDetails(show dbModels.MovieShowDetails) models.MovieShowDetails {
	jsonShowDetails := models.MovieShowDetails{
		ID:         show.ID,
		Name:       show.Name,
		Dimensions: show.Dimensions,
		Language:   show.Languages,
		City: models.City{
			ID:   show.City.ID,
			Name: show.City.Name,
		},
		DurationInMinutes: show.DurationInMinutes,
	}

	groupedShows := make(map[int][]dbModels.ShowDetails)

	for i := range show.ShowDetails {
		groupedShows[show.ShowDetails[i].TheaterID] = append(groupedShows[show.ShowDetails[i].TheaterID], show.ShowDetails[i])
	}

	for tID, timings := range groupedShows {
		var webShowDetails models.ShowDetails

		webShowDetails.Theater.ID = tID
		webShowDetails.Theater.Name = timings[0].TheaterName

		dayByTimings := make(map[time.Time][]dbModels.ShowDetails)

		for i := range timings {
			dayByTimings[utils.ToDate(timings[i].ShowStartTime)] = append(dayByTimings[utils.ToDate(timings[i].ShowStartTime)], timings[i])
		}

		var webShowTimings models.ShowTimings

		for date, timingsForADay := range dayByTimings {
			webShowTimings.Date = date

			for j := range timingsForADay {
				webShowTimings.Shows = append(webShowTimings.Shows, models.Shows{
					ID:        timingsForADay[j].ShowID,
					StartTime: timingsForADay[j].ShowStartTime,
					EndTime:   timingsForADay[j].ShowEndTime,
					Dimension: timingsForADay[j].Dimension,
					Language:  timingsForADay[j].Language,
				})
			}
			webShowDetails.Shows = append(webShowDetails.Shows, webShowTimings)
		}
		jsonShowDetails.ShowDetails = append(jsonShowDetails.ShowDetails, webShowDetails)
	}

	return jsonShowDetails
}

func (c *converter) ToBooking(booking dbModels.Booking) models.BookingDetails {
	webBooking := models.BookingDetails{
		Movie: models.BookedMovieDetail{
			ID:   booking.MovieID,
			Name: booking.MovieName,
		},
		Status: c.ToBookingStatus(booking.Status),
		Theater: models.Theater{
			ID:   booking.TheaterID,
			Name: booking.TheaterName,
		},
		Hall:       booking.HallName,
		BookingUID: booking.BookingID,
		ShowTiming: models.ShowTiming{
			StartTime: booking.ShowStartTime,
			EndTime:   booking.ShowEndTime,
		},
	}

	for i := range booking.Seats {
		webBooking.SeatCodes = append(webBooking.SeatCodes, models.SeatCodes{
			Name:       booking.Seats[i].SeatsCode,
			IsRecliner: booking.Seats[i].IsRecliner,
		})
	}

	return webBooking
}

func (c *converter) ToBookings(bookings []dbModels.Booking) []models.BookingDetails {
	if len(bookings) <= 0 {
		return []models.BookingDetails{}
	}

	webBookings := make([]models.BookingDetails, len(bookings))
	for i := range bookings {
		webBookings[i] = c.ToBooking(bookings[i])
	}

	return webBookings
}

func (c *converter) ToSeatsDetails(seatsDetails dbModels.ShowSeats) models.ShowSeats {
	webSeatDetails := models.ShowSeats{
		HallID:      seatsDetails.HallID,
		TotalRows:   seatsDetails.TotalRows,
		TotalColumn: seatsDetails.TotalColumn,
		TotalSeats:  seatsDetails.TotalSeats,
	}

	webSeatDetails.SeatDetails = make([]models.MovieShowSeatsDetails, len(seatsDetails.SeatDetails))
	for i := range seatsDetails.SeatDetails {
		webSeatDetails.SeatDetails[i] = models.MovieShowSeatsDetails{
			ID:         seatsDetails.SeatDetails[i].ID,
			SeatCode:   seatsDetails.SeatDetails[i].SeatCode,
			RowPos:     seatsDetails.SeatDetails[i].RowPos,
			ColumnPos:  seatsDetails.SeatDetails[i].ColumnPos,
			IsRecliner: seatsDetails.SeatDetails[i].IsRecliner,
			IsBooked:   seatsDetails.SeatDetails[i].IsBooked,
		}
	}

	return webSeatDetails
}
