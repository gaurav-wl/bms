package utils

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
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

	return "URL:" + bucket + "::" + path
}

func ToInt64Slice(ints []int) []int64 {
	int64Slice := make([]int64, len(ints))
	for index, i := range ints {
		int64Slice[index] = int64(i)
	}
	return int64Slice
}

func ToDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
