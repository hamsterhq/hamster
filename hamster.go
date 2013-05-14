package hamster

/*Hamster Server*/

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/adnaan/routes"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

//http://blog.semantics3.com/building-a-paid-api-offering/
//db/apps/:name/classes/:name
//base path:db/apps/:name-->identified by each handler from headers
//relative path:/classes/:name-->from request params
//class:json from request body
// A server controls access to all classes

//the fcgi server. Not entirely sure why fcgi?. But it sounds cool:/
type Server struct {
	listener   net.Listener
	logger     *log.Logger
	httpServer *http.Server
	route      *routes.RouteMux
	db         *Db
	config     *Config
}

//dbUrl:"mongodb://adnaan:pass@localhost:27017/hamster"
//serverUrl:fmt.Sprintf("%s:%d", address, port)
//creates a new server, setups logging etc.
func NewServer(port int, dbUrl string) *Server {
	f, err := os.OpenFile("hamster.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("hamster.log faied to open")

	}
	//log.SetOutput(f)
	//log.SetOutput(os.Stdout)

	r := routes.New()
	var cfg Config
	if _, err := toml.DecodeFile("hamster.toml", &cfg); err != nil {
		fmt.Println(err)
		return nil
	}
	s := &Server{
		httpServer: &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: r},
		route:      r,
		logger:     log.New(f, "", log.LstdFlags),
		db:         &Db{Url: dbUrl},
		config:     &cfg,
	}

	s.addHandlers()

	return s

}

//listen and serve a fastcgi server

func (s *Server) ListenAndServe() error {

	listener, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		s.logger.Printf("error listening: %v \n", err)
		return err
	}
	s.listener = listener

	go s.httpServer.Serve(s.listener)

	s.logger.Print("********Server Startup*********\n")
	s.logger.Print("********++++++++++++++*********\n")
	s.logger.Printf("hamster is now listening on http://localhost%s\n", s.httpServer.Addr)

	//index the collections
	s.IndexDevelopers()

	return nil
}

// stops the server.
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

// no log
func (s *Server) Quiet() {
	s.logger = log.New(ioutil.Discard, "", log.LstdFlags)
}
