package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	jose "gopkg.in/square/go-jose.v2"
	"log"
	"testing"
)

func TestJwt(t *testing.T) {

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(publicKey), nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Checking token validity
	if !token.Valid {
		log.Fatal("invalid token")
	}

}

func TestLol(t *testing.T) {
	_, err := jose.ParseSigned(accessToken)
	if err != nil {
		log.Fatal(err)
	}

}
