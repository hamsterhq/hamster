package hamster

//get application id and key for incoming request
import (
	"code.google.com/p/go.crypto/bcrypt"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"hash"
	"sync"
)

type Hmac struct {
	hashFunc hash.Hash
	m        sync.Mutex
}

/*Generates hash*/
func (h *Hmac) generateHash(data []byte) []byte {
	h.m.Lock()
	defer h.m.Unlock()

	h.hashFunc.Reset()
	h.hashFunc.Write(data)
	return h.hashFunc.Sum(nil)

}

/*Encrypts the password and returns password hash*/
func (h *Hmac) Encrypt(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(h.generateHash(password), cost)
}

func (h *Hmac) Compare(hash, password []byte) error {
	return bcrypt.CompareHashAndPassword(hash, h.generateHash(password))
}

/*Returns a hmac type*/
func New(hash func() hash.Hash, salt []byte) *Hmac {
	hm := &Hmac{
		hashFunc: hmac.New(hash, salt),
		m:        sync.Mutex{},
	}
	return hm
}

/*Encrypts password. Returns hash+salt*/
func encryptPassword(password string) (string, string, error) {

	salt, err0 := GenUUID(16)
	if err0 != nil {
		return "", "", err0
	}

	hm := New(sha512.New, []byte(salt))
	pass := []byte(password)
	encrypted, err := hm.Encrypt(pass, bcrypt.DefaultCost)

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
	hm := New(sha512.New, s)
	err := hm.Compare(h, p)
	if err != nil {
		return false
	} else {
		//matched!
		return true
	}

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
func GenUUID(size int) (string, error) {
	uuid := make([]byte, size)
	n, err := rand.Read(uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}

	uuid[8] = 0x80
	uuid[4] = 0x40

	return hex.EncodeToString(uuid), nil
}
