package core

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"hash"
	"net/http"
	"strings"
	"sync"
	"time"

	"code.google.com/p/go.crypto/bcrypt"
	"github.com/garyburd/redigo/redis"
	"github.com/kr/fernet"
	"labix.org/v2/mgo/bson"
)

//Allow only verified ip's from config
func (s *Server) baseAuth(w http.ResponseWriter, r *http.Request) {
	//ip := strings.Split(r.RemoteAddr, ":")
	//fmt.Println(r.RemoteAddr)
	//s.logger.SetPrefix("BaseAuth:")
	//s.logger.Printf("%v  %v", r.Method, r.URL.Path)

	//TODO:IPV6 local loop back address resolves to ::1. Find a better solution
	// if !s.ipAllowed(ip[0]) {
	// 	http.Error(w, fmt.Sprintf("Unauthorized ip:%s not allowed", ip[0]), http.StatusUnauthorized)
	// 	return
	//
	// }

}

//Allow ip
func (s *Server) ipAllowed(IP string) bool {
	for _, client := range s.config.Clients {
		if IP == client.IP {
			return true
		}
	}
	return false
}

//Authenticate developer. Access Token is generated using a shared secret between
//Hamster and the client. The shared secret is manually configured in hamster.toml.
//TODO: find a better implementation
func (s *Server) developerAuth(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("DeveloperAuth:")

	//skip access token check if logging in
	if r.Method == "POST" && r.URL.Path == "/api/v1/developers/login/" {

		return
	}

	accessToken := r.Header.Get("X-Access-Token")

	if accessToken == "" {

		s.unauthorized(r, w, errors.New("token is empty"), "access token invalid")
		return

	}

	//check POST /api/v1/developers
	if r.Method == "POST" && r.URL.Path == "/api/v1/developers/" {

		//check auth
		if !s.validateSharedToken(accessToken) {
			s.unauthorized(r, w, errors.New("shared token is old"), "access token invalid")
		}
		return
	}

	if _, ok := s.validateAccessLoginToken(accessToken); !ok {
		s.unauthorized(r, w, errors.New("token is old"), "access token invalid")
		return

	}
}

//Validate access token
func (s *Server) validateAccessLoginToken(token string) (string, bool) {

	btok, err := base64.URLEncoding.DecodeString(token)

	if err != nil {
		return "", false

	}
	k := fernet.MustDecodeKeys(s.config.Clients["browser"].Secret)
	email := fernet.VerifyAndDecrypt(btok, 60*time.Second, k)

	c := s.redisConn()
	defer c.Close()

	status, err := redis.String(c.Do("GET", email))
	if err != nil {
		return "", false
	}

	if status == "loggedin" {
		return string(email), true
	}

	return "", false
}

//Validate access token
func (s *Server) validateSharedToken(token string) bool {

	btok, err := base64.URLEncoding.DecodeString(token)

	if err != nil {
		return false

	}
	k := fernet.MustDecodeKeys(s.config.Clients["browser"].Secret)
	sharedToken := fernet.VerifyAndDecrypt(btok, 60*10*time.Second, k)
	if string(sharedToken) == string(s.config.Clients["browser"].Token) {
		return true
	}

	return false

}

//Authenticates object level requests
func (s *Server) objectAuth(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("ObjectAuth:")
	apiToken := r.Header.Get("X-Api-Token")
	apiSecret := r.Header.Get("X-Api-Secret")

	if apiToken == "" || apiSecret == "" {
		s.unauthorized(r, w, errors.New("token or secret invalid"), "access token invalid")
		return
	}

	if !s.validateAPIToken(apiToken, apiSecret) {
		s.unauthorized(r, w, errors.New("token match failed"), "access token failed")
		return
	}

}

func (s *Server) validateAPIToken(token string, secret string) bool {

	if ok, hash, salt := s.getHashSalt(token); ok {
		//fmt.Println("found key in redis")
		if matchPassword(decodeToken(secret), hash, salt) {
			return true
		}
	}
	if ok, hash, salt := s.getHashSaltFromDB(token); ok {
		if matchPassword(decodeToken(secret), hash, salt) {
			return true
		}
	}
	return false
}

func (s *Server) getHashSalt(token string) (bool, string, string) {
	c := s.redisConn()
	defer c.Close()

	hash, err := redis.String(c.Do("GET", token+":hash"))
	if err != nil {
		return false, "", ""
	}
	salt, err := redis.String(c.Do("GET", token+":salt"))
	if err != nil {
		return false, "", ""
	}

	return true, hash, salt
}

func (s *Server) getHashSaltFromDB(token string) (bool, string, string) {
	//get collection

	appID := decodeToken(token)
	session := s.db.GetSession()
	defer session.Close()
	c := session.DB("").C(aName)

	app := app{}
	//TODO:select fields
	if err := c.FindId(bson.ObjectIdHex(appID)).One(&app); err != nil {
		//fmt.Println("app not found\n")
		return false, "", ""
	}

	//cache hash-salt in redis
	rc := s.redisConn()
	defer rc.Close()

	rc.Do("SET", token+":hash", app.Hash)
	rc.Do("SET", token+":salt", app.Salt)

	return true, app.Hash, app.Salt

}

//Get Basic user password
func getUserPassword(r *http.Request) (string, string) {

	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(s) != 2 || s[0] != "Basic" {
		return "", ""
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return "", ""
	}
	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		return "", ""
	}

	return pair[0], pair[1]

}

type hmacFactory struct {
	hashFunc hash.Hash
	m        sync.Mutex
}

/*Generates hash*/
func (h *hmacFactory) generateHash(data []byte) []byte {
	h.m.Lock()
	defer h.m.Unlock()

	h.hashFunc.Reset()
	h.hashFunc.Write(data)
	return h.hashFunc.Sum(nil)

}

/*Encrypts the password and returns password hash*/
func (h *hmacFactory) encrypt(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(h.generateHash(password), cost)
}

func (h *hmacFactory) compare(hash, password []byte) error {
	return bcrypt.CompareHashAndPassword(hash, h.generateHash(password))
}

/*Returns a hmac type*/
func newHMAC(hash func() hash.Hash, salt []byte) *hmacFactory {
	hm := &hmacFactory{
		hashFunc: hmac.New(hash, salt),
		m:        sync.Mutex{},
	}
	return hm
}

/*Encrypts password. Returns hash+salt*/
func encryptPassword(password string) (string, string, error) {

	salt, err0 := genUUID(16)
	if err0 != nil {
		return "", "", err0
	}

	hm := newHMAC(sha512.New, []byte(salt))
	pass := []byte(password)
	encrypted, err := hm.encrypt(pass, bcrypt.DefaultCost)

	if err != nil {

		return "", "", err
	}

	return string(encrypted), salt, nil
}

/*Match encrypted string*/
func matchPassword(password string, hash string, salt string) bool {
	p := []byte(password)
	h := []byte(hash)
	s := []byte(salt)
	hm := newHMAC(sha512.New, s)
	err := hm.compare(h, p)
	if err != nil {
		return false
	}
	return true
}

/*Encode to base64*/
func encodeBase64Token(hexVal string) string {

	token := base64.URLEncoding.EncodeToString([]byte(hexVal))

	return token

}

/*Decode from base64*/
func decodeToken(token string) string {

	hexVal, err := base64.URLEncoding.DecodeString(token)
	if err != nil {

		return ""

	}

	return string(hexVal)

}

/*Generate uuid*/
func genUUID(size int) (string, error) {
	uuid := make([]byte, size)
	n, err := rand.Read(uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}

	uuid[8] = 0x80
	uuid[4] = 0x40

	return hex.EncodeToString(uuid), nil
}

/*
 Return a random 16-byte base64 alphabet string
*/
func randomKey() string {
	k := make([]byte, 12)
	for bytes := 0; bytes < len(k); {
		n, err := rand.Read(k[bytes:])
		if err != nil {
			panic("rand.Read() failed")
		}
		bytes += n
	}
	return base64.StdEncoding.EncodeToString(k)
}

//Generate time based access token using shared secret. See fernet project
//for more details
func (s *Server) genAccessToken(email string) (string, error) {
	//encrypt token
	k := fernet.MustDecodeKeys(s.config.Clients["browser"].Secret)
	tok, err := fernet.EncryptAndSign([]byte(email), k[0])
	if err != nil {

	}
	token := base64.URLEncoding.EncodeToString(tok)

	//cache it
	c := s.redisConn()
	defer c.Close()

	c.Do("SET", email, "loggedin")

	return token, nil

}

//Logout developer
//TODO: find a better way to handle login/logout
func (s *Server) logout(email string) error {

	c := s.redisConn()
	defer c.Close()

	status, err := redis.String(c.Do("GET", email))
	if err != nil {
		return err
	}

	if status == "loggedin" {
		c.Do("SET", email, "") //log out
	} else {
		return errors.New("user is not logged in")
	}

	return nil

}
