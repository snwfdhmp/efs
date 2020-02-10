package services

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type AuthenticationService interface {
	Init(ctx) error //loads a new basic into memory and inits AuthenticationServerID
	jwtAuth(jwtToken []byte) error
	pgpHello(clientID []byte) (helloCode string, err error)
	pgpAuth(clientID []byte, helloCode string) error
}

func NewAuthenticationService() AuthenticationService {
	return &authenticationService{}
}

type authenticationService struct {
	instanceID     uuid.UUID
	instanceSecret uuid.UUID
	starterKey     []byte
	basicKey       [32]byte
	jwtKey         []byte
}

func (l *authenticationService) Init(ctx ctx) error {
	var err error

	l.instanceID, err = uuid.NewV4()
	if err != nil {
		return err
	}

	l.instanceSecret, err = uuid.NewV4()
	if err != nil {
		return err
	}

	l.basicKey = sha256.Sum256([]byte(fmt.Sprintf("%s:%s", l.instanceID, l.instanceSecret)))
	l.jwtKey = make([]byte, 128)
	_, err = rand.Read(l.jwtKey)
	if err != nil {
		return err
	}

	return nil
}

func (l *authenticationService) jwtAuth(jwtToken []byte) error {
	return nil
}

func (l *authenticationService) pgpHello(clientID []byte) (helloCode string, err error) {
	return "testcode", nil
}
func (l *authenticationService) pgpAuth(clientID []byte, helloCode string) error {
	return nil
}
