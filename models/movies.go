package models

import (
	"github.com/lib/pq"
	"github.com/volatiletech/null"
	"time"
)

type MovieSearchRequest struct {
	Name       string         `json:"name"`
	Dimensions pq.StringArray `json:"dimension"`
	Languages  pq.StringArray `json:"languages"`
	CityID     null.Int       `json:"cityID"`
}

type MovieShowDetails struct {
	ID                int           `json:"id"`
	Name              string        `json:"name"`
	Dimensions        []string      `json:"dimensions"`
	Language          []string      `json:"languages"`
	City              City          `json:"city"`
	DurationInMinutes int           `json:"duration_in_minutes"`
	ShowDetails       []ShowDetails `json:"showDetails"`
}

type Movie struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	Languages         []string  `json:"languages"`
	Dimensions        []string  `json:"dimensions"`
	ReleaseDate       time.Time `json:"releaseDate"`
	DurationInMinutes int       `json:"duration_in_minutes"`
	Banners           []string  `json:"banners"` // Image URLS
}

type MovieDetails struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	Dimensions        []string  `json:"dimensions"`
	Language          []string  `json:"languages"`
	ReleaseDate       time.Time `json:"releaseDate"`
	DurationInMinutes int       `json:"duration_in_minutes"`
	Banners           []string  `json:"banners"` // Image URLS
	Cast              []Cast    `json:"cast"`
}

type Cast struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

type City struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ShowDetails struct {
	Theater Theater       `json:"theater"`
	Shows   []ShowTimings `json:"shows"`
}

type Theater struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	City     string `json:"city"`
	Location *LongLat  `json:"location"`
	Address  string `json:"string"`
}

type ShowTimings struct {
	Date  time.Time `json:"date"`
	Shows []Shows   `json:"shows"`
}

type Shows struct {
	ID        int       `json:"id"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Dimension string    `json:"dimension"`
	Language  string    `json:"language"`
}

type LongLat struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}
