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
	Salt     string        `bson:"salt" json:"salt"`
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
		return

	}

	//check if email is not empty
	if developer.Email == "" {
		s.logger.Printf("error in inserting data for: %v, email is empty \n", developer)
		http.Error(w, "email is empty", http.StatusInternalServerError)
		return
	}

	//encrypt password
	encrypted_password, salt, encryptPasswordError := encryptPassword(string(developer.Password))
	if encryptPasswordError != nil {

		s.logger.Printf("error encrypting password: %v \n", developer, encryptPasswordError)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return

	}

	s.logger.Printf("password-salt \n", encrypted_password, salt)

	//get db session
	session := s.db.GetSession()
	defer session.Close()
	c := session.DB("").C(cName)

	//set fields
	developer.Id = bson.NewObjectId()
	//s.logger.Printf("bson objectId %v\n", developer.Id)

	// UrlToken: public object id used for routes and queries
	//todo:make it shorter and user friendly
	developer.UrlToken = encodeBase64Token(developer.Id.Hex())
	developer.Password = encrypted_password
	developer.Salt = salt
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
		response := Response{C: encodeBase64Token(developer.Id.Hex()), S: 200}
		s.serveJson(w, &response)
	}

	return

}

//login developer. Generate access token(email+timestamp)
//store access_token in session
func (s *Server) LoginDev(w http.ResponseWriter, r *http.Request) { /*

		//get objectId

		id := r.URL.Query().Get(":id")

		//get db session
		session := s.db.GetSession()
		defer session.Close()
		c := session.DB("").C(cName)

		//find object
		var result Developer
		ctx := ""

		if objectId != "" {
			if findErr := c.FindId(decodeToken(objectId)).One(&result); findErr != nil {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
		} else if username != "" {
			if findErr := c.Find().One(&result); findErr != nil {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
		} else {
			http.Error(w, "no identifier", http.StatusNotFound)
			return
		}

		//respond with developer profile
		reponse := &Response{C: ctx, S: 200, D: &Developer{Name: result.Name, Email: result.Email, Verified: result.Verified}}*/

}

//logout developer

func (s *Server) LogoutDev(w http.ResponseWriter, r *http.Request) {

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
