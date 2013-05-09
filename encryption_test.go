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
