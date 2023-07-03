package user

import (
	"github.com/dgrijalva/jwt-go"
)

// Claims struct represents the claims for JWT
type Claims struct {
	Userid   int    `json:"userid"`
	Username string `json:"username"`
	jwt.StandardClaims
}
