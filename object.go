package hamster

import (
	"errors"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"time"
)

var oName = "objects"

type HamsterObject struct {
	//default fields
	Id             bson.ObjectId `bson:"_id" json:"-"`
	ParentId       bson.ObjectId `bson:"parentId" json:"-"`
	ObjectId       string        `bson:"-" json:"object_id"`
	ParentObjectId string        `bson:"-" json:"parent_id"`
	Created        time.Time     `bson:"created" json:"created"`
	Updated        time.Time     `bson:"updated" json:"updated"`
	//unknown fields
	M map[string]interface{} `bson:",inline" json:"-"`
}

//create new user defined object
func (s *Server) CreateObject(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("CreateObject: ")
	//get objectName
	object_name := s.getObjectName(w, r)
	//get app object id
	app_id := s.getAppObjectId(w, r)
	//verify app object id valid
	//get apps collection
	session := s.db.GetSession()
	defer session.Close()
	a := session.DB("").C(aName)

	//change
	var change = mgo.Change{
		ReturnNew: true,
		Update: bson.M{
			"$set": bson.M{
				"updated": time.Now(),
			},
			"$addToSet": bson.M{
				"objects": object_name,
			}}}

	//if app id exists, update
	app := App{}
	if _, err := a.FindId(bson.ObjectIdHex(app_id)).Apply(change, &app); err != nil {
		s.notFound(r, w, err, app_id+" : id not found")
		return
	}

	//fmt.Printf("app: %v\n", app)

	//parse body
	var obj interface{}
	s.readJson(&obj, r, w)

	//and marshal into hamster object
	hamster_obj := HamsterObject{}

	h, err := bson.Marshal(obj)
	if err != nil {
		s.internalError(r, w, err, "error marshalling hamster object")
	}
	bson.Unmarshal(h, &hamster_obj)

	//set fields
	hamster_obj.Id = bson.NewObjectId()
	hamster_obj.ParentId = bson.ObjectIdHex(app_id)
	hamster_obj.Created = time.Now()
	hamster_obj.Updated = time.Now()

	//get objects collection
	c := session.DB("").C(object_name)

	//then insert object
	if insert_err := c.Insert(hamster_obj); insert_err != nil {

		s.internalError(r, w, insert_err, "error inserting: "+fmt.Sprintf("%v", hamster_obj))

	} else {

		//find inlined object
		var result map[string]interface{}
		if err := c.FindId(hamster_obj.Id).Limit(1).One(&result); err != nil {
			s.notFound(r, w, err, hamster_obj.Id.Hex()+" : id not found")
			return
		}

		//append object_id,parent_id
		delete(result, "_id")
		result["object_id"] = encodeBase64Token(hamster_obj.Id.Hex())
		delete(result, "parentId")
		result["parent_id"] = encodeBase64Token(hamster_obj.ParentId.Hex())

		s.logger.Printf("created new object: %+v, id: %v\n", result)
		s.serveJson(w, &result)
	}

}

func (s *Server) CreateObjects(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("CreateObjects: ")
	//get param
	app_id := s.getAppObjectId(w, r)
	//parse body
	var obj map[string]interface{}
	s.readJson(&obj, r, w)

	if obj["__op"] != "InsertBatch" {
		s.badRequest(r, w, errors.New("expected batch insert op"), fmt.Sprintf("%v", obj["__op"]))
	}

	var object_name string
	object_name = obj["objectName"].(string)
	if obj["objectName"] == "" {
		s.badRequest(r, w, errors.New("expected objectName"), "")
	}

	//get apps collection
	session := s.db.GetSession()
	defer session.Close()
	a := session.DB("").C(aName)

	//change
	var change = mgo.Change{
		ReturnNew: true,
		Update: bson.M{
			"$set": bson.M{
				"updated": time.Now(),
			},
			"$addToSet": bson.M{
				"objects": object_name,
			}}}

	//if app id exists, update
	app := App{}
	if _, err := a.FindId(bson.ObjectIdHex(app_id)).Apply(change, &app); err != nil {
		s.notFound(r, w, err, app_id+" : id not found")
		return
	}

	var objects []map[string]interface{}
	objects = obj["objects"].([]map[string]interface{})

	hamster_objects := []HamsterObject{}

	//marshal incoming objects
	for _, object := range objects {
		hamster_obj := HamsterObject{}
		h, err := bson.Marshal(obj)
		if err != nil {
			s.internalError(r, w, err, "error marshalling hamster object: "+fmt.Sprintf("%v", object))
		}
		bson.Unmarshal(h, &hamster_obj)
		hamster_objects = append(hamster_objects, hamster_obj)

	}

	//set fields
	time_now := time.Now()
	for _, hamster_obj := range hamster_objects {
		hamster_obj.Id = bson.NewObjectId()
		hamster_obj.ParentId = bson.ObjectIdHex(app_id)
		hamster_obj.Created = time_now
		hamster_obj.Updated = time_now

	}

	//get objects collection
	c := session.DB("").C(object_name)

	//then insert object
	if insert_err := c.Insert(hamster_objects); insert_err != nil {

		s.internalError(r, w, insert_err, "error inserting: "+fmt.Sprintf("%v", hamster_obj))

	} else {

		//find inlined object
		var result []map[string]interface{}
		if err := c.FindId(hamster_obj.Id).Limit(1).One(&result); err != nil {
			s.notFound(r, w, err, hamster_obj.Id.Hex()+" : id not found")
			return
		}

		//append object_id,parent_id
		delete(result, "_id")
		result["object_id"] = encodeBase64Token(hamster_obj.Id.Hex())
		delete(result, "parentId")
		result["parent_id"] = encodeBase64Token(hamster_obj.ParentId.Hex())

		s.logger.Printf("created new object: %+v, id: %v\n", result)
		s.serveJson(w, &result)
	}

}

//query object
func (s *Server) QueryObject(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("QueryObject: ")
	object_name, object_id := s.getObjectParams(w, r)

	//get object collection
	session := s.db.GetSession()
	defer session.Close()
	c := session.DB("").C(object_name)

	//find object
	var result map[string]interface{}
	if err := c.FindId(bson.ObjectIdHex(object_id)).Limit(1).One(&result); err != nil {
		s.notFound(r, w, err, object_id+" : id not found")
		return
	}

	//append object_id,parent_id
	//convert object id to base64
	result["object_id"] = encodeBase64Token(result["_id"].(bson.ObjectId).Hex())
	delete(result, "_id")
	result["parent_id"] = encodeBase64Token(result["parentId"].(bson.ObjectId).Hex())
	delete(result, "parentId")

	s.serveJson(w, &result)

}

//query objects
func (s *Server) QueryObjects(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("QueryObjects: ")

	object_name := s.getObjectName(w, r)

	//get collection
	session := s.db.GetSession()
	defer session.Close()
	c := session.DB("").C(object_name)

	//find apps
	var result []map[string]interface{}
	iter := c.Find(nil).Iter()
	err := iter.All(&result)
	if err != nil {
		s.internalError(r, w, err, "error iterating "+object_name+" documents")
	}

	//convert object id to base64
	for _, object := range result {
		object["object_id"] = encodeBase64Token(object["_id"].(bson.ObjectId).Hex())
		delete(object, "_id")
		object["parent_id"] = encodeBase64Token(object["parentId"].(bson.ObjectId).Hex())
		delete(object, "parentId")

	}

	s.serveJson(w, &result)

}

//update object
func (s *Server) UpdateObject(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("UpdateObject: ")
	object_name, object_id := s.getObjectParams(w, r)

	//parse body
	var u map[string]interface{}
	if err := s.readJson(&u, r, w); err != nil {
		s.badRequest(r, w, err, "malformed update request body")
		return
	}

	//add update field
	u["updated"] = time.Now()

	//get object collection
	session := s.db.GetSession()
	defer session.Close()
	c := session.DB("").C(object_name)

	//change
	var change = mgo.Change{
		ReturnNew: true,
		Update: bson.M{
			"$set": u}}

	//find and update
	var result map[string]interface{}
	if _, err := c.FindId(bson.ObjectIdHex(object_id)).Apply(change, &result); err != nil {
		s.notFound(r, w, err, object_id+" : id not found")
		return
	}

	//append object_id,parent_id
	//convert object id to base64
	result["object_id"] = encodeBase64Token(result["_id"].(bson.ObjectId).Hex())
	delete(result, "_id")
	result["parent_id"] = encodeBase64Token(result["parentId"].(bson.ObjectId).Hex())
	delete(result, "parentId")

	s.serveJson(w, &result)

}

//delete object
func (s *Server) DeleteObject(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("DeleteObject: ")

	//get params
	object_name, object_id := s.getObjectParams(w, r)

	//get object collection
	session := s.db.GetSession()
	defer session.Close()
	c := session.DB("").C(object_name)

	//delete
	if err := c.RemoveId(bson.ObjectIdHex(object_id)); err != nil {
		s.notFound(r, w, err, object_id+" : id not found")
		return
	}

	//respond
	response := DeleteResponse{Status: "ok"}
	s.serveJson(w, &response)

}
