package db

import (
	"database/sql"
	"fmt"
	"log"
)

// CreateMigrationsTable creates the migrations table in the database.
func CreateMigrationsTable() error {
	query, err := Conn.Prepare(createMigrationsTableSQL())
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = query.Exec()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// DropMigrationsTable drops the migrations table from the database.
func DropMigrationsTable() error {
	query, err := Conn.Prepare(dropMigrationsTableSQL())
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = query.Exec()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// InsertMigration inserts a new migration into the migrations table.
func InsertMigration(id string) error {
	query, err := Conn.Prepare(insertMigrationSQL())
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = query.Exec(id)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// DeleteMigration deletes a migration from the migrations table.
func DeleteMigration(id string) error {
	query, err := Conn.Prepare(deleteMigrationSQL())
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = query.Exec(id)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// MigrationActive queries the migrations table for the migration ID and returns true if a result is found.
func MigrationActive(id string) (bool, error) {
	query, err := Conn.Prepare(selectMigrationSQL())
	if err != nil {
		log.Println(err)
		return false, err
	}

	var rowID int

	err = query.QueryRow(id).Scan(&rowID)
	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		log.Println(err)
		return false, err
	default:
		return true, nil
	}
}

// createMigrationsTableSQL returns the SQL for creating the migrations table.
func createMigrationsTableSQL() string {
	return fmt.Sprintf(
		"CREATE TABLE %s (id INT NOT NULL AUTO_INCREMENT, migration_id VARCHAR(255) NOT NULL UNIQUE, PRIMARY KEY(id))",
		MigrationsTableName,
	)
}

// dropMigrationsTableSQL returns the SQL for dropping the migrations table.
func dropMigrationsTableSQL() string {
	return fmt.Sprintf(
		"DROP TABLE %s",
		MigrationsTableName,
	)
}

// insertMigrationSQL returns the SQL for inserting a new migration into the migrations table.
func insertMigrationSQL() string {
	return fmt.Sprintf(
		"INSERT INTO %s (migration_id) VALUES (?)",
		MigrationsTableName,
	)
}

// selectMigrationSQL returns the SQL for selecting a migration from the migrations table.
func selectMigrationSQL() string {
	return fmt.Sprintf(
		"SELECT id FROM %s WHERE migration_id=?",
		MigrationsTableName,
	)
}

// deleteMigrationSQL returns the SQL for deleting a migration from the migrations table.
func deleteMigrationSQL() string {
	return fmt.Sprintf(
		"DELETE FROM %s WHERE migration_id=?",
		MigrationsTableName,
	)
}
