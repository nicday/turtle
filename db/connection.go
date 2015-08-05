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
	err := InitConn()
	if err != nil {
		log.Println("[Error]", err)
	}

	config.DontRunInTest(VerifyConnection)
}

// InitConn initializes the database connection. This is called by init but is also exported to allow testing.
func InitConn() error {
	config.InitEnv()
	if config.IsTestEnv() {
		return nil
	}

	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)

	c, err := sql.Open("mysql", connString)
	if err != nil {
		log.Println("[Error]", ErrUnableToParseDBConnection)
		return err
	}

	Conn = c

	return nil
}

// VerifyConnection pings the database to verify a connection is established. If the connection cannot be established,
// it will retry with an exponential back off.
func VerifyConnection() {
	pingDB := func() error {
		return Conn.Ping()
	}

	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.MaxElapsedTime = BackoffTimeout

	err := backoff.Retry(pingDB, expBackoff)
	if err != nil {
		log.Println("[Error]", ErrUnableToConnectToDB)
		panic(err)
	}
}
