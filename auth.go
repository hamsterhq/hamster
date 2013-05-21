package hamster

import (
	"net/http"
)

//check access token + csrf secret(generated ont the client-server using the shared secret)
var DevAuth = func(w http.ResponseWriter, r *http.Request) {
	/*if r.URL.User == nil || r.URL.User.Username() != "admin" {
		http.Error(w, "", http.StatusUnauthorized)
	}*/

}

//check access token and csrf secret(fernet)
var AppAuth = func(w http.ResponseWriter, r *http.Request) {
	/*if r.URL.User == nil || r.URL.User.Username() != "admin" {
		http.Error(w, "", http.StatusUnauthorized)
	}*/

}

//check user agent
//for browser clients: check access-token and csrf token
//for other clients: check app-id and app-key

var APIAuth = func(w http.ResponseWriter, r *http.Request) {

}

//generates csrftoken using shared email + secret + timestamp + domain
func GenerateCSRFToken() {

}

func ValidateCSRFToken() {

}

//saves access_token to session
func SaveAccessToken() {

}

//validates an incoming access_token from session
func ValidatetAccessToken() {

}

//generate appid: appname + createdat
func GenerateAppID() {

}

func GenerateApiKey() {

}

func FindAppID() {

}

func ValidateApiKey() {

}
