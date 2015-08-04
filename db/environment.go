package db

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	defaultMigrationsTableName = "migrations"
	defaultDBPort              = "3306"
	defaultDBUser              = "root"
)

var (
	MigrationsTableName string
	dbHost              string
	dbPort              string
	dbName              string
	dbUser              string
	dbPassword          string

	// ErrNoDBHost is raised when there is no DB_HOST in the environment variables
	ErrNoDBHost = errors.New("DB_HOST not found in environment variables")

	// ErrNoDBPassword is raised when there is no DB_PASSWORD in the environment variables
	ErrNoDBPassword = errors.New("DB_PASSWORD not found in environment variables")

	// ErrNoDBName is raised when there is no DB_NAME in the environment variables
	ErrNoDBName = errors.New("DB_NAME not found in environment variables")
)

func initDBEnv() {
	if IsTestEnv() {
		return
	}

	// Don't worry about an error here, .env might not be present; So long as we have the environment variables required.
	godotenv.Load()

	MigrationsTableName = os.Getenv("MIGRATIONS_TABLE_NAME")
	if MigrationsTableName == "" {
		MigrationsTableName = defaultMigrationsTableName
	}

	dbHost = os.Getenv("DB_HOST")
	if dbHost == "" {
		log.Fatal(ErrNoDBHost)
	}

	dbPort = os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = defaultDBPort
	}

	dbName = os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal(ErrNoDBName)
	}

	dbUser = os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = defaultDBUser
	}

	dbPassword = os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		log.Fatal(ErrNoDBPassword)
	}
}

// IsTestEnv returns true when the ENV=test
func IsTestEnv() bool {
	env := os.Getenv("ENV")
	if env == "test" {
		return true
	}

	return false
}
