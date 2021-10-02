package server

import (
	"encoding/json"
	"errors"
	"github.com/gauravcoco/bms/errors"
	"github.com/gauravcoco/bms/models"
	"github.com/gauravcoco/bms/utils"
	"net/http"
	"strconv"
)

func (srv *Server) getAllBookings(resp http.ResponseWriter, req *http.Request) {
	uc := srv.GetUserContext(req)
	if uc == nil {
		bmsError.RespondClientErr(resp, req, errors.New("error getting user context"), http.StatusUnauthorized, "Error getting user")
		return
	}

	bookings, err := srv.DBHelper.GetAllUserBookings(req.Context(), uc.ID)
	if err != nil {
		bmsError.RespondClientErr(resp, req, err, http.StatusInternalServerError, "Error getting movie details", "Error getting bookings#"+strconv.Itoa(uc.ID))
		return
	}

	utils.EncodeJSONBody(resp, http.StatusOK, map[string]interface{}{
		"movieDetails": srv.Converter.ToBookings(bookings),
	})
}


func (srv *Server) book(resp http.ResponseWriter, req *http.Request) {
	uc := srv.GetUserContext(req)
	if uc == nil {
		bmsError.RespondClientErr(resp, req, errors.New("error getting user context"), http.StatusUnauthorized, "Error getting user")
		return
	}

	var bookRequest models.BookingRequest

	err := json.NewDecoder(req.Body).Decode(&bookRequest)
	if err != nil {
		bmsError.RespondClientErr(resp, req, errors.New("error decoding request"), http.StatusInternalServerError, "Failed to book tickets")
		return
	}

	bookingDetails, err := srv.DBHelper.BookMovieTicket(req.Context(), bookRequest)
	if err != nil {
		bmsError.RespondClientErr(resp, req, err, http.StatusInternalServerError, "Error getting movie details", "Error getting bookings#"+strconv.Itoa(uc.ID))
		return
	}

	utils.EncodeJSONBody(resp, http.StatusCreated, map[string]interface{}{
		"movieDetails": srv.Converter.ToBooking(bookingDetails),
	})
}
