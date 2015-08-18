package core

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

var oName = "objects"

type hamsterObject struct {
	//default fields
	ID             bson.ObjectId `bson:"_id" json:"-"`
	ParentID       bson.ObjectId `bson:"parentId" json:"-"`
	ObjectID       string        `bson:"-" json:"objectID"`
	ParentObjectID string        `bson:"-" json:"parent_id"`
	Created        time.Time     `bson:"created" json:"created"`
	Updated        time.Time     `bson:"updated" json:"updated"`
	//unknown fields
	M map[string]interface{} `bson:",inline" json:"-"`
}

//create new user defined object
//POST:/api/v1/objects/:objectName  handler
func (s *Server) createObject(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("CreateObject: ")
	//get objectName
	objectName := s.getObjectName(w, r)
	//get app object id
	appID := s.getAppObjectID(w, r)
	//verify app object id valid
	//get apps collection
	session := s.db.GetSession()
	defer session.Close()
	a := session.DB("").C(aName)

	//change:update new object type name in app if it doesn't exist
	var change = mgo.Change{
		ReturnNew: true,
		Update: bson.M{
			"$set": bson.M{
				"updated": time.Now(),
			},
			"$addToSet": bson.M{
				"objects": objectName,
			}}}

	//if app id exists, update
	app := app{}
	if _, err := a.FindId(bson.ObjectIdHex(appID)).Apply(change, &app); err != nil {
		s.notFound(r, w, err, appID+" : id not found")
		return
	}

	//fmt.Printf("app: %v\n", app)

	//parse body
	var obj interface{}
	s.readJSON(&obj, r, w)

	//and marshal into hamster object
	hamsterObj := hamsterObject{}

	h, err := bson.Marshal(obj)
	if err != nil {
		s.internalError(r, w, err, "error marshalling hamster object")
	}
	bson.Unmarshal(h, &hamsterObj)

	//set fields
	hamsterObj.ID = bson.NewObjectId()
	hamsterObj.ParentID = bson.ObjectIdHex(appID)
	hamsterObj.Created = time.Now()
	hamsterObj.Updated = time.Now()

	//get objects collection
	c := session.DB("").C(objectName)

	//then insert object
	if insertErr := c.Insert(hamsterObj); insertErr != nil {

		s.internalError(r, w, insertErr, "error inserting: "+fmt.Sprintf("%v", hamsterObj))

	} else {

		//find inlined object
		var result map[string]interface{}
		if err := c.FindId(hamsterObj.ID).Limit(1).One(&result); err != nil {
			s.notFound(r, w, err, hamsterObj.ID.Hex()+" : id not found")
			return
		}

		//append objectID,parent_id
		delete(result, "_id")
		result["objectID"] = encodeBase64Token(hamsterObj.ID.Hex())
		delete(result, "parentId")
		result["parent_id"] = encodeBase64Token(hamsterObj.ParentID.Hex())

		s.logger.Printf("created new object: %+v\n", result)
		s.serveJSON(w, &result)
	}

}

//create multiple user defined object
//POST:/api/v1/objects/ handler
func (s *Server) createObjects(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("CreateObjects: ")
	//get param
	appID := s.getAppObjectID(w, r)
	//parse body
	var obj map[string]interface{}
	s.readJSON(&obj, r, w)

	if obj["__op"] != "InsertBatch" {
		s.badRequest(r, w, errors.New("expected batch insert op"), fmt.Sprintf("%v", obj["__op"]))
	}

	var objectName string
	objectName = obj["objectName"].(string)
	fmt.Printf("object name: %v \n", obj["objectName"])
	if obj["objectName"] == "" {
		s.badRequest(r, w, errors.New("expected objectName"), "")
	}

	//get apps collection
	session := s.db.GetSession()
	defer session.Close()
	a := session.DB("").C(aName)

	//change:update new object type name in app if it doesn't exist
	var change = mgo.Change{
		ReturnNew: true,
		Update: bson.M{
			"$set": bson.M{
				"updated": time.Now(),
			},
			"$addToSet": bson.M{
				"objects": objectName,
			}}}

	//if app id exists, update
	app := app{}
	if _, err := a.FindId(bson.ObjectIdHex(appID)).Apply(change, &app); err != nil {
		s.notFound(r, w, err, appID+" : id not found")
		return
	}

	var objects []map[string]interface{}
	objects = obj["objects"].([]map[string]interface{})

	hamsterObjects := []hamsterObject{}

	//marshal incoming objects
	for _, object := range objects {
		hamsterObj := hamsterObject{}
		h, err := bson.Marshal(obj)
		if err != nil {
			s.internalError(r, w, err, "error marshalling hamster object: "+fmt.Sprintf("%v", object))
		}
		bson.Unmarshal(h, &hamsterObj)
		hamsterObjects = append(hamsterObjects, hamsterObj)

	}

	//set fields
	timeNow := time.Now()
	for _, hamsterObj := range hamsterObjects {
		hamsterObj.ID = bson.NewObjectId()
		hamsterObj.ParentID = bson.ObjectIdHex(appID)
		hamsterObj.Created = timeNow
		hamsterObj.Updated = timeNow

	}

	//get objects collection
	c := session.DB("").C(objectName)

	//then insert object
	if insertErr := c.Insert(hamsterObjects); insertErr != nil {

		s.internalError(r, w, insertErr, "error inserting: "+fmt.Sprintf("%v", hamsterObjects))

	} else {

		//find inlined objects
		var results []map[string]interface{}
		if err := c.Find(bson.M{"created": bson.M{"$gte": timeNow, "$lt": time.Now()}}).All(&results); err != nil {
			s.notFound(r, w, err, " : objects not found")
			return
		}

		for _, result := range results {
			//append objectID,parent_id
			//convert object id to base64
			result["objectID"] = encodeBase64Token(result["_id"].(bson.ObjectId).Hex())
			delete(result, "_id")
			result["parent_id"] = encodeBase64Token(result["parentId"].(bson.ObjectId).Hex())
			delete(result, "parentId")
		}

		s.logger.Printf("created new objects: %+v", results)
		s.serveJSON(w, &results)
	}

}

//GET: /api/v1/objects/:objectName/:objectId handler
func (s *Server) queryObject(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("QueryObject: ")
	objectName, objectID := s.getObjectParams(w, r)

	//get object collection
	session := s.db.GetSession()
	defer session.Close()
	c := session.DB("").C(objectName)

	//find object
	var result map[string]interface{}
	if err := c.FindId(bson.ObjectIdHex(objectID)).Limit(1).One(&result); err != nil {
		s.notFound(r, w, err, objectID+" : id not found")
		return
	}

	//append objectID,parent_id
	//convert object id to base64
	result["objectID"] = encodeBase64Token(result["_id"].(bson.ObjectId).Hex())
	delete(result, "_id")
	result["parent_id"] = encodeBase64Token(result["parentId"].(bson.ObjectId).Hex())
	delete(result, "parentId")

	s.serveJSON(w, &result)

}

//GET: /api/v1/objects/:objectName handler
func (s *Server) queryObjects(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("QueryObjects: ")

	objectName := s.getObjectName(w, r)

	//get collection
	session := s.db.GetSession()
	defer session.Close()
	c := session.DB("").C(objectName)

	//find apps
	var result []map[string]interface{}
	iter := c.Find(nil).Iter()
	err := iter.All(&result)
	if err != nil {
		s.internalError(r, w, err, "error iterating "+objectName+" documents")
	}

	//convert object id to base64
	for _, object := range result {
		object["objectID"] = encodeBase64Token(object["_id"].(bson.ObjectId).Hex())
		delete(object, "_id")
		object["parent_id"] = encodeBase64Token(object["parentId"].(bson.ObjectId).Hex())
		delete(object, "parentId")

	}

	s.serveJSON(w, &result)

}

//PUT: /api/v1/objects/:objectName/:objectId
func (s *Server) updateObject(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("UpdateObject: ")
	objectName, objectID := s.getObjectParams(w, r)

	//parse body
	var u map[string]interface{}
	if err := s.readJSON(&u, r, w); err != nil {
		s.badRequest(r, w, err, "malformed update request body")
		return
	}

	//add update field
	u["updated"] = time.Now()

	//get object collection
	session := s.db.GetSession()
	defer session.Close()
	c := session.DB("").C(objectName)

	//change
	var change = mgo.Change{
		ReturnNew: true,
		Update: bson.M{
			"$set": u}}

	//find and update
	var result map[string]interface{}
	if _, err := c.FindId(bson.ObjectIdHex(objectID)).Apply(change, &result); err != nil {
		s.notFound(r, w, err, objectID+" : id not found")
		return
	}

	//append objectID,parent_id
	//convert object id to base64
	result["objectID"] = encodeBase64Token(result["_id"].(bson.ObjectId).Hex())
	delete(result, "_id")
	result["parent_id"] = encodeBase64Token(result["parentId"].(bson.ObjectId).Hex())
	delete(result, "parentId")

	s.serveJSON(w, &result)

}

//DELETE:/api/v1/objects/:objectName/:objectId
func (s *Server) deleteObject(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("DeleteObject: ")

	//get params
	objectName, objectID := s.getObjectParams(w, r)

	//get object collection
	session := s.db.GetSession()
	defer session.Close()
	c := session.DB("").C(objectName)

	//delete
	if err := c.RemoveId(bson.ObjectIdHex(objectID)); err != nil {
		s.notFound(r, w, err, objectID+" : id not found")
		return
	}

	//respond
	response := deleteResponse{Status: "ok"}
	s.serveJSON(w, &response)

}
