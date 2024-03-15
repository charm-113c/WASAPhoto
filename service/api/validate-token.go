package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Returns true if token is valid. Requires as input the username of the requesting user as well as the secret key
func validateToken(r *http.Request, userID string, secret string) error {
	// get token
	tokenstr := r.Header.Get("Authorization")
	if tokenstr == "" {
		return errors.New("no authentication token provided")
	}
	tokenstr = strings.TrimPrefix(tokenstr, "Bearer ")

	// parse and validate token
	token, err := jwt.Parse(tokenstr, func(token *jwt.Token) (interface{}, error) {
		// validating expected algorithm, i.e. HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	// check for errors/validity
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("unauthorized, token invalid")
	}

	// check that requesting user == token's subject
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if claims["userID"] != userID {
			log.Println("Unauthorized operation attempt: requesting user mismatch with token subject")
			return errors.New("unauthorized, token invalid")
		}
		return nil
	} else {
		return errors.New("token claims could not be validated")
	}
}
