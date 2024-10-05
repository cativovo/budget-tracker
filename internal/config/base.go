package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
}

// https://github.com/joho/godotenv?tab=readme-ov-file#precedence--conventions
func loadEnv() []string {
	env := os.Getenv("BUDGET_TRACKER_ENV")
	if env == "" {
		env = "development"
	}

	var loadedFiles []string

	log.Println("loading env files")
	log.Println("env:", env)

	f := ".env." + env + ".local"
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

	f = ".env." + env
	if err := godotenv.Load(f); err != nil {
		log.Println(err)
	} else {
		loadedFiles = append(loadedFiles, f)
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
		Port: os.Getenv("PORT"),
	}, nil
}
