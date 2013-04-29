package hamster

import (
	"net/http"
	"net/http/pprof"
)

/*----------------------------------------------------------*/
/*Handlers*/
/*
/1/classes/<className>	POST	Creating Objects
/1/classes/<className>/<objectId>	GET	Retrieving Objects
/1/classes/<className>/<objectId>	PUT	Updating Objects
/1/classes/<className>	GET	Queries
/1/classes/<className>/<objectId>	DELETE	Deleting Objects
*/

//Authenticate rest api calls with application-id and application-key
var Auth = func(w http.ResponseWriter, r *http.Request) {
	/*if r.URL.User == nil || r.URL.User.Username() != "admin" {
		http.Error(w, "", http.StatusUnauthorized)
	}*/
}

//add path handlers
func (s *Server) addHandlers() {

	//authentication
	s.router.Filter(Auth)

	//add pprof path handlers
	s.router.AddRoute("GET", "/debug/pprof", pprof.Index)
	s.router.AddRoute("GET", "/debug/pprof/cmdline", pprof.Cmdline)
	s.router.AddRoute("GET", "/debug/pprof/profile", pprof.Profile)
	s.router.AddRoute("GET", "/debug/pprof/symbol", pprof.Symbol)

	//create a developer.
	s.router.Post("/developers", s.CreateDev)
	//get a developer
	s.router.Get("/developers/:objectId", s.GetDev)
	//login
	s.router.Get("/developers/:objectId/login", s.LoginDev)
	//update developer
	s.router.Put("/developers/:objectId", s.UpdateDev)
	//queries
	s.router.Get("/developers", s.QueryDev)
	//delete object
	s.router.Del("/developers:objectId", s.DeleteDev)

	//create an app.
	s.router.Post("/apps/", s.CreateApp)
	//get an app
	s.router.Get("/apps/:objectId", s.GetApp)
	//update app
	s.router.Put("/apps/:objectId", s.UpdateApp)
	//queries
	s.router.Get("/apps", s.QueryApp)
	//delete app
	s.router.Del("/apps:objectId", s.DeleteApp)

	//create an object. create collection(class) if it does not exist else add document(object)
	s.router.Post("/classes/:className", s.CreateObject)
	//get an object
	s.router.Get("/classes/:className/:objectId", s.GetObject)
	//update object
	s.router.Put("/classes/:className/:objectId", s.UpdateObject)
	//queries
	s.router.Get("/classes/:className", s.QueryObject)
	//delete object
	s.router.Del("/classes/:className/:objectId", s.DeleteObject)

}
