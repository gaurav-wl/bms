package dbModels

import (
	"github.com/lib/pq"
	"time"
)

type Movie struct {
	ID                int            `db:"id"`
	Name              string         `db:"name"`
	Languages         pq.StringArray `db:"languages"`
	ReleaseDate       time.Time      `db:"releaseDate"`
	DurationInMinutes int            `db:"duration_in_minutes"`
	Banners           []Image        `db:"banners"`
}

type MovieShowDetails struct {
	ID                int            `db:"id"`
	Name              string         `db:"name"`
	Dimension         pq.StringArray `db:"dimension"`
	Languages         pq.StringArray `db:"languages"`
	City              City           `db:"city"`
	DurationInMinutes int            `db:"duration_in_minutes"`
	ShowDetails       []ShowDetails  `db:"-"`
}

type MovieDetails struct {
	ID                int       `db:"id"`
	Name              string    `db:"name"`
	Dimension         []string  `db:"dimension"`
	Language          []string  `db:"languages"`
	ReleaseDate       time.Time `db:"releaseDate"`
	DurationInMinutes int       `db:"duration_in_minutes"`
	Banners           []Image   `db:"banners"`
	Cast              []Cast    `db:"cast"`
}

type Cast struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Image Image  `db:"image"`
}

type Image struct {
	ID     int    `db:"id"`
	Bucket string `db:"bucket"`
	Path   string `db:"path"`
}

type City struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type Point struct {
	Lat  float64 `db:"lat"`
	Long float64 `db:"long"`
}

type Theater struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	City     string `db:"city"`
	Location Point  `db:"location"`
	Address  string `db:"string"`
}

type ShowDetails struct {
	MovieID       int       `db:"movie_id"`
	MovieName     string    `db:"movie_name"`
	ShowID        int       `db:"show_id"`
	ShowStartTime time.Time `db:"show_start_time"`
	ShowEndTime   time.Time `db:"show_end_time"`
	TheaterID     int       `db:"theater_id"`
	TheaterName   string    `db:"theater_name"`
	Dimension     string    `db:"dimension"`
	Language      string    `db:"language"`
}

type ShowSeats struct {
	HallID      int `db:"hall_id"`
	TotalRows   int `db:"total_rows"`
	TotalColumn int `db:"total_column"`
	TotalSeats  int `db:"total_seats"`
	SeatDetails MovieShowSeatsDetails
}

type MovieShowSeatsDetails struct {
	ID         int    `db:"id"`
	SeatCode   string `db:"seat_code"`
	RowPos     int    `db:"row_pos"`
	ColumnPos  int    `db:"column_pos"`
	IsRecliner bool   `db:"is_recliner"`
}
