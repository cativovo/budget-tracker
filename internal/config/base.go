package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	Env    string
	Port   string
	DBPath string
}

const envKey = "BUDGET_TRACKER_ENV"

// https://github.com/joho/godotenv?tab=readme-ov-file#precedence--conventions
func loadEnv(logger *zap.SugaredLogger) []string {
	var loadedFiles []string
	env := os.Getenv(envKey)

	logger.Info("loading env files")

	if env == "" {
		env = "development"
		os.Setenv(envKey, env)
	}

	logger.Infof("env: %s", env)

	f := ".env." + env
	if err := godotenv.Load(f); err != nil {
		logger.Warn(err)
	} else {
		loadedFiles = append(loadedFiles, f)
	}

	if env == "production" {
		return loadedFiles
	}

	f = ".env." + env + ".local"
	if err := godotenv.Load(".env." + env + ".local"); err != nil {
		logger.Warn(err)
	} else {
		loadedFiles = append(loadedFiles, f)
	}

	if env != "test" {
		f = ".env.local"
		err := godotenv.Load(f)
		if err != nil {
			logger.Warn(err)
		} else {
			loadedFiles = append(loadedFiles, f)
		}
	}

	// The Original .env
	if err := godotenv.Load(); err != nil {
		logger.Warn(err)
	} else {
		loadedFiles = append(loadedFiles, ".env")
	}

	return loadedFiles
}

func LoadConfig(logger *zap.SugaredLogger) (Config, error) {
	loadedFiles := loadEnv(logger)

	if len(loadedFiles) == 0 {
		return Config{}, errors.New("no env file found")
	}

	for _, f := range loadedFiles {
		logger.Info("env vars loaded from", f)
	}

	return Config{
		Port:   os.Getenv("PORT"),
		DBPath: os.Getenv("DB_PATH"),
		Env:    os.Getenv(envKey),
	}, nil
}
