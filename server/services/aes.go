package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

type AESService interface {
	decrypt(encryptedFile []byte) ([]byte, error)
	encrypt(content []byte) ([]byte, error)
}

type aesService struct {
	aesKey []byte
}

func NewAESService(key []byte) AESService {
	return &aesService{
		aesKey: key,
	}
}

func (l *aesService) decrypt(encryptedFile []byte) ([]byte, error) {
	cipherText, err := base64.URLEncoding.DecodeString(string(encryptedFile))
	if err != nil {
		return []byte{}, err
	}

	block, err := aes.NewCipher(l.aesKey)
	if err != nil {
		return []byte{}, err
	}

	if len(cipherText) < aes.BlockSize {
		return []byte{}, errors.New("Ciphertext block size is too short!")
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)

	return cipherText, nil
}

func (l *aesService) encrypt(content []byte) ([]byte, error) {
	block, err := aes.NewCipher(l.aesKey)
	if err != nil {
		return []byte{}, err
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	cipherText := make([]byte, aes.BlockSize+len(content))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return []byte{}, nil
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], content)

	//returns to base64 encoded string
	return []byte(base64.URLEncoding.EncodeToString(cipherText)), nil
}
