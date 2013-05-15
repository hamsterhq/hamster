package hamster

import (
	"net/http"
)

//check username-password, return with access-token
var DevAuth = func(w http.ResponseWriter, r *http.Request) {
	/*if r.URL.User == nil || r.URL.User.Username() != "admin" {
		http.Error(w, "", http.StatusUnauthorized)
	}*/

}

//check username-password, return with access-token
var AppAuth = func(w http.ResponseWriter, r *http.Request) {
	/*if r.URL.User == nil || r.URL.User.Username() != "admin" {
		http.Error(w, "", http.StatusUnauthorized)
	}*/

}

//check user agent
//for browser clients: check access-token and shared-secret
//for other clients: check app-id and app-key

var APIAuth = func(w http.ResponseWriter, r *http.Request) {

}
