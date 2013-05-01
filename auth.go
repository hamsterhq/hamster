package hamster

//get application id and key for incoming request
import (
	"code.google.com/p/go.crypto/bcrypt"
	"crypto/hmac"
	//"crypto/sha512"
	"hash"
	"sync"
)

//generate salt
/*salt := make([]byte, 64)
if n, err := io.ReadFull(rand.Reader, salt); err != nil {

}*/

var (
	Salt = []byte("putyoursalt/here")
)

func (s *Server) GetAppIdKey(app_name string) (app_id string, app_key string) {

	return "", ""

}

//https://github.com/abbot/go-http-auth
//http://security.stackexchange.com/questions/19809/how-should-api-keys-be-generated
func (s *Server) GenerateAppIdKey() (app_id string, app_key string) {

	return "", ""

}

type Hmac struct {
	hashFunc hash.Hash
	m        sync.Mutex
}

func (h *Hmac) generateHash(data []byte) []byte {
	h.m.Lock()
	defer h.m.Unlock()

	h.hashFunc.Reset()
	h.hashFunc.Write(data)
	return h.hashFunc.Sum(nil)

}

func (h *Hmac) Encrypt(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(h.generateHash(password), cost)
}

func (h *Hmac) Compare(hash, password []byte) error {
	return bcrypt.CompareHashAndPassword(hash, h.generateHash(password))
}

func New(hash func() hash.Hash, salt []byte) *Hmac {
	hm := &Hmac{
		hashFunc: hmac.New(hash, salt),
		m:        sync.Mutex{},
	}
	return hm
}

func enrypt(password string, hm Hmac) (string, error) {

	//	hm := New(sha512.New, Salt)
	pass := []byte(password)
	digest, err := hm.Encrypt(pass, bcrypt.DefaultCost)

	if err != nil {

		return "", err
	}

	return string(digest), nil
}
