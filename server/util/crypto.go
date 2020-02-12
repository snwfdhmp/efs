package util

import (
	"github.com/kevinburke/nacl"
	"github.com/kevinburke/nacl/secretbox"
)

var (
	naclKeyString = "56b2a39f271334de4cbfafedbf40d50364cfa10370b9f9522fc696859a6e48d1"
	key           nacl.Key
)

func init() {
	var err error
	key, err = nacl.Load(naclKeyString)
	if err != nil {
		panic(err)
	}
}

func SealED25519(key nacl.Key, msg []byte) ([]byte, error) {
	encrypted := secretbox.EasySeal([]byte(msg), key)
	return encrypted, nil
}

func OpenED25519(key nacl.Key, box []byte) ([]byte, error) {
	decrypted, err := secretbox.EasyOpen(box, key)
	if err != nil {
		return nil, err
	}
	return decrypted, nil
}
