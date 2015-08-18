package core

import (
	"net/http/pprof"
)

//add path handlers
func (s *Server) addHandlers() {

	//api server info
	s.route.Get("/", s.info)
	//add pprof path handlers
	s.route.AddRoute("GET", "/debug/pprof", pprof.Index)
	s.route.AddRoute("GET", "/debug/pprof/cmdline", pprof.Cmdline)
	s.route.AddRoute("GET", "/debug/pprof/profile", pprof.Profile)
	s.route.AddRoute("GET", "/debug/pprof/symbol", pprof.Symbol)

	//Route filters
	s.route.FilterPrefixPath("/", s.baseAuth)
	s.route.FilterPrefixPath("/api/v1/developers/", s.developerAuth)
	s.route.FilterPrefixPath("/api/v1/objects/", s.objectAuth)
	s.route.FilterPrefixPath("/api/v1/files/", s.objectAuth)

	/*Developer*/
	s.route.Post("/api/v1/developers/", s.createDev)
	//get a developer objectId, email or username
	s.route.Post("/api/v1/developers/login/", s.loginDev)
	//login
	s.route.Post("/api/v1/developers/logout/", s.logoutDev)
	//update developer
	s.route.Put("/api/v1/developers/:objectId", s.updateDev)
	//queries
	s.route.Get("/api/v1/developers/:objectId", s.queryDev)
	//delete object
	s.route.Del("/api/v1/developers/:objectId", s.deleteDev)

	/*App*/
	s.route.Post("/api/v1/developers/:developerId/apps/", s.createApp)
	//get an app
	s.route.Get("/api/v1/developers/apps/:objectId", s.queryApp)
	//queries
	s.route.Get("/api/v1/developers/:developerId/apps/", s.queryAllApps)
	//update app
	s.route.Put("/api/v1/developers/apps/:objectId", s.updateApp)
	//delete app
	s.route.Del("/api/v1/developers/apps/:objectId", s.deleteApp)

	/*Object*/
	s.route.Post("/api/v1/objects/batch/:objectName", s.createObjects)
	s.route.Post("/api/v1/objects/:objectName", s.createObject)
	//get an object
	s.route.Get("/api/v1/objects/:objectName/:objectId", s.queryObject)
	//queries
	s.route.Get("/api/v1/objects/:objectName", s.queryObjects)
	//update object
	s.route.Put("/api/v1/objects/:objectName/:objectId", s.updateObject)
	//delete object
	s.route.Del("/api/v1/objects/:objectName/:objectId", s.deleteObject)

	/*File*/
	s.route.Post("/api/v1/files/:fileName", s.saveFile)
	s.route.Get("/api/v1/files/:fileName/:fileId", s.getFile)

}
