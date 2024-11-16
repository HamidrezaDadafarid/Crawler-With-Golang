package middlewares

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type AuthTestSuite struct {
	suite.Suite
}

func (suite *AuthTestSuite) SetupSuite() {

	password := "test_password"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// Set the environment variable
	err = os.Setenv("PASSWORD", string(hashedPassword))
	if err != nil {
		log.Fatalf("Failed to set environment variable: %v", err)
	}

}

func (suite *AuthTestSuite) TearDownSuite() {

	// Unset the environment variable
	os.Unsetenv("PASSWORD")

}

func (suite *AuthTestSuite) TestValidPassword() {
	// Test a valid password
	pass := "test_password"

	_, result := Authentication(pass)

	suite.True(result, "Expected the password to be valid")
}

func (suite *AuthTestSuite) TestInvalidPassword() {
	// Test with an invalid password
	pass := "invalid_password"

	_, result := Authentication(pass)

	suite.False(result, "Expected the password to be invalid")
}

func (suite *AuthTestSuite) TestNoPasswordConfigured() {
	// Unset the environment variable to simulate the case where no password is set
	os.Unsetenv("PASSWORD")

	// Test with some password
	pass := "some_password"

	_, result := Authentication(pass)

	suite.False(result, "Expected the authentication to fail due to no configured password")
}

func (suite *AuthTestSuite) TestLoadEnvFailure() {
	// Simulate a failure to load .env by giving an incorrect path or no .env file
	err := godotenv.Load("nonexistent/.env")
	suite.Error(err, "Expected error while loading nonexistent .env file")
}

func TestAuthTestSuite(t *testing.T) {
	// Run the test suite
	suite.Run(t, new(AuthTestSuite))
}
