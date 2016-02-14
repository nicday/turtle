package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/nicday/turtle/config"

	"github.com/cenkalti/backoff"
)

var (
	// Conn the active database connection
	Conn *sql.DB

	// ErrUnableToParseDBConnection is raised when there are missing or invalid details for the database connection.
	ErrUnableToParseDBConnection = errors.New("unable to parse database connection details")

	// ErrUnableToConnectToDB is raised when a connection to the database cannot be established.
	ErrUnableToConnectToDB = errors.New("unable to connect to the database")
)

// InitConnection initializes the database connection
func InitConnection() {
	err := config.InitEnv()
	if err != nil {
		log.Fatal(err)
	}

	if config.IsTestEnv() {
		return
	}
	connString := ConnString()

	c, err := sql.Open("mysql", connString)
	if err != nil {
		log.Println("[Error]", err)
		log.Fatal(ErrUnableToParseDBConnection)
	}

	err = VerifyConnection(c)
	if err != nil {
		log.Println("[Error]", err)
		log.Fatal(ErrUnableToConnectToDB)
	}

	Conn = c
}

// ConnString returns the connection string for the database driver.
func ConnString() string {
	switch config.DBDriver {
	case "mysql":
		return mysqlConnString()
	case "postgres":
		return postgresConnString()
	default:
		log.Fatal(config.ErrUnknownDBDriver)
		return ""
	}
}

func connCredentials() string {
	if config.DBPassword != "" {
		return fmt.Sprintf("%s:%s", config.DBUser, config.DBPassword)
	}
	return config.DBUser
}

func mysqlConnString() string {
	return fmt.Sprintf("%s@tcp(%s:%s)/", connCredentials(), config.DBHost, config.DBPort)
}

func postgresConnString() string {
	return fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", connCredentials(), config.DBHost, config.DBPort, config.DBName)
}

// VerifyConnection pings the database to verify a connection is established. If the connection cannot be established,
// it will retry with an exponential back off.
func VerifyConnection(c *sql.DB) error {
	pingDB := func() error {
		return c.Ping()
	}

	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.MaxElapsedTime = time.Duration(30) * time.Second

	err := backoff.Retry(pingDB, expBackoff)
	if err != nil {
		return err
	}

	return nil
}

// UseDB runs the `USE` SQL command, ensuring that all future SQL commands on the database connection use the named
// database.
func UseDB() error {
	_, err := Conn.Exec(fmt.Sprintf("USE %s", config.DBName))

	if err != nil {
		return err
	}

	return nil
}
