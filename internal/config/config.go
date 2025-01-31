package config

import (
	"github.com/joho/godotenv"
)

func Load(path string) error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}

	return nil
}
