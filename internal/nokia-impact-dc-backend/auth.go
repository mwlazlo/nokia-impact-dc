package nokia_impact_dc_backend

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// it's not really token based auth
type Token struct {
	Id     string
	Passwd string
}

func AuthMiddleware(appCtx *AppContext) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.Method, "-", r.URL)

			authHdr, ok := r.Header["Authorization"]
			if !ok {
				log.Println("no auth header")
				HttpError(w, http.StatusUnauthorized)
				return
			}

			authArr := strings.Split(authHdr[0], " ")
			if len(authArr) != 2 || authArr[0] != "Basic" {
				HttpError(w, http.StatusBadRequest)
				return
			}

			if tok, err := verifyToken(appCtx, authArr[1], r); err != nil {
				HttpError(w, http.StatusUnauthorized)
				return
			} else {
				ctx := context.WithValue(r.Context(), "uid", tok.Id)
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		})
	}
}

func verifyToken(appCtx *AppContext, token string, r *http.Request) (*Token, error) {
	cfg := appCtx.Config
	data, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}
	cred := strings.Split(string(data), ":")

	if cred[0] != cfg.CallbackUsername || cred[1] != cfg.CallbackPassword {
		log.Println("Credential Mismatch", cred[0], "!=", cfg.CallbackUsername, " ",
			cred[1], "!=", cfg.CallbackPassword)
		return nil, errors.New("Credential Mismatch")
	}

	log.Println("Request authenticated")

	return &Token{
		cred[0],
		cred[1],
	}, nil
}
