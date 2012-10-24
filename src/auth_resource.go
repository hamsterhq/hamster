package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	BadRequest   = 400
	Unauthorized = 401
	NotFound     = 404
)

type Credentials string

func Authenticate(w http.ResponseWriter, r *http.Request) {
	credentials, err := newCredentials(r)
	if err != nil {
		http.Error(w, err.Error(), BadRequest)
		return
	}
	user, err := GetUserByApiKey(db, credentials)
	if err != nil {
		http.Error(w, err.Error(), Unauthorized)
		return
	}
	if user == nil {
		http.Error(w, "No user could be found", NotFound)
		return
	}
	bytes, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), Unauthorized)
		return
	}
	w.Write(bytes)
}

func newCredentials(r *http.Request) (Credentials, error) {
	_, authValue, err := getDecodedAuthorizationHeader(r.Header)
	if err != nil {
		return "", err
	}
	return Credentials(authValue), nil
}

func getDecodedAuthorizationHeader(headers http.Header) (aType string, aValue string, err error) {
	auth := strings.Split(headers.Get("Authorization"), " ")
	if len(auth) != 2 {
		return aType, aValue, errors.New("Bad auth header")
	}
	if auth[0] == "Basic" {
		reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(auth[1]))
		bytes, err := ioutil.ReadAll(reader)
		if err != nil {
			return aType, aValue, err
		}
		auth[1] = string(bytes)
	}
	return auth[0], auth[1], nil
}
