package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env    string
	Port   string
	DBPath string
}

const envKey = "BUDGET_TRACKER_ENV"

// https://github.com/joho/godotenv?tab=readme-ov-file#precedence--conventions
func loadEnv() []string {
	var loadedFiles []string
	env := os.Getenv(envKey)

	log.Println("loading env files")
	log.Println("env:", env)

	if env == "" {
		env = "development"
		os.Setenv(envKey, env)
	}

	f := ".env." + env
	if err := godotenv.Load(f); err != nil {
		log.Println(err)
	} else {
		loadedFiles = append(loadedFiles, f)
	}

	if env == "production" {
		return loadedFiles
	}

	f = ".env." + env + ".local"
	if err := godotenv.Load(".env." + env + ".local"); err != nil {
		log.Println(err)
	} else {
		loadedFiles = append(loadedFiles, f)
	}

	if env != "test" {
		f = ".env.local"
		err := godotenv.Load(f)
		if err != nil {
			log.Println(err)
		} else {
			loadedFiles = append(loadedFiles, f)
		}
	}

	// The Original .env
	if err := godotenv.Load(); err != nil {
		log.Println(err)
	} else {
		loadedFiles = append(loadedFiles, ".env")
	}

	return loadedFiles
}

func LoadConfig() (Config, error) {
	loadedFiles := loadEnv()

	if len(loadedFiles) == 0 {
		return Config{}, errors.New("no env file found")
	}

	for _, f := range loadedFiles {
		log.Println("env vars loaded from", f)
	}

	return Config{
		Port:   os.Getenv("PORT"),
		DBPath: os.Getenv("DB_PATH"),
		Env:    os.Getenv(envKey),
	}, nil
}
