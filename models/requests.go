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
	UserID        int   `json:"userID"`
	MovieID       int   `json:"movieID"`
	TheaterID     int   `json:"theaterID"`
	ShowIDs       []int `json:"showIDs"`
	SeatCodes     []int `json:"seatCodes"`
	NumberOfSeats int   `json:"numberOfSeats"`
}

type BookedMovieDetail struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Language          string `json:"language"`
	DurationInMinutes int    `json:"duration_in_minutes"`
}

type BookingDetails struct {
	Movie      BookedMovieDetail `json:"movie"`
	Status     string            `json:"status"`
	Theater    Theater           `json:"theater"`
	Hall       string            `json:"hall"`
	BookingUID string            `json:"bookingUID"`
	SeatCodes  []int             `json:"seatCodes"`
	ShowTiming ShowTiming        `json:"showTiming"`
}

type ShowTiming struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}
