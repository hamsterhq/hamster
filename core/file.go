package core

import (
	"bufio"
	"bytes"
	"encoding/json"
	"image"
	"image/png"
	"io"
	"net/http"

	"labix.org/v2/mgo/bson"
)

//Saves file read from request body
//POST:/api/v1/files/:fileName
func (s *Server) saveFile(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("SaveFile: ")

	//get file params and file data reader
	fileName := s.getFileName(w, r)
	fileReader := bufio.NewReader(r.Body)
	defer r.Body.Close()

	//get meta data
	metaDataJSON := r.Header.Get("X-Meta-Data")

	//parse meta data
	var metadata map[string]interface{}
	json.Unmarshal([]byte(metaDataJSON), &metadata)

	//get session
	session := s.db.GetSession()
	defer session.Close()
	db := session.DB("")

	//create file
	file, err := db.GridFS("fs").Create(fileName)
	if err != nil {
		s.internalError(r, w, err, "could not create file")
	}

	//copy incoming data to file
	_, err = io.Copy(file, fileReader)
	if err != nil {
		s.internalError(r, w, err, "could not copy file")
	}

	//set content type and meta deta
	file.SetContentType("image/png")
	file.SetMeta(metadata)

	//encode file id and serve
	fileID := encodeBase64Token(file.Id().(bson.ObjectId).Hex())
	response := SaveFileResponse{FileID: fileID, FileName: file.Name()}

	err = file.Close()
	if err != nil {
		s.internalError(r, w, err, "could not close file")
	}

	//respond

	s.serveJSON(w, &response)

}

//Gets file from GridFS and writes to response body
//GET:/api/v1/files/:fileName/:fileId handler
func (s *Server) getFile(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("GetFile: ")

	//get file params
	_, fileID := s.getFileParams(w, r)

	//get session
	session := s.db.GetSession()
	defer session.Close()
	db := session.DB("")

	//open file from GridFS
	file, err := db.GridFS("fs").OpenId(bson.ObjectIdHex(fileID))
	if err != nil {
		s.internalError(r, w, err, "could not open file")
	}

	//copy buffer
	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		s.internalError(r, w, err, "could not copy buffer")
	}

	//decode buffer
	img, _, err := image.Decode(&buf)

	if err != nil {
		s.internalError(r, w, err, "could not decode image")
	}

	contentType := file.ContentType()
	//fmt.Printf("content type: %v \n", contentType)

	err = file.Close()
	if err != nil {
		s.internalError(r, w, err, "could not close file")
	}

	//set content type and write to response body
	w.Header().Set("Content-Type", contentType)
	//TODO encode by mime type
	png.Encode(w, img)

}
