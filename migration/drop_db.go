package migration

import (
	"fmt"
	"log"

	"github.com/nicday/turtle/config"
	"github.com/nicday/turtle/db"
)

// DropDB removes the database from the host.
func DropDB() {
	_, err := db.Conn.Exec(fmt.Sprintf("DROP DATABASE %s", config.DBName))
	if err != nil {
		log.Fatal(err)
	}
}
