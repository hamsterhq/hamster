/*
Generates private and public keys
*/
package encryption

import (
	"encryption/key"
	"fmt"
	"os"
	"strconv"
)

func main() {
	filenamePrivateKey := "key.pri"
	filenamePublicKey := "key.pub"
	bits := 1024
	if len(os.Args) > 1 {
		bitsZ, err := strconv.ParseInt(os.Args[1], 10, 0)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		bits = int(bitsZ)
		if bits < 1024 {
			fmt.Printf("Expected a bit size larger >= 1024 but got %d instead.\n", bits)
			os.Exit(1)
		}
	}

	fmt.Printf("Generating %d-bit RSA private key...\n", bits)
	prikey := key.GenPrivateKeyPair(bits)

	err := key.WritePrivateKeyToPEMFile(prikey, filenamePrivateKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Wrote %d-bit RSA private key to %s.\n", bits, filenamePrivateKey)

	err = key.WritePublicKeyToPEMFile(prikey.PublicKey, filenamePublicKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Wrote %d-bit RSA public key to %s.\n", bits, filenamePublicKey)
}
