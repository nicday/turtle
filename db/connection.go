package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"

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

func init() {
	initDBEnv()
	if IsTestEnv() {
		return
	}

	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

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
	expBackoff.MaxElapsedTime = time.Duration(30) * time.Second

	err := backoff.Retry(pingDB, expBackoff)
	if err != nil {
		log.Fatal(ErrUnableToConnectToDB)
	}

	return nil
}