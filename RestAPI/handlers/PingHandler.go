package handlers

import (
	"Golang-Templates/RestAPI/models"
	"net/http"

	dbhelper "github.com/JojiiOfficial/GoDBHelper"
)

//Ping handles ping request
func Ping(db *dbhelper.DBhelper, handlerData handlerData, w http.ResponseWriter, r *http.Request) {
	var request models.PingRequest
	if !parseUserInput(handlerData.config, w, r, &request) {
		return
	}

	payload := "pong"

	auth := NewAuthHandler(r, nil)
	if len(auth.GetBearer()) > 0 {
		payload = "Authorized pong"
	}

	response := models.StringResponse{
		String: payload,
	}
	sendResponse(w, models.ResponseSuccess, "", response)
}
