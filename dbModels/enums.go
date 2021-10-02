package dbModels

type BookingStatus string
type SeatStatus string

const (
	BookingStatusConfirmed   BookingStatus = "confirmed"
	BookingStatusCancelled   BookingStatus = "cancelled"
	BookingStatusUnConfirmed BookingStatus = "unconfirmed"

	SeatStatusFunctional    BookingStatus = "functional"
	SeatStatusNotFunctional BookingStatus = "not_functional"
)
