package dbModels

import (
	"github.com/lib/pq"
	"time"
)

type Movie struct {
	ID                int            `db:"id"`
	Name              string         `db:"name"`
	Languages         pq.StringArray `db:"languages"`
	Dimensions        pq.StringArray `db:"dimensions"`
	ReleaseDate       time.Time      `db:"release_date"`
	DurationInMinutes int            `db:"duration_in_minutes"`
	Banners           []Image        `db:"banners"`
}

type MovieShowDetails struct {
	ID                int            `db:"id"`
	Name              string         `db:"name"`
	ReleaseDate       time.Time      `db:"release_date"`
	Dimensions        pq.StringArray `db:"dimensions"`
	Languages         pq.StringArray `db:"languages"`
	City              City           `db:"city"`
	DurationInMinutes int            `db:"duration_in_minutes"`
	ShowDetails       []ShowDetails  `db:"-"`
}

type MovieDetails struct {
	ID                int            `db:"id"`
	Name              string         `db:"name"`
	Dimensions        pq.StringArray `db:"dimensions"`
	Language          pq.StringArray `db:"languages"`
	ReleaseDate       time.Time      `db:"release_date"`
	DurationInMinutes int            `db:"duration_in_minutes"`
	Banners           []Image        `db:"banners"`
	Cast              []Cast         `db:"cast"`
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

type Theater struct {
	ID       int      `db:"id"`
	Name     string   `db:"name"`
	City     string   `db:"city"`
	Address  string   `db:"string"`
}

type ShowDetails struct {
	ShowID          int       `db:"show_id"`
	ShowStartTime   time.Time `db:"show_start_time"`
	ShowEndTime     time.Time `db:"show_end_time"`
	TheaterID       int       `db:"theater_id"`
	TheaterName     string    `db:"theater_name"`
	Dimension       string    `db:"dimension"`
	Language        string    `db:"language"`
}

type ShowSeats struct {
	HallID      int `db:"hall_id"`
	TotalRows   int `db:"total_rows"`
	TotalColumn int `db:"total_columns"`
	TotalSeats  int `db:"total_seats"`
	SeatDetails []MovieShowSeatsDetails
}

type MovieShowSeatsDetails struct {
	ID         int        `db:"id"`
	SeatCode   string     `db:"seat_code"`
	Status     SeatStatus `db:"status"`
	RowPos     int        `db:"row_pos"`
	ColumnPos  int        `db:"column_pos"`
	IsRecliner bool       `db:"is_recliner"`
	IsBooked   bool       `db:"is_booked"`
}

type Booking struct {
	ID            int            `db:"id"`
	BookingID     string         `db:"booking_id"`
	Status        BookingStatus  `db:"booking_status"`
	MovieID       int            `db:"movie_id"`
	MovieName     string         `db:"movie_name"`
	TheaterID     int            `db:"theater_id"`
	TheaterName   string         `db:"theater_name"`
	HallID        int            `db:"hall_id"`
	HallName      string         `db:"hall_name"`
	ShowID        int            `db:"show_id"`
	ShowStartTime time.Time      `db:"show_start_time"`
	ShowEndTime   time.Time      `db:"show_end_time"`
	Seats         []BookingSeats `db:"-"`
}

type BookingSeats struct {
	ID         int    `db:"id"`
	SeatsCode  string `db:"seat_code"`
	IsRecliner bool   `db:"is_recliner"`
	BookingID  int    `db:"booking_id"`
}
