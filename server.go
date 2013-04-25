package hamster

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"hash/fnv"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"regexp"
	"runtime"
	"time"
)

// A server controls access to all objects
type Server struct {
	logger     *log.Logger
	httpServer *http.Server
	path       string
	router     *mux.Router
	objects    map[string]*Object
}
