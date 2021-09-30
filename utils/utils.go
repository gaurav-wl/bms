package utils

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

func EncodeJSONBody(resp http.ResponseWriter, statusCode int, data interface{}) {
	resp.WriteHeader(statusCode)
	err := json.NewEncoder(resp).Encode(data)
	if err != nil {
		logrus.Errorf("Error encoding response %v", err)
	}
	return
}

// Mock method for image URL
func GetImageURL(bucket, path string) string {
	return "URL"
}
