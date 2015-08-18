package core

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (s *Server) readJSON(d interface{}, r *http.Request, w http.ResponseWriter) error {

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {

		s.logger.Printf("error in reading body for: %v, err: %v\n ", r.Body, err)
		http.Error(w, "Bad Data!", http.StatusBadRequest)
		return err
	}

	return json.Unmarshal(body, &d)

}

func (s *Server) serveJSON(w http.ResponseWriter, v interface{}) {
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

func (s *Server) getObjectID(w http.ResponseWriter, r *http.Request) string {
	oid := r.URL.Query().Get(":objectId")
	if oid == "" {
		s.notFound(r, w, errors.New("objectId is empty"), "val: "+oid)
	}

	objectID := decodeToken(oid)
	if objectID == "" {
		s.notFound(r, w, errors.New("objectId cannot be decoded"), "val: "+oid)
	}

	return objectID

}

func (s *Server) getObjectName(w http.ResponseWriter, r *http.Request) string {
	oname := r.URL.Query().Get(":objectName")
	if oname == "" {
		s.notFound(r, w, errors.New("objectName is empty"), "val: "+oname)
	}

	return oname

}

func (s *Server) getFileName(w http.ResponseWriter, r *http.Request) string {
	fname := r.URL.Query().Get(":fileName")
	if fname == "" {
		s.notFound(r, w, errors.New("fileName is empty"), "val: "+fname)
	}

	return fname

}

func (s *Server) getFileParams(w http.ResponseWriter, r *http.Request) (string, string) {
	fname := r.URL.Query().Get(":fileName")
	if fname == "" {
		s.notFound(r, w, errors.New("fileName is empty"), "val: "+fname)
	}

	fid := r.URL.Query().Get(":fileId")
	if fid == "" {
		s.notFound(r, w, errors.New("fileId is empty"), "val: "+fid)
	}

	fileID := decodeToken(fid)
	if fileID == "" {
		s.notFound(r, w, errors.New("object params cannot be decoded"), "val: "+fname+" , "+fileID)
	}

	return fname, fileID

}

func (s *Server) getObjectParams(w http.ResponseWriter, r *http.Request) (string, string) {
	objectName := r.URL.Query().Get(":objectName")
	oid := r.URL.Query().Get(":objectId")
	if oid == "" || objectName == "" {
		s.notFound(r, w, errors.New("object params are invalid"), "val: "+objectName+" , "+oid)
	}

	objectID := decodeToken(oid)
	if objectID == "" {
		s.notFound(r, w, errors.New("object params cannot be decoded"), "val: "+objectName+" , "+oid)
	}

	return objectName, objectID

}

func (s *Server) getAppObjectID(w http.ResponseWriter, r *http.Request) string {
	atok := r.Header.Get("X-Api-Token")

	if atok == "" {
		s.unauthorized(r, w, errors.New("token is empty"), "api token invalid")
	}

	objectID := decodeToken(atok)
	if objectID == "" {
		s.notFound(r, w, errors.New("app objectid cannot be decoded"), "val: "+atok)
	}

	return objectID

}

func (s *Server) getAppParams(w http.ResponseWriter, r *http.Request) (string, string) {
	did := r.URL.Query().Get(":developerId")
	oid := r.URL.Query().Get(":objectId")
	if oid == "" || did == "" {
		s.notFound(r, w, errors.New("app params are invalid"), "val: "+did+" , "+oid)
	}

	developerID := decodeToken(did)
	objectID := decodeToken(oid)
	if objectID == "" || developerID == "" {
		s.notFound(r, w, errors.New("app params cannot be decoded"), "val: "+did+" , "+oid)
	}

	return developerID, objectID

}
