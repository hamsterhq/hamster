package hamster

import (
	"net/http"
)

type File struct {
	FileName string
}

func (f *File) SaveFile(w http.ResponseWriter, r *http.Request) {

}

func (f *File) GetFile(w http.ResponseWriter, r *http.Request) {

}
