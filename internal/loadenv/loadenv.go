package loadenv

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	env, err := godotenv.Read()
	if err != nil {
		env = map[string]string{}
	}
	for key, value := range env {
		str, exists := os.LookupEnv(key)
		slog.Info("lookup", "lookup", str)
		if !exists {
			os.Setenv(key, value)
		}
	}
}