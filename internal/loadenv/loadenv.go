package loadenv

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	env, err := godotenv.Read()
	if err != nil {
		env = map[string]string{}
	}
	for key, value := range env {
		_, exists := os.LookupEnv(key)
		if !exists {
			if err := os.Setenv(key, value); err != nil {
				panic(fmt.Errorf("setenv: %w", err))
			}
		}
	}
}
