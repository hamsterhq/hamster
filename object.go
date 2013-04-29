package hamster

import (
	"net/http"
)

//create new user defined object
func (s *Server) CreateObject(w http.ResponseWriter, r *http.Request) {
	//login to db(username:password?)
	//defer db session
	//find collection: className, if found create document with request body
	//otherwise create collection

}

//get object
func (s *Server) GetObject(w http.ResponseWriter, r *http.Request) {

}

//update object
func (s *Server) UpdateObject(w http.ResponseWriter, r *http.Request) {

}

//query object
func (s *Server) QueryObject(w http.ResponseWriter, r *http.Request) {

}

//get object
func (s *Server) DeleteObject(w http.ResponseWriter, r *http.Request) {

}
