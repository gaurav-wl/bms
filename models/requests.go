package models

import "time"

type NewUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type BookingRequest struct {
	UserID    int   `json:"userID"`
	ShowID    int   `json:"showID"`
	SeatIDs   []int `json:"seatIDs"`
}

type BookedMovieDetail struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Language          string `json:"language"`
	DurationInMinutes int    `json:"duration_in_minutes"`
}

type BookingDetails struct {
	Movie      BookedMovieDetail `json:"movie"`
	Status     BookingStatus     `json:"status"`
	Theater    Theater           `json:"theater"`
	Hall       string            `json:"hall"`
	BookingUID string            `json:"bookingUID"`
	SeatCodes  []SeatCodes       `json:"seatCodes"`
	ShowTiming ShowTiming        `json:"showTiming"`
}

type SeatCodes struct {
	Name       string `json:"name"`
	IsRecliner bool   `json:"isRecliner"`
}

type ShowTiming struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

type ShowSeats struct {
	HallID      int `db:"hall_id"`
	TotalRows   int `db:"total_rows"`
	TotalColumn int `db:"total_column"`
	TotalSeats  int `db:"total_seats"`
	SeatDetails []MovieShowSeatsDetails
}

type MovieShowSeatsDetails struct {
	ID         int    `db:"id"`
	SeatCode   string `db:"seat_code"`
	RowPos     int    `db:"row_pos"`
	ColumnPos  int    `db:"column_pos"`
	IsRecliner bool   `db:"is_recliner"`
	IsBooked   bool   `db:"is_booked"`
}
