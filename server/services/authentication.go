package services

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	uuid "github.com/satori/go.uuid"
)

type AuthenticationService interface {
	Init(ctx) error //loads a new basic into memory and inits AuthenticationServerID
	handleHello(clientID []byte) (MsgHelloAuthenticate, error)
	handleAuthenticate(clientID []byte, password []byte) (*MsgWelcome, error)
	handleJWT(jwtToken string) (*JWTClaims, error)
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

type MsgHelloAuthenticate struct {
	InstanceID string
	UserID     string
}

func (l *authenticationService) handleHello(clientID []byte) (MsgHelloAuthenticate, error) {
	now := time.Now()
	return MsgHelloAuthenticate{
		InstanceID: l.instanceID.String(),
		UserID:     fmt.Sprintf("%x-%x-%d", []byte(l.instanceID.String()), clientID, now.Unix()),
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

func (l *authenticationService) handleAuthenticate(clientID []byte, password []byte) (*MsgWelcome, error) {
	if string(password) != "dev" {
		return nil, fmt.Errorf("dev mode enabled. log in with password 'dev'")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"client_id": string(clientID),
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
