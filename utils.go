package hamster

import (
	//"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"strconv"
	//"os"
	"strings"
	//"testing"
)

var (
	port        = 8686
	host        = "http://localhost:8686"
	mongoHost   = "mongodb://adnaan:pass@localhost:27017/hamster"
	contentType = "application/json"
)

func testHttpRequest(verb string, resource string, body string) (*http.Response, error) {
	client := &http.Client{Transport: &http.Transport{DisableKeepAlives: true}}
	r, _ := http.NewRequest(verb, fmt.Sprintf("%s%s", host, resource), strings.NewReader(body))
	r.Header.Add("Content-Type", contentType)
	return client.Do(r)
}

func testServer(f func(s *Server)) {

	server := NewServer(port, mongoHost)
	//server.Quiet()
	server.ListenAndServe()
	defer server.Shutdown()
	f(server)
}

func (s *Server) readJson(d interface{}, r *http.Request, w http.ResponseWriter) error {

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {

		s.logger.Printf("error in reading body for: %v, err: %v\n ", r.Body, err)
		http.Error(w, "Bad Data!", http.StatusBadRequest)
		return err
	}

	return json.Unmarshal(body, &d)

}

func (s *Server) serveJson(w http.ResponseWriter, v interface{}) {
	content, err := json.MarshalIndent(v, "", "  ")
	if err != nil {

		s.logger.Printf("error in serving json err: %v  \n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	w.Header().Set("Content-Type", "application/json")
	w.Write(content)
}

func (s *Server) encodeBase64Token(hexVal string) string {

	token := base64.URLEncoding.EncodeToString([]byte(hexVal))

	s.logger.Printf("encoded token: %s \n", token)
	return token

}

func (s *Server) decodeToken(token string) string {

	hexVal, err := base64.URLEncoding.DecodeString(token)
	if err != nil {

		s.logger.Printf("decoded token error: %v \n", err)
		return ""

	}

	s.logger.Printf("decoded token: %s \n", string(hexVal))
	return string(hexVal)

}
