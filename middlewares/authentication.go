package middlewares

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func Authentication(pass string) bool {
	// Loads .env file from main directory
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Couldn't open .env file")
		return false
	}

	// Using (PASSWORD) as a key to search for hashed pass in env file
	storedH := os.Getenv("PASSWORD")

	if storedH == "" {
		log.Println("NO PASSWORD HAS BEEN CONFIGURED!")
		return false
	}
	// Compares hash with given password
	err = bcrypt.CompareHashAndPassword([]byte(storedH), []byte(pass))
	if err != nil {
		log.Println("Invalid password!")
		return false
	}
	log.Println("Vallid password")
	return true

}
