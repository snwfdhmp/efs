package util

import (
	"testing"

	"github.com/kevinburke/nacl"
)

var (
	testMsg      = "bonjour"
	testKeyWrong = "56b2a39f271334de4cbfafedbf40d50364cfa10371c9f9522fc696859a6e48d1"
)

func Test_SealED25519(t *testing.T) {
	encrypted, err := SealED25519(key, []byte(testMsg))
	if err != nil { //should work
		t.Error(err)
	}

	_, err = OpenED25519(key, encrypted)
	if err != nil { //should work
		t.Error(err)
	}

	//generate different key
	keyWrong, err := nacl.Load(testKeyWrong)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = OpenED25519(keyWrong, encrypted)
	if err == nil { //should fail
		t.Error("verification success with wrong key")
	}
}
