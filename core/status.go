package core

import (
	"net/http"
)

func (s *Server) serveError(w http.ResponseWriter, err error, userMessage string, baseMessage string, status int) {

	logMessage := baseMessage + " : " + userMessage
	s.logger.Printf("Error %v:  %v \n", logMessage, err)
	http.Error(w, baseMessage, status)

}
func (s *Server) badRequest(r *http.Request, w http.ResponseWriter, err error, msg string) {
	s.serveError(w, err, msg, "Bad Request", http.StatusBadRequest)
}

func (s *Server) unauthorized(r *http.Request, w http.ResponseWriter, err error, msg string) {
	s.serveError(w, err, msg, "Unauthorized", http.StatusUnauthorized)
}

func (s *Server) notFound(r *http.Request, w http.ResponseWriter, err error, msg string) {
	s.serveError(w, err, msg, "Not found", http.StatusNotFound)

}

func (s *Server) internalError(r *http.Request, w http.ResponseWriter, err error, msg string) {
	s.serveError(w, err, msg, "Internal Server Error", http.StatusInternalServerError)

}
