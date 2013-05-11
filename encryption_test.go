package hamster

import (
	"testing"
)

func Testbcrypt(t *testing.T) {

	password := "password"

	//encrypt, get hash and salt
	hash, salt, err := encryptPassword(password)

	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	//decrypt and match

	if matched := matchPassword(password, hash, salt); !matched {
		t.Fatalf("match failed: %v", err)
	}
}

func Testbase64(t *testing.T) {
	token := "518b65cdcde9e8116e000001"

	//encode it
	encoded_token := encodeBase64Token(token)

	//decode and match
	if decoded_token := decodeToken(encoded_token); token != decoded_token {

		t.Fatalf("decoding match failed: %v")

	}

}
