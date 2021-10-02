package converter

import (
	"github.com/gauravcoco/bms/dbModels"
	"github.com/gauravcoco/bms/models"
)

func (c *converter) ToBookingStatus(status dbModels.BookingStatus) models.BookingStatus {
	switch status {
	case dbModels.BookingStatusConfirmed:
		return models.BookingStatusConfirmed
	case dbModels.BookingStatusCancelled:
		return models.BookingStatusCancelled
	}
	return ""
}
