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
		slog.Info("env", key, value)
		os.Setenv(key, value)
	}
}