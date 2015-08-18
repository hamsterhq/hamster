/*
The Hamster Server. The Server type holds instances of all the components,
*effectively making it possible to collapse all the code into one file. The separation
* of code is only for readability. To use it as a package simply:
* import ("github.com/adnaan/hamster")
* server := hamster.NewServer()
* //server.Quiet()//disable logging
* server.ListenAndServe()
* Also change hamster.toml for custom configuration.
* TODO: Pass hamster.toml as argument to the server
* TODO: make handler methods local, model method exported for pkg/rpc support
*/
package core

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/adnaan/routes"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/sessions"
)

//Server The server type holds instances of all components
type Server struct {
	listener   net.Listener
	logger     *log.Logger
	httpServer *http.Server
	route      *routes.RouteMux
	db         *db
	config     *config
	cookie     *sessions.CookieStore //unused
	redisConn  func() redis.Conn
}

//NewServer creates a new server
//dbUrl:"mongodb://adnaan:pass@localhost:27017/hamster"
//db.addUser( { user: "adnaan",pwd: "pass",roles: [ "readWrite" ] } )
//serverUrl:fmt.Sprintf("%s:%d", address, port)
//creates a new server, setups logging etc.
func NewServer(configPath string) *Server {
	f, err := os.OpenFile("hamster.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("hamster.log faied to open")

	}
	//log.SetOutput(f)
	//log.SetOutput(os.Stdout)
	//router
	r := routes.New()
	//toml config
	var cfg config
	if _, err := toml.DecodeFile("hamster.toml", &cfg); err != nil {
		fmt.Println(err)
		return nil
	}
	//cookie store
	ck := sessions.NewCookieStore([]byte(cfg.Servers["local"].CookieSecret))

	//redis
	var getRedis = func() redis.Conn {

		c, err := redis.Dial("tcp", os.Getenv("REDIS_URL"))
		if err != nil {
			panic(err)
		}

		return c

	}

	//initialize server
	s := &Server{
		httpServer: &http.Server{Addr: ":" + os.Getenv("SERVER_PORT"), Handler: r},
		route:      r,
		logger:     log.New(f, "", log.LstdFlags),
		db:         &db{URL: os.Getenv("MONGODB_URL")},
		config:     &cfg,
		cookie:     ck,
		redisConn:  getRedis,
	}

	s.logger.SetFlags(log.Lshortfile)
	s.addHandlers()

	return s

}

//ListenAndServe: listen and serve a fastcgi server
func (s *Server) ListenAndServe() error {

	listener, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		s.logger.Printf("error listening: %v \n", err)
		return err
	}
	s.listener = listener

	go s.httpServer.Serve(s.listener)

	//s.logger.Print("********Server Startup*********\n")
	//s.logger.Print("********++++++++++++++*********\n")
	//s.logger.Printf("hamster is now listening on http://localhost%s\n", s.httpServer.Addr)

	//index the collections
	s.indexDevelopers()

	return nil
}

// Shutdown the server.
func (s *Server) Shutdown() error {

	if s.listener != nil {
		// Then stop the server.
		err := s.listener.Close()
		s.listener = nil
		if err != nil {
			return err
		}
	}

	return nil
}

// Quiet down the log
func (s *Server) Quiet() {
	s.logger = log.New(ioutil.Discard, "", log.LstdFlags)
}
