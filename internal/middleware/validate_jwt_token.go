package middleware

import (
	"errors"
	"github.com/StarkovPO/Go-shop-final/internal/service"
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

func CheckToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")
		if token == "" {
			next.ServeHTTP(w, r)
			return
		}
		UID, err := parseToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		r.Header.Add("User-ID", UID)

		next.ServeHTTP(w, r)
	})
}

func parseToken(accessToken string) (string, error) {

	token, err := jwt.ParseWithClaims(accessToken, &service.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", errors.New("invalid signing method")
		}
		return []byte("qiausydigswig104#hlk[pzxn"), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*service.TokenClaims)
	if !ok {
		return "", errors.New("token claims are not of type")
	}
	return claims.UserID, nil
}
