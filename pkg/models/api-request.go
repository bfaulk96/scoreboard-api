package models

import (
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"io/ioutil"
	"errors"
)

type APIRequest struct {
	*http.Request
}

func NewAPIRequest(r *http.Request) (ar *APIRequest) {
	ar = new(APIRequest)
	ar.Request = r
	log.Println("Call received: \"" + ar.Method + " " + ar.URL.Path + "\"")
	return ar
}

func (ar *APIRequest) GetRouteVariables() (routeVariables map[string]string) {
	return mux.Vars(ar.Request)
}

func (ar *APIRequest) GetQueryParameters() (queryParameters map[string][]string) {
	return ar.URL.Query()
}

func (ar *APIRequest) GetRequestBody() (rawRequestBody []byte, err error) {
	rawRequestBody, err = ioutil.ReadAll(ar.Body)
	if err != nil {
		log.Printf("Request Error: %v\n", err.Error())
		return nil, errors.New("Could not read request body.")
	}
	return rawRequestBody, nil
}