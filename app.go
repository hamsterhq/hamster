package hamster

import (
	"labix.org/v2/mgo/bson"
	"net/http"
	"time"
)

type App struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"-"`
	ParentId  string        `bson:"parentId" json:"parentId"`
	Name      string        `bson:"name" json:"name"`
	OS        string        `bson:"os" json:"os"`
	ApiToken  string        `bson:"apitoken" json:"apitoken"`
	ApiSecret string        `bson:"apisecret" json:"apisecret"`
	Salt      string        `bson:"salt" json:"salt"`
	Created   time.Time     `bson:"created" json:"created"`
	Updated   time.Time     `bson:"updated" json:"updated"`
}

func (s *Server) CreateApp(w http.ResponseWriter, r *http.Request) {
	session := s.db.GetSession()
	defer session.Close()
	c := session.DB("").C("apps")

	app := &App{Name: "traverik", OS: "android"}

	err := c.Insert(app)
	if err != nil {

	} else {
		//serve created app json
	}

}

//get app
func (s *Server) GetApp(w http.ResponseWriter, r *http.Request) {

}

//update app
func (s *Server) UpdateApp(w http.ResponseWriter, r *http.Request) {

}

//query app
func (s *Server) QueryApp(w http.ResponseWriter, r *http.Request) {

}

//delete app
func (s *Server) DeleteApp(w http.ResponseWriter, r *http.Request) {

}
