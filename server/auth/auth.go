package auth

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"heyui/server/token"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var thejwt = token.JWT{}

func Initialize() {
	prvKey, err := ioutil.ReadFile(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		fmt.Print(err)
	}
	pubKey, err := ioutil.ReadFile(os.Getenv("PUBLIC_KEY"))
	if err != nil {
		fmt.Print(err)
	}

	thejwt = token.NewJWT(prvKey, pubKey)
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func CreateToken(acct string) (string, error) {
	return thejwt.Create(acct)
}

func ValidateToken(r *http.Request) (interface{}, error) {
	tokenString := extractToken(r)
	return thejwt.Validate(tokenString)
}

func extractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}
