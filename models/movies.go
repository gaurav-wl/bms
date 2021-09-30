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
	Dimension         []string      `json:"dimension"`
	Language          []string      `json:"languages"`
	City              City          `json:"city"`
	DurationInMinutes int           `json:"duration_in_minutes"`
	ShowDetails       []ShowDetails `json:"showDetails"`
}

type Movie struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	Languages         []string  `json:"languages"`
	ReleaseDate       time.Time `json:"releaseDate"`
	DurationInMinutes int       `json:"duration_in_minutes"`
	Banners           []string  `json:"banners"` // Image URLS
}

type MovieDetails struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	Dimension         []string  `json:"dimension"`
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

type Point struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type ShowDetails struct {
	Theater Theater       `json:"theater"`
	Movie   Movie         `json:"movie"`
	Shows   []ShowTimings `json:"shows"`
}

type Theater struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	City     string `json:"city"`
	Location Point  `json:"location"`
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
