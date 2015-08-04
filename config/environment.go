package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	defaultMigrationsTableName = "migrations"
	defaultMigrationsPath      = "migrations"
	defaultDBPort              = "3306"
	defaultDBUser              = "root"
)

var (
	// MigrationsTableName is the table name where migrations are logged in the database.
	MigrationsTableName = defaultMigrationsTableName

	// MigrationsPath is the location that migration files will loaded from the filesystem.
	MigrationsPath = defaultMigrationsPath

	// DBHost is the host address when the database is running.
	DBHost string

	// DBPort is the port the database is running on.
	DBPort string

	// DBName is the database name to preform migrations on.
	DBName string

	// DBUser is the username to use when preforming migrations.
	DBUser string

	// DBPassword is the password to use for the database user.
	DBPassword string

	// ErrNoDBHost is raised when there is no DB_HOST in the environment variables
	ErrNoDBHost = errors.New("DB_HOST not found in environment variables")

	// ErrNoDBPassword is raised when there is no DB_PASSWORD in the environment variables
	ErrNoDBPassword = errors.New("DB_PASSWORD not found in environment variables")

	// ErrNoDBName is raised when there is no DB_NAME in the environment variables
	ErrNoDBName = errors.New("DB_NAME not found in environment variables")
)

// InitEnv initializes the environment variables. An attempt will be made to load variables from a `.env`, this can
// silently fail, so long as validation passes for the required variables.
func InitEnv() {
	if IsTestEnv() {
		return
	}

	// Don't worry about an error here, .env might not be present; So long as we have the environment variables required.
	godotenv.Load()

	MigrationsTableName = os.Getenv("MIGRATIONS_TABLE_NAME")
	if MigrationsTableName == "" {
		MigrationsTableName = defaultMigrationsTableName
	}

	MigrationsPath = os.Getenv("MIGRATIONS_PATH")
	if MigrationsPath == "" {
		MigrationsPath = defaultMigrationsPath
	}

	DBHost = os.Getenv("DB_HOST")
	if DBHost == "" {
		log.Fatal(ErrNoDBHost)
	}

	DBPort = os.Getenv("DB_PORT")
	if DBPort == "" {
		DBPort = defaultDBPort
	}

	DBName = os.Getenv("DB_NAME")
	if DBName == "" {
		log.Fatal(ErrNoDBName)
	}

	DBUser = os.Getenv("DB_USER")
	if DBUser == "" {
		DBUser = defaultDBUser
	}

	DBPassword = os.Getenv("DB_PASSWORD")
	if DBPassword == "" {
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
