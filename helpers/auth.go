package helpers

import (
	"fmt"
	"log"
	"os"
	"time"

	"aidanwoods.dev/go-paseto"

	password "github.com/dwin/goSecretBoxPassword"
	"github.com/joho/godotenv"
)

func GetSecret() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading Env File")
	}
	return os.Getenv("SECRET")
}

func GetSecretKey() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading Env File")
	}
	return os.Getenv("TOKENSECRET")
}

func GenerateHash(pass string) string {
	pwHash, err := password.Hash(pass, GetSecret(), 0, password.ScryptParams{N: 32768, R: 16, P: 1}, password.DefaultParams)
	if err != nil {
		fmt.Println("Hash fail. ", err)
	}
	return pwHash
}

func ValidateHash(hash, pass string) bool {
	err := password.Verify(pass, GetSecret(), hash)
	if err != nil {
		return false
	}
	return true
}

func GenerateToken(data string) string {
	key, err := paseto.V4SymmetricKeyFromHex(GetSecretKey())
	if err != nil {
		return ""
	}
	token := paseto.NewToken()

	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(2 * time.Hour))

	token.SetString("user-id", data)

	encrypted := token.V4Encrypt(key, nil)
	return encrypted
}

func GenerateRefreshToken(data string) string {
	key, err := paseto.V4SymmetricKeyFromHex(GetSecretKey())
	if err != nil {
		return ""
	}
	token := paseto.NewToken()
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(24 * time.Hour))
	token.SetString("user-id", data)

	encrypted := token.V4Encrypt(key, nil)
	return encrypted
}

func ValidateAccessToken(token string) (string, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())
	parser.AddRule(paseto.ValidAt(time.Now()))
	log.Println("Token at Validation")
	key, err := paseto.V4SymmetricKeyFromHex(GetSecretKey())
	if err != nil {
		log.Println("error fetching key")
		return "", err
	}

	parsedToken, err := parser.ParseV4Local(key, token, nil)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	return parsedToken.GetString("user-id")

}

func ValidateRefreshToken(token string) (string, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())
	parser.AddRule(paseto.ValidAt(time.Now()))
	log.Println("Token at Validation")
	key, err := paseto.V4SymmetricKeyFromHex(GetSecretKey())
	if err != nil {
		log.Println("error fetching key")
		return "", err
	}

	parsedToken, err := parser.ParseV4Local(key, token, nil)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	return parsedToken.GetString("user-id")

}

func TokenParser(token string) (string, error) {
	parser := paseto.NewParser()
	key, err := paseto.V4SymmetricKeyFromHex(GetSecretKey())
	if err != nil {
		log.Println("error fetching key")
		return "", err
	}
	parsedToken, err := parser.ParseV4Local(key, token, nil)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	return parsedToken.GetString("user-id")
}
