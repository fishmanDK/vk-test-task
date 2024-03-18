package handlers

import (
	"encoding/json"
	"github.com/fishmanDK/vk-test-task/internal/service/response"
	"net/http"
)

type keyLogger string

const (
	parseDataUser = "parseDataUser"
	admin         = "admin"

	accessDenied = "access denied"

	deleteSUCCESS = "removal was successful"
	changeSUCCESS = "change was successful"
	createSUCCESS = "creation was successful"

	loggerKey keyLogger = "logger"
)

func newErrorResponse(w http.ResponseWriter, status int, messageErr string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response.Error(messageErr))
	return
}
