package migration

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

const (
	timeFormat   = "20060102150405.000"
	migrationDir = "migrations"
)

// Generate creates up and down migration files.
func Generate(name string) error {
	assertMigrationDir()

	baseFilename := fmt.Sprintf("%s_%s", timestamp(), name)

	for _, direction := range []string{"up", "down"} {
		filename := fmt.Sprintf("%s_%s.sql", baseFilename, direction)
		createMigrationFile(filename)
	}

	return nil
}

// createMigrationFile creates the migration file in the migration directory.
func createMigrationFile(name string) error {
	_, err := os.Create(path.Join(migrationDir, name))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// assertMigrationDir ensures that the migration direction exists and raises an error if it cannot be created.
func assertMigrationDir() error {
	_, err := os.Stat(migrationDir)
	if err != nil {
		err := os.Mkdir(migrationDir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

// timestamp returns a timestamp with millisecond accuracy and no decimal place.
func timestamp() string {
	t := time.Now().Format(timeFormat)
	return strings.Replace(t, ".", "", 1)
}
