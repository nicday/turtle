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

	// BackoffTimeout is the total time the backoff with wait before failing.
	BackoffTimeout = time.Duration(30) * time.Second
)

func init() {
	config.InitEnv()
	if config.IsTestEnv() {
		return
	}

	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)

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

// VerifyConnection pings the database to verify a connection is established. If the connection cannot be established,
// it will retry with an exponential back off.
func VerifyConnection(c *sql.DB) error {
	pingDB := func() error {
		return c.Ping()
	}

	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.MaxElapsedTime = BackoffTimeout

	err := backoff.Retry(pingDB, expBackoff)
	if err != nil {
		log.Fatal(ErrUnableToConnectToDB)
	}

	return nil
}
