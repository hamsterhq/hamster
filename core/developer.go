package core

import (
	"fmt"
	"net/http"
	"time"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

var (
	dName = "developers"
)

//Stores developer account info
type developer struct {
	ID       bson.ObjectId `bson:"_id" json:"id"`
	ParentID string        `bson:"parentId" json:"parentId"` //unused. change string to bson.ObjectId
	Name     string        `bson:"name" json:"name"`
	Email    string        `bson:"email" json:"email"`
	Verified bool          `bson:"verified" json:"verified"`
	Password string        `json:"password"` //only used for parsing incoming json
	Hash     string        `bson:"hash"`
	Salt     string        `bson:"salt"`
	Created  time.Time     `bson:"created" json:"created"`
	Updated  time.Time     `bson:"updated" json:"updated"`
	URLToken string        `bson:"urltoken" json:"urltoken"`
}

//pre-index developers collection on startup. probably should do it through command line?
func (s *Server) indexDevelopers() {
	s.logger.SetPrefix("IndexDev: ")
	//ensure email exists and is unique
	index := mgo.Index{
		Key:        []string{"email"},
		Unique:     true,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	}

	//get collection
	session := s.db.GetSession()
	defer session.Close()
	c := session.DB("").C(dName)

	//ensure index key exists and is unique

	if err := c.EnsureIndex(index); err != nil {

		s.logger.Printf("failed indexing developers, err: %v \n", err)
	}

}

//POST: /api/v1/developers/ handler
func (s *Server) createDev(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("CreateDev: ")
	//get the request body
	developer := &developer{}

	if err := s.readJSON(developer, r, w); err != nil {
		s.badRequest(r, w, err, "malformed developer json")
		return

	}

	//check if email is not empty
	if developer.Email == "" {
		s.internalError(r, w, nil, "empty email ")
		return
	}

	//get collection
	session := s.db.GetSession()
	defer session.Close()
	c := session.DB("").C(dName)

	//set fields
	developer.ID = bson.NewObjectId() //todo:make it shorter and user friendly
	developer.URLToken = encodeBase64Token(developer.ID.Hex())
	//encrypt password
	hash, salt, err := encryptPassword(developer.Password)
	if err != nil {
		s.internalError(r, w, err, "encypt password")
		return

	}
	developer.Hash = hash
	developer.Salt = salt
	developer.Created = time.Now()
	developer.Updated = time.Now()

	//insert new document

	if insertErr := c.Insert(developer); insertErr != nil {

		s.internalError(r, w, insertErr, "error inserting: "+fmt.Sprintf("%v", developer))

	} else {
		s.logger.Printf("created new developer: %+v", developer)
		//serve created developer json
		accessToken, err := s.genAccessToken(developer.Email)
		if err != nil {
			s.internalError(r, w, err, "error generating access token")
		}
		response := NewDeveloperResponse{ObjectID: developer.URLToken, AccessToken: accessToken}
		s.serveJSON(w, &response)
	}

	return

}

//POST: /api/v1/developers/login/ handler
func (s *Server) loginDev(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("LoginDev: ")

	//get email password from request
	email, password := getUserPassword(r)

	//get collection
	session := s.db.GetSession()
	defer session.Close()
	c := session.DB("").C(dName)

	developer := developer{}

	//find and login developer
	s.logger.Printf("find developer: %s %s", email, password)
	if email != "" && password != "" {
		if findErr := c.Find(bson.M{"email": email}).One(&developer); findErr != nil {
			s.notFound(r, w, findErr, email+" user not found")
			return
		}
		//match password and set session
		if matchPassword(password, developer.Hash, developer.Salt) {
			accessToken, err := s.genAccessToken(developer.Email)
			if err != nil {
				s.internalError(r, w, err, email+" generate access token")
			}
			//respond with developer profile
			response := loginResponse{ObjectID: developer.URLToken, AccessToken: accessToken, Status: "ok"}
			s.serveJSON(w, &response)

		} else {

			s.notFound(r, w, nil, email+" password match failed")
		}

	} else {
		s.notFound(r, w, nil, "email empty")
	}

}

//POST:/api/v1/developers/logout/ handler
func (s *Server) logoutDev(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("LogoutDev: ")

	//parse body
	logoutRequest := &logoutRequest{}

	if err := s.readJSON(logoutRequest, r, w); err != nil {
		s.badRequest(r, w, err, "malformed logout request")
		return
	}

	//logout
	if logoutErr := s.logout(logoutRequest.Email); logoutErr != nil {
		s.internalError(r, w, logoutErr, logoutRequest.Email+" : could not logout")
	}

	//response
	response := logoutResponse{Status: "ok"}
	s.serveJSON(w, &response)

}

//GET:/api/v1/developers/:objectID handler
func (s *Server) queryDev(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("QueryDev: ")

	//getObjectId
	objectID := s.getObjectID(w, r)

	//get collection
	session := s.db.GetSession()
	defer session.Close()
	c := session.DB("").C(dName)

	//find and serve data
	developer := developer{}
	if err := c.FindId(bson.ObjectIdHex(objectID)).Limit(1).One(&developer); err != nil {
		s.notFound(r, w, err, objectID+" : id not found")
		return
	}

	//respond
	response := queryDevResponse{Name: developer.Name, Email: developer.Email}
	s.serveJSON(w, &response)

}

//PUT:/api/v1/developers/:objectID handler
func (s *Server) updateDev(w http.ResponseWriter, r *http.Request) {

	s.logger.SetPrefix("UpdateDev: ")

	//getObjectId
	objectID := s.getObjectID(w, r)

	//parse body
	updateRequest := &updateRequest{}
	if err := s.readJSON(updateRequest, r, w); err != nil {
		s.badRequest(r, w, err, "malformed update request body")
		return
	}

	//get collection
	session := s.db.GetSession()
	defer session.Close()
	c := session.DB("").C(dName)

	//change
	var change = mgo.Change{
		ReturnNew: true,
		Update: bson.M{
			"$set": bson.M{
				"updated": time.Now(),
				"name":    updateRequest.Name,
			}}}

	//find and update
	developer := developer{}
	if _, err := c.FindId(bson.ObjectIdHex(objectID)).Apply(change, &developer); err != nil {
		s.notFound(r, w, err, objectID+" : id not found")
		return
	}

	//respond
	response := queryDevResponse{Name: developer.Name, Email: developer.Email}
	s.serveJSON(w, &response)

}

//DELETE:/api/v1/developers/:objectID
func (s *Server) deleteDev(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("DeleteDev: ")

	//getObjectId
	objectID := s.getObjectID(w, r)

	//get collection
	session := s.db.GetSession()
	defer session.Close()
	c := session.DB("").C(dName)

	//delete
	if err := c.RemoveId(bson.ObjectIdHex(objectID)); err != nil {
		s.notFound(r, w, err, objectID+" : id not found")
		return
	}

	//respond
	response := deleteResponse{Status: "ok"}
	s.serveJSON(w, &response)

}

func (s *Server) info(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hamster api")
}
