package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthClaim struct {
  Id string `json:"id"`
  User string `json:"user"`
  Admin bool `json:"role"`
  jwt.RegisteredClaims
}

func CreateNewAuthToken(id string, email string, isAdmin bool) (string, error)  {
  claims := AuthClaim{
	Id: id,
	User: email,
	Admin: isAdmin,
	RegisteredClaims: jwt.RegisteredClaims{
	  ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	  Issuer: "searchengine.com",
	},
  }
  // create token

}
