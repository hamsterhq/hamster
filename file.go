package hamster

import (
	"bufio"
	"bytes"
	"encoding/json"
	//"fmt"
	"image"
	"image/png"
	"io"
	"labix.org/v2/mgo/bson"
	"net/http"
)

type File struct {
	FileName string
}

func (s *Server) SaveFile(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("SaveFile: ")
	file_name := s.getFileName(w, r)
	fileReader := bufio.NewReader(r.Body)
	defer r.Body.Close()

	//get meta data
	meta_data_json := r.Header.Get("X-Meta-Data")

	var metadata map[string]interface{}
	json.Unmarshal([]byte(meta_data_json), &metadata)

	session := s.db.GetSession()
	defer session.Close()
	db := session.DB("")

	file, err := db.GridFS("fs").Create(file_name)
	if err != nil {
		s.internalError(r, w, err, "could not create file")
	}

	_, err = io.Copy(file, fileReader)
	if err != nil {
		s.internalError(r, w, err, "could not copy file")
	}

	file.SetContentType("image/png")
	file.SetMeta(metadata)

	file_id := encodeBase64Token(file.Id().(bson.ObjectId).Hex())
	response := SaveFileResponse{FileId: file_id, FileName: file.Name()}

	err = file.Close()
	if err != nil {
		s.internalError(r, w, err, "could not close file")
	}

	//respond

	s.serveJson(w, &response)

}

func (s *Server) GetFile(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("GetFile: ")
	_, file_id := s.getFileParams(w, r)

	session := s.db.GetSession()
	defer session.Close()
	db := session.DB("")

	file, err := db.GridFS("fs").OpenId(bson.ObjectIdHex(file_id))
	if err != nil {
		s.internalError(r, w, err, "could not open file")
	}

	//copy buffer
	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		s.internalError(r, w, err, "could not copy buffer")
	}

	//get image
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

	w.Header().Set("Content-Type", contentType)
	//TODO encode by mime type
	png.Encode(w, img)

}
