/*Generate public key, write to pem,
load from pem*/
package key

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

func PublicKeyFromPEMFile(filename string) (publicKey *rsa.PublicKey, err error) {
	pubkeyFile, err := os.Open(filename)
	if err != nil {
		return
	}
	pubkeyFileStat, err := pubkeyFile.Stat()
	if err != nil {
		return
	}
	pubkeyValue := make([]byte, pubkeyFileStat.Size())
	n, err := pubkeyFile.Read(pubkeyValue)
	if err != nil {
		return
	}
	if n != int(pubkeyFileStat.Size()) {
		return nil, errors.New("public key file read size error")
	}
	publicKey = PublicKeyFromPEM(string(pubkeyValue))
	return
}

func PublicKeyFromPEM(pubkeyPEM string) *rsa.PublicKey {
	pubkeyMarshaled, _ := pem.Decode([]byte(pubkeyPEM))
	pubkey, err := x509.ParsePKIXPublicKey(pubkeyMarshaled.Bytes)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return pubkey.(*rsa.PublicKey)
}

func PublicKeyToPEM(pubkey rsa.PublicKey) string {
	pubkeyMarshaled, err := x509.MarshalPKIXPublicKey(&pubkey)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	pubkeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Headers: nil, Bytes: pubkeyMarshaled})
	return string(pubkeyPEM)
}

func WritePublicKeyToPEMFile(pubkey rsa.PublicKey, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	err = file.Chmod(os.FileMode(0644))
	if err != nil {
		return err
	}
	_, err = file.WriteString(PublicKeyToPEM(pubkey))
	if err != nil {
		return err
	}
	return nil
}
