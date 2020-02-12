package services

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	starterKey = "AEFDAEFDAEFDAEFD"
)

type AuthenticationService interface {
	Init(ctx) error //loads a new basic into memory and inits AuthenticationServerID
	handleAuthenticate(clientID string, password string) (*MsgWelcome, error)
	handleJWT(jwtToken string) (*JWTClaims, error)
}

func NewAuthenticationService() AuthenticationService {
	return &authenticationService{
		starterKey: []byte(starterKey),
	}
}

type authenticationService struct {
	instanceID string
	starterKey []byte
	jwtKey     []byte
}

func (l *authenticationService) Init(ctx ctx) error {
	var err error

	instanceID := make([]byte, 128)
	_, err = rand.Read(instanceID)
	if err != nil {
		return err
	}
	l.instanceID = fmt.Sprintf("%x", instanceID)

	l.jwtKey = make([]byte, 128)
	_, err = rand.Read(l.jwtKey)
	if err != nil {
		return err
	}

	return nil
}

type MsgHelloAuthenticate struct {
	InstanceID string
	UserID     string
}

func (l *authenticationService) handleHello(clientID []byte) (MsgHelloAuthenticate, error) {
	now := time.Now()
	return MsgHelloAuthenticate{
		InstanceID: l.instanceID,
		UserID:     fmt.Sprintf("%x-%x-%d", []byte(l.instanceID), clientID, now.Unix()),
	}, nil
}

type MsgWelcome struct {
	JWTToken string
}

type JWTClaims struct {
	ClientID   string `json:"client_id"`
	Expiration int    `json:"exp"`
	jwt.StandardClaims
}

func (l *authenticationService) handleAuthenticate(clientID string, password string) (*MsgWelcome, error) {
	if password != "dev" {
		return nil, fmt.Errorf("dev mode enabled. log in with password 'dev'")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"client_id": clientID,
		"exp":       time.Now().Add(time.Minute * 30).Unix(),
	})
	tokenString, err := token.SignedString(l.jwtKey)
	if err != nil {
		return nil, err
	}
	return &MsgWelcome{
		JWTToken: tokenString,
	}, nil
}

func (l *authenticationService) handleJWT(jwtToken string) (*JWTClaims, error) {
	claims := JWTClaims{}
	token, err := jwt.ParseWithClaims(jwtToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return l.jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		fmt.Printf("%v", claims.StandardClaims.ExpiresAt)
	} else {
		fmt.Println(err)
	}

	return &claims, nil
}
