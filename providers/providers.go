package providers

import (
	"context"
	"github.com/bms/dbModels"
	"github.com/bms/models"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type ConfigProvider interface {
	Read() error
	GetString(key string) string
	GetPSQLConnectionString() string
	GetPSQLMaxConnection() int
	GetPSQLMaxIdleConnection() int
	GetInt(key string) int
	GetAny(key string) interface{}
	GetJWTKey() string
	GetServerPort() string
}

type DBProvider interface {
	Ping() error
	PSQLProvider
}

type PSQLProvider interface {
	DB() *sqlx.DB
}

type MiddlewareProvider interface {
	Default() chi.Middlewares
	AUTH() chi.Middlewares
}

type KeyProvider interface {
	GenerateUniqueKey() string
}

type DBHelpProvider interface {
	CreateUser(ctx context.Context, user *dbModels.User) (int, error)
	GetUserByID(ctx context.Context, id int) (*dbModels.User, error)
	GetAllMovies(ctx context.Context, req models.MovieSearchRequest) ([]dbModels.Movie, error)
	GetMovieDetails(ctx context.Context, id int) (*dbModels.MovieDetails, error)
	GetMovieShowDetails(ctx context.Context, id int) (*dbModels.MovieShowDetails, error)
	GetShowSeats(ctx context.Context, showID int) ([]dbModels.ShowSeats, error)
	BookMovieTicket(ctx context.Context, request models.BookingRequest) (models.BookingDetails, error)
	GetMovieBannerImages(ctx context.Context, ids pq.Int64Array) (map[int][]dbModels.Image, error)
	GetMovieCastImages(ctx context.Context, id int) (map[int]dbModels.Image, error)
}

type Converter interface {
	ToUser(user *dbModels.User) *models.User
	ToMovies(movies []dbModels.Movie) []models.Movie
	ToMovie(movie dbModels.Movie) models.Movie
	ToMovieDetails(movieDetails dbModels.MovieDetails) models.MovieDetails
	ToMovieCasts(casts []dbModels.Cast) []models.Cast
	ToMovieCast(cast dbModels.Cast) models.Cast
	ToMovieShowDetails(show dbModels.MovieShowDetails) models.ShowDetails
}

type MigrationProvider interface {
	Up()
}
