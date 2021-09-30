package converter

import (
	"github.com/gauravcoco/bms/dbModels"
	"github.com/gauravcoco/bms/models"
	"github.com/gauravcoco/bms/providers"
	"github.com/gauravcoco/bms/utils"
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
		Dimension:         details.Dimension,
		Language:          details.Language,
		ReleaseDate:       details.ReleaseDate,
		DurationInMinutes: details.DurationInMinutes,
		Cast:              c.ToMovieCasts(details.Cast),
	}

	for i := range details.Banners {
		jsonMovieDetails.Banners = append(jsonMovieDetails.Banners, utils.GetImageURL(details.Banners[i].Bucket, details.Banners[i].Path))
	}
	panic("implement me")
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

func (c *converter) ToMovieShowDetails(_ dbModels.MovieShowDetails) models.ShowDetails {
	//jsonShowDetails := models.ShowDetails{
	//	Theater: models.Theater{
	//		ID:   show.TheaterID,
	//		Name: show.TheaterName,
	//	},
	//	Movie: models.Movie{
	//		ID:   show.MovieID,
	//		Name: show.MovieName,
	//	},
	//}
	//
	//groupedShows := make(map[time.Time][]models.ShowTimings)
	//
	//for i := range
	panic("panic")
}
