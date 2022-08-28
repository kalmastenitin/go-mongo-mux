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
	key := paseto.NewV4SymmetricKey()
	token := paseto.NewToken()

	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(2 * time.Hour))

	token.SetString("user-id", "<uuid>")

	encrypted := token.V4Encrypt(key, nil)
	return encrypted
}
