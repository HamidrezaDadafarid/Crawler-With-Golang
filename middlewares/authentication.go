package middlewares

import (
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func Authentication(pass string) (error, bool) {

	err := godotenv.Load("./.env")
	if err != nil {
		return err, false
	}

	storedH := os.Getenv("SUPER_ADMIN_PASSWORD")

	err = bcrypt.CompareHashAndPassword([]byte(storedH), []byte(pass))
	if err != nil {
		return err, false
	}
	return nil, true
}
