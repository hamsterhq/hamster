package core

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

var aName = "apps"

//The App type
type app struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"-"`
	ParentID  bson.ObjectId `bson:"parentId" json:"parentId"`
	Name      string        `bson:"name" json:"name"`
	OS        string        `bson:"os" json:"os"`
	APIToken  string        `bson:"apitoken" json:"apitoken"`
	APISecret string        `bson:"apisecret" json:"apisecret"`
	Hash      string        `bson:"hash" json:"hash"`
	Salt      string        `bson:"salt" json:"salt"`
	Created   time.Time     `bson:"created" json:"created"`
	Updated   time.Time     `bson:"updated" json:"updated"`
	Objects   []string      `bson:"objects" json:"objects"`
}

//POST: "/api/v1/developers/:developerId/apps/" handler
func (s *Server) createApp(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("CreateApp: ")

	//get params
	did := r.URL.Query().Get(":developerId")
	if did == "" {
		s.notFound(r, w, errors.New("app params are invalid"), "val: "+did)
	}
	developerID := decodeToken(did)
	if developerID == "" {
		s.notFound(r, w, errors.New("app params are invalid"), "val: "+developerID)
	}

	//get collection developer
	session := s.db.GetSession()
	defer session.Close()
	d := session.DB("").C(dName)

	//if developer id exists

	if err := d.FindId(bson.ObjectIdHex(developerID)).Limit(1).One(nil); err != nil {
		s.notFound(r, w, err, developerID+" : id not found")
		return
	}

	//parse body
	app := &app{}
	if err := s.readJSON(app, r, w); err != nil {
		s.badRequest(r, w, err, "malformed app json")
		return

	}

	//set fields
	app.ID = bson.NewObjectId() //todo:make it shorter and user friendly
	app.ParentID = bson.ObjectIdHex(developerID)
	app.Created = time.Now()
	app.Updated = time.Now()
	app.APIToken = encodeBase64Token(app.ID.Hex())
	secret, err := genUUID(16)
	if err != nil {
		s.internalError(r, w, err, "gen uuid")
		return
	}
	app.APISecret = encodeBase64Token(secret)
	hash, salt, err := encryptPassword(secret)
	if err != nil {
		s.internalError(r, w, err, "encypt secret")
		return

	}
	app.Hash = hash
	app.Salt = salt

	//get apps collection
	c := session.DB("").C(aName)

	//then insert app and respond
	if insertErr := c.Insert(app); insertErr != nil {

		s.internalError(r, w, insertErr, "error inserting: "+fmt.Sprintf("%v", app))

	} else {
		response := AppResponse{APIToken: app.APIToken, APISecret: app.APISecret, Name: app.Name, OS: app.OS}
		s.logger.Printf("created new app: %+v", response)
		s.serveJSON(w, &response)
	}

}

//GET "/api/v1/developers/apps/:objectId"
func (s *Server) queryApp(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("QueryApp: ")

	//getObjectId
	objectID := s.getObjectID(w, r)

	//get collection
	session := s.db.GetSession()
	defer session.Close()
	c := session.DB("").C(aName)

	//find and serve data
	app := app{}
	if err := c.FindId(bson.ObjectIdHex(objectID)).Limit(1).One(&app); err != nil {
		s.notFound(r, w, err, objectID+" : id not found")
		return
	}

	//respond
	response := AppResponse{APIToken: app.APIToken, APISecret: app.APISecret, Name: app.Name, OS: app.OS}
	//s.logger.Printf("query app: %+v, id: %v\n", response)
	s.serveJSON(w, &response)

}

//GET "/api/v1/developers/:developerId/apps/"
func (s *Server) queryAllApps(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("QueryAllApps: ")

	//get params
	did := r.URL.Query().Get(":developerId")
	if did == "" {
		s.notFound(r, w, errors.New("app params are invalid"), "val: "+did)
	}
	developerID := decodeToken(did)
	if developerID == "" {
		s.notFound(r, w, errors.New("app params are invalid"), "val: "+developerID)
	}

	//get collection
	session := s.db.GetSession()
	defer session.Close()
	c := session.DB("").C(aName)

	//find apps
	var apps []app
	iter := c.Find(bson.M{"parentId": bson.ObjectIdHex(developerID)}).Iter()
	err := iter.All(&apps)
	if err != nil {
		s.internalError(r, w, err, "error iterating app documents")
	}

	//respond
	var re []AppResponse
	for _, app := range apps {

		re = append(re, AppResponse{APIToken: app.APIToken, APISecret: app.APISecret, Name: app.Name, OS: app.OS})

	}

	s.serveJSON(w, &re)
}

//PUT "/api/v1/developers/apps/:objectId"
func (s *Server) updateApp(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("UpdateApp: ")

	//getObjectId
	objectID := s.getObjectID(w, r)

	//parse body
	updateRequest := &updateAppRequest{}
	if err := s.readJSON(updateRequest, r, w); err != nil {
		s.badRequest(r, w, err, "malformed update request body")
		return
	}

	//get collection
	session := s.db.GetSession()
	defer session.Close()
	c := session.DB("").C(aName)

	//change
	var change = mgo.Change{
		ReturnNew: true,
		Update: bson.M{
			"$set": bson.M{
				"updated": time.Now(),
				"name":    updateRequest.Name,
				"os":      updateRequest.OS,
			}}}

	//find and update
	app := app{}
	if _, err := c.FindId(bson.ObjectIdHex(objectID)).Apply(change, &app); err != nil {
		s.notFound(r, w, err, objectID+" : id not found")
		return
	}

	//respond
	response := AppResponse{APIToken: app.APIToken, APISecret: app.APISecret, Name: app.Name, OS: app.OS}
	s.serveJSON(w, &response)

}

//DELETE "/api/v1/developers/apps/:objectId"
func (s *Server) deleteApp(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("DeleteApp: ")

	//getObjectId
	objectID := s.getObjectID(w, r)

	//get collection
	session := s.db.GetSession()
	defer session.Close()
	c := session.DB("").C(aName)

	//delete
	if err := c.RemoveId(bson.ObjectIdHex(objectID)); err != nil {
		s.notFound(r, w, err, objectID+" : id not found")
		return
	}

	//respond
	response := deleteResponse{Status: "ok"}
	s.serveJSON(w, &response)

}
