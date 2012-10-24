/*load public and private pem file
decrypt the data with key.pri
encrypt the data with key.pub*/
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/json"
	"encryption/key"
	"fmt"
	"os"
	"path/filepath"
)

type KeyConfig struct {
	PrivateKeyFile string // relative to the config file's 
	PublicKeyFile  string
}

var config *KeyConfig = new(KeyConfig)
var configFileName = "conf.json"

//load filepath from conf.json
func loadConfig() {

	if len(os.Args) > 1 {
		configFileName = os.Args[1]
	}
	configFile, err := os.Open(configFileName)
	if err != nil {
		fmt.Println("[ERROR] " + err.Error())

		os.Exit(1)
	}
	configDecoder := json.NewDecoder(configFile)
	err = configDecoder.Decode(config)
	if err != nil {
		fmt.Println("[CONFIG FILE FORMAT ERROR] " + err.Error())
		fmt.Println("Please ensure that your config file is in valid JSON format.")
		os.Exit(1)
	}

	switch {
	case config.PrivateKeyFile == "":
		fmt.Println("[CONFIG ERROR] PrivateKeyFile is missing.")
		os.Exit(1)
	case config.PublicKeyFile == "":
		fmt.Println("[CONFIG ERROR] PublicKeyFile is missing.")
		os.Exit(1)
	}
}

func main() {

	loadConfig()

	//read private key from pem file
	privateKeyFile := filepath.Clean(config.PrivateKeyFile)
	if !filepath.IsAbs(privateKeyFile) {
		privateKeyFile = filepath.Clean(filepath.Join(filepath.Dir(configFileName), privateKeyFile))
	}

	privateKey, err := key.PrivateKeyFromPEMFile(privateKeyFile)

	if err != nil {
		fmt.Printf("[PRIVATEKEY ERROR] %s\n", err)
		os.Exit(1)
	}

	//read public key from pem file
	publicKeyFile := filepath.Clean(config.PublicKeyFile)
	if !filepath.IsAbs(publicKeyFile) {
		publicKeyFile = filepath.Clean(filepath.Join(filepath.Dir(configFileName), publicKeyFile))
	}

	publicKey, err := key.PublicKeyFromPEMFile(publicKeyFile)

	if err != nil {
		fmt.Printf("[PUBLICKEY ERROR] %s\n", err)
		os.Exit(1)
	}

	//ok,lets encrypt this data

	data_in := "Hello World!"

	//encryption
	sha1 := sha1.New()

	out, err := rsa.EncryptOAEP(sha1, rand.Reader, publicKey, []byte(data_in), nil)
	if err != nil {
		fmt.Printf("error: %s", err)
	}

	fmt.Println("encrypted\n", string(out))
	//

	//decryption

	data_out, err2 := rsa.DecryptOAEP(sha1, nil, privateKey, out, nil)
	if err2 != nil {
		fmt.Printf("error: %s", err2)
	}

	fmt.Println("\ndecrypted\n", string(data_out))
}
