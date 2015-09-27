package migration

import (
	"fmt"
	"log"

	"github.com/nicday/turtle/config"
	"github.com/nicday/turtle/db"
)

// CreateDB creates the database on the host.
func CreateDB() {
	_, err := db.Conn.Exec(fmt.Sprintf("CREATE DATABASE %s", config.DBName))
	if err != nil {
		log.Fatal(err)
	}
}
