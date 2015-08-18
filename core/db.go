package core

import (
	"log"

	"labix.org/v2/mgo"
)

//mongodb
type db struct {
	URL        string
	MgoSession *mgo.Session
}

//get new db session
func (d *db) GetSession() *mgo.Session {
	if d.MgoSession == nil {
		var err error
		d.MgoSession, err = mgo.Dial(d.URL)
		if err != nil {
			log.Fatalf("dialing mongo url %v failed with %v", d.URL, err)
			//panic(err)
		}
	}
	return d.MgoSession.Clone()

}
