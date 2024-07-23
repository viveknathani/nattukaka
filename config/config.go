package config

import (
	"github.com/joho/godotenv"
)

// LoadEnvFromFile attempts to load an .env file.
// Use Os.Getenv to read the variables after this.
func LoadEnvFromFile(relativePath string) error {
	err := godotenv.Load(relativePath)
	if err != nil {
		return err
	}
	return nil
}
