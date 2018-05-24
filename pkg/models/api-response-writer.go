package models

import (
	"net/http"
	"encoding/json"
	"log"
	"strconv"
)

type APIResponseWriter struct {
	http.ResponseWriter
}

func NewAPIResponseWriter(w http.ResponseWriter) (aw *APIResponseWriter) {
	aw = new(APIResponseWriter)
	aw.ResponseWriter = w
	return aw
}

func (aw *APIResponseWriter) Respond(r *APIRequest, response interface{}, responseStatus int) () {
	responseBody, err := json.Marshal(response)
	if err != nil {
		log.Printf("JSON Marshal error: %v\n", err.Error())
		aw.Header().Set("Content-Type", "text/plain")
		aw.WriteHeader(http.StatusInternalServerError)
		aw.Write([]byte("{\n\t\"error\": \"Could not process response body.\"\n}"))
		return
	}

	aw.Header().Set("Content-Type", "application/json")
	aw.WriteHeader(responseStatus)
	aw.Write([]byte(responseBody))
	log.Printf("Response sent: " + strconv.Itoa(responseStatus) + ": \"" + r.Method + " " + r.URL.Path + "\"\n")
}