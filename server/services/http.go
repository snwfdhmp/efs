package services

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type ApiService interface {
	Listen(addr string) // creates new router and initializes it
}

type apiService struct {
	ctx Ctx
}

func NewAPIService(ctx Ctx) ApiService {
	return &apiService{ctx: ctx}
}

func (l *apiService) Listen(addr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authentication"), " ")
		if len(authHeader) != 2 {
			w.WriteHeader(401)
			io.WriteString(w, "401 - Bad Authentication\n")
			return
		}
		switch strings.ToLower(authHeader[0]) { //protocol
		case "efs-handshake-hello":
			helloCode, err := l.ctx.Auth().pgpHello([]byte(authHeader[1]))
			if err != nil {
				l.ctx.Log().Error(err.Error())
				w.WriteHeader(500)
				io.WriteString(w, "500 - Internal error\n")
				return
			}
			_, err = io.WriteString(w, fmt.Sprintf("%s\n", helloCode))
			if err != nil {
				l.ctx.Log().Error(err.Error())
				w.WriteHeader(500)
				io.WriteString(w, "500 - Internal error\n")
				return
			}
		case "efs-handshake-authenticate":
			// verify code signature
			// return jwt token
		case "efs-jwt":
			// parse jwt and verify
			// update context accordingly
		default:
			w.WriteHeader(401)
			io.WriteString(w, "401 - Bad Authentication Method\n")
			return
		}

		switch r.Method {
		case "GET":
			decrypted, err := l.ctx.GetFile(r.URL.String())
			if err != nil {
				w.WriteHeader(500)
				io.WriteString(w, err.Error())
				return
			}
			if _, err = w.Write(decrypted); err != nil {
				w.WriteHeader(500)
				io.WriteString(w, err.Error())
				return
			}
			return
		case "POST":
			decrypted, err := ioutil.ReadAll(r.Body)
			if err != nil {
				io.WriteString(w, err.Error())
				w.WriteHeader(500)
				return
			}
			l.ctx.PostFile(r.URL.String(), decrypted)
			if err != nil {
				io.WriteString(w, err.Error())
				w.WriteHeader(500)
				return
			}
			w.WriteHeader(200)
			return
		default:
			w.WriteHeader(405)
			return
		}
		return
	})
	http.ListenAndServe(addr, mux)
}
