/*Generate private key, write to pem,
load from pem*/

package key

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

func GenPrivateKeyPair(bits int) *rsa.PrivateKey {
	prikey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return prikey
}

func PrivateKeyFromPEM(prikeyPEM string) *rsa.PrivateKey {
	prikeyMarshaled, _ := pem.Decode([]byte(prikeyPEM))
	prikey, err := x509.ParsePKCS1PrivateKey(prikeyMarshaled.Bytes)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	//myprikey.Precompute()
	return prikey
}

func PrivateKeyToPEM(prikey *rsa.PrivateKey) string {
	privateKeyMarshaled := x509.MarshalPKCS1PrivateKey(prikey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Headers: nil, Bytes: privateKeyMarshaled})
	return string(privateKeyPEM)
}

func WritePrivateKeyToPEMFile(prikey *rsa.PrivateKey, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	err = file.Chmod(os.FileMode(0600))
	_, err = file.WriteString(PrivateKeyToPEM(prikey))
	if err != nil {
		return err
	}
	return nil
}

func PrivateKeyFromPEMFile(filename string) (privateKey *rsa.PrivateKey, err error) {
	prikeyFile, err := os.Open(filename)
	if err != nil {
		return
	}
	prikeyFileStat, err := prikeyFile.Stat()
	if err != nil {
		return
	}
	prikeyValue := make([]byte, prikeyFileStat.Size())
	n, err := prikeyFile.Read(prikeyValue)
	if err != nil {
		return
	}
	if n != int(prikeyFileStat.Size()) {
		return nil, errors.New("private key file read size error")
	}
	privateKey = PrivateKeyFromPEM(string(prikeyValue))
	return
}
