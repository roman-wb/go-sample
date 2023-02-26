package env

import (
	"fmt"

	"github.com/joho/godotenv"
)

const EnvFile = ".env"

func Init(filename string) error {
	if err := godotenv.Load(filename); err != nil {
		return fmt.Errorf("godotenv load %s file: %w", filename, err)
	}

	return nil
}
