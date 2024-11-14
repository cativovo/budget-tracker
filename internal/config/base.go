package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Password string
	User     string
	DB       string
	Host     string
	Port     string
	SSL      string
}

type Config struct {
	Env  string
	Port string
	DB   DBConfig
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

	dbCfg := DBConfig{
		Password: os.Getenv("POSTGRES_PASSWORD"),
		User:     os.Getenv("POSTGRES_USER"),
		DB:       os.Getenv("POSTGRES_DB"),
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		SSL:      os.Getenv("POSTGRES_SSL"),
	}

	return Config{
		Port: os.Getenv("PORT"),
		DB:   dbCfg,
		Env:  os.Getenv(envKey),
	}, nil
}
