package hamster

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"time"
)

var (
	cName = "developers"
)

//stores developer account info
type Developer struct {
	Id       bson.ObjectId `bson:"_id" json:"id"`
	ParentId string        `bson:"parentId" json:"parentId"`
	Name     string        `bson:"name" json:"name"`
	Email    string        `bson:"email" json:"email"`
	Verified bool          `bson:"verified" json:"verified"`
	Password string        `bson:"password" json:"password"`
	Created  time.Time     `bson:"created" json:"created"`
	Updated  time.Time     `bson:"updated" json:"updated"`
	UrlToken string        `bson:"urltoken" json:"urltoken"`
}

//pre-index developers collection on startup. probably should do it through command line?
func (s *Server) IndexDevelopers() {
	s.logger.SetPrefix("IndexDev: ")
	//ensure email exists and is unique
	index := mgo.Index{
		Key:        []string{"email"},
		Unique:     true,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	}
	session := s.db.GetSession()
	defer session.Close()
	c := session.DB("").C(cName)

	//ensure index key exists and is unique

	err := c.EnsureIndex(index)
	if err != nil {

		s.logger.Printf("failed indexing developers, err: %v \n", err)
	} else {
		s.logger.Printf("developers collection indexed!")
	}

}

//POST: /developers/
func (s *Server) CreateDev(w http.ResponseWriter, r *http.Request) {

	s.logger.SetPrefix("CreateDev: ")
	//get the request body
	developer := &Developer{}

	marshalError := s.readJson(developer, r, w)
	if marshalError != nil {

		s.logger.Printf("error in parsing body for: %v, err: %v \n", r.Body, marshalError)

		http.Error(w, "Bad Data!", http.StatusBadRequest)

	}

	//check if email is not empty
	if developer.Email == "" {
		s.logger.Printf("error in inserting data for: %v, email is empty \n", developer)
		http.Error(w, "email is empty", http.StatusInternalServerError)
		return
	}

	//encrypt password

	//profile details
	session := s.db.GetSession()
	defer session.Close()
	c := session.DB("").C(cName)

	//set fields
	developer.Id = bson.NewObjectId()
	s.logger.Printf("bson objectId %v\n", developer.Id)

	// UrlToken: public object id used for routes and queries
	//todo:make it shorter and user friendly

	developer.UrlToken = s.encodeBase64Token(developer.Id.Hex())
	developer.Created = time.Now()
	developer.Updated = time.Now()
	//insert new document
	insertError := c.Insert(developer)
	if insertError != nil {

		s.logger.Printf("error in inserting data for: %v, err: %v \n", developer, insertError)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return

	} else {
		s.logger.Printf("created new developer: %+v, id: %v\n", developer)
		//serve created developer json
		response := Response{S: 200, D: &developer}
		s.serveJson(w, &response)
	}

	return

}

//get developer
func (s *Server) GetDev(w http.ResponseWriter, r *http.Request) {

}

//login developer
func (s *Server) LoginDev(w http.ResponseWriter, r *http.Request) {

}

//query developer
func (s *Server) QueryDev(w http.ResponseWriter, r *http.Request) {

}

//update developer
func (s *Server) UpdateDev(w http.ResponseWriter, r *http.Request) {

}

//get developer
func (s *Server) DeleteDev(w http.ResponseWriter, r *http.Request) {

}
