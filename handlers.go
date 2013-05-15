package hamster

import (
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

//add path handlers
func (s *Server) addHandlers() {

	//account handlers
	//authentication
	s.route.FilterPrefixPath("/developers", DevAuth)

	//add pprof path handlers
	s.route.AddRoute("GET", "/debug/pprof", pprof.Index)
	s.route.AddRoute("GET", "/debug/pprof/cmdline", pprof.Cmdline)
	s.route.AddRoute("GET", "/debug/pprof/profile", pprof.Profile)
	s.route.AddRoute("GET", "/debug/pprof/symbol", pprof.Symbol)

	//create a developer.
	s.route.Post("/developers/v1", s.CreateDev)
	//get a developer objectId, email or username
	s.route.Get("/developers/v1/login", s.LoginDev)
	//login
	s.route.Get("/developers/v1/logout", s.LogoutDev)
	//update developer
	s.route.Put("/developers/v1/:objectId", s.UpdateDev)
	//queries
	s.route.Get("/developers/v1", s.QueryDev)
	//delete object
	s.route.Del("/developers/v1/:objectId", s.DeleteDev)

	//create an app.
	s.route.Post("/apps/v1", s.CreateApp)
	//get an app
	s.route.Get("/apps/v1/:objectId", s.GetApp)
	//update app
	s.route.Put("/apps/v1/:objectId", s.UpdateApp)
	//queries
	s.route.Get("/apps/v1", s.QueryApp)
	//delete app
	s.route.Del("/apps/v1/:objectId", s.DeleteApp)

	//create an object. create collection(class) if it does not exist else add document(object)
	s.route.Post("/api/v1/classes/:className", s.CreateObject)
	//get an object
	s.route.Get("/api/v1/classes/:className/:objectId", s.GetObject)
	//update object
	s.route.Put("/api/v1/classes/:className/:objectId", s.UpdateObject)
	//queries
	s.route.Get("/api/v1/classes/:className", s.QueryObject)
	//delete object
	s.route.Del("/api/v1/classes/:className/:objectId", s.DeleteObject)

	//files

	//s.rouTe.Post("/files/<filename>", s.SaveFile)
	//s.rouTe.Get("/files/<filename>", s.GetFile)

}
