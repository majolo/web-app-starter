package jwt

import (
	"fmt"
	jwt "github.com/golang-jwt/jwt/v5"
)

func VerifyToken(encodedToken string, publicKey []byte) (*jwt.Token, error) {
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwt.ParseRSAPublicKeyFromPEM(publicKey)
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Now claims should contain the user information
		fmt.Println(claims)
	} else {
		return nil, fmt.Errorf("invalid token")
	}
	return token, nil
}

func VerifyTokenWithSecret(encodedToken string, secret []byte) (*jwt.Token, error) {
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extract user info from claims
		fmt.Println(claims)
	} else {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}
