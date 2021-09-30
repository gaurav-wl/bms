package bmsError

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

const (
	NotAuthorizedMsg = "User Not Authorized for this resource."
	ServerErrorMsg   = "Internal Server Error occurred when processing request."
	NoUserMsg        = "Couldn't get user information"
	BadRequestFormat = "Invalid Request Body"
)

type clientError struct {
	ID            string `json:"id"`
	MessageToUser string `json:"messageToUser"`
	DeveloperInfo string `json:"developerInfo"`
	Err           string `json:"error"`
	StatusCode    int    `json:"statusCode"`
	IsClientError bool   `json:"isClientError"`
}

func (c clientError) Error() string {
	return c.Err
}

func (c *clientError) LogFields() logrus.Fields {
	return logrus.Fields{
		"ID":            c.ID,
		"MessageToUser": c.MessageToUser,
		"DeveloperInfo": c.DeveloperInfo,
	}
}

func (c *clientError) LogMessage() string {
	return fmt.Sprintf("[BMSClientErr]: id(%s) %s : %+v", c.ID, c.MessageToUser+" "+c.DeveloperInfo, c.Err)
}

func NewClientError(err error, statusCode int, req *http.Request, messageToUser string, additionalInfoForDevs ...string) *clientError {
	additionalInfoJoined := strings.Join(additionalInfoForDevs, "\n")
	if len(additionalInfoJoined) == 0 {
		additionalInfoJoined = messageToUser
	}

	var errString string
	if err != nil {
		errString = err.Error()
	}
	return &clientError{
		ID:            middleware.GetReqID(req.Context()),
		MessageToUser: messageToUser,
		DeveloperInfo: additionalInfoJoined,
		Err:           errString,
		StatusCode:    statusCode,
		IsClientError: true,
	}
}

func RespondClientErr(resp http.ResponseWriter, req *http.Request, err error, statusCode int, messageToUser string, additionalInfoForDevs ...string) {
	resp.WriteHeader(statusCode)
	clientError := NewClientError(err, statusCode, req, messageToUser, additionalInfoForDevs...)
	if statusCode >= 400 && statusCode < 500 {
		logrus.WithContext(req.Context()).Warn(clientError.LogMessage())
	} else {
		logrus.WithContext(req.Context()).Error(clientError.LogMessage())
	}
	if err := json.NewEncoder(resp).Encode(clientError); err != nil {
		logrus.Error(err)
	}
}

func RespondGenericServerErr(resp http.ResponseWriter, req *http.Request, err error, additionalInfoForDevs ...string) {
	resp.WriteHeader(http.StatusInternalServerError)
	additionalInfoJoined := strings.Join(additionalInfoForDevs, "\n")
	clintErr := &clientError{
		ID:            middleware.GetReqID(req.Context()),
		MessageToUser: ServerErrorMsg,
		DeveloperInfo: additionalInfoJoined,
		Err:           err.Error(),
		StatusCode:    http.StatusInternalServerError,
		IsClientError: false,
	}
	logrus.WithContext(req.Context()).Error(clintErr.LogMessage())
	if err := json.NewEncoder(resp).Encode(clintErr); err != nil {
		logrus.Error(err)
	}
}
