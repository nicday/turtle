package migration

import (
	"fmt"
	"log"
	"path"
	"regexp"
	"strings"

	"github.com/nicday/turtle/config"
	"github.com/nicday/turtle/db"
)

var (
	upMigrationRegex   = regexp.MustCompile(`(\d+)_([\w-]+)_up\.sql`)
	downMigrationRegex = regexp.MustCompile(`(\d+)_([\w-]+)_down\.sql`)
	migrationIDRegex   = regexp.MustCompile(`(\d+)_([\w-]+)`)
)

// Migration is a SQL migration
type Migration struct {
	ID       string
	UpPath   string
	DownPath string

	active bool
}

// AddPath adds or updates a path for a migration direction.
func (m *Migration) AddPath(path string) {
	if direction(path) == "up" {
		m.UpPath = path
	} else {
		m.DownPath = path
	}
}

// Apply runs the up migration on the database.
func (m Migration) Apply() error {
	// Return early if the migration is already active
	active, err := db.MigrationActive(m.ID)
	if err != nil {
		return err
	}
	if active {
		return nil
	}

	sql, err := FS.ReadFile(m.UpPath)
	if err != nil {
		return err
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(string(sql))
	if err != nil {
		log.Printf("[Error] Unable to apply migration (%s): %v", m.ID, err)
		if err := tx.Rollback(); err != nil {
			log.Printf("[Error] Unable to roll back transaction: %v", err)
			return err
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("[Error] Unable to commit transaction: %v", err)
		return err
	}

	// Update the migration log
	err = db.InsertMigration(m.ID)
	if err != nil {
		return err
	}

	fmt.Printf("Migration(%s) applied\n", m.ID)

	return nil
}

// Revert runs the down migration on the database. True will be returned if the migration was completed.
func (m Migration) Revert() (bool, error) {
	// Return early if the migration isn't active
	active, err := db.MigrationActive(m.ID)
	if err != nil {
		return false, err
	}
	if active == false {
		return false, nil
	}

	sql, err := FS.ReadFile(m.DownPath)
	if err != nil {
		return false, err
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		return false, err
	}

	_, err = tx.Exec(string(sql))
	if err != nil {
		log.Printf("[Error] Unable to revert migration (%s): %v", m.ID, err)
		if err := tx.Rollback(); err != nil {
			log.Printf("[Error] Unable to roll back transaction: %v", err)
			return false, err
		}
		return false, err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("[Error] Unable to commit transaction: %v", err)
		return false, err
	}

	// Update the migration log
	err = db.DeleteMigration(m.ID)
	if err != nil {
		return false, err
	}

	fmt.Printf("Migration (%s) reverted\n", m.ID)

	return true, nil
}

// ApplyAll applies all migrations in chronological order.
func ApplyAll() error {
	err := assertMigrationTable()
	if err != nil {
		return err
	}

	migrations, err := all()
	if err != nil {
		return err
	}

	ordered := SortMigrations(migrations, "asc")

	for _, m := range ordered {
		err = m.Apply()
		if err != nil {
			return err
		}
	}

	return nil
}

// RevertAll reverts all migrations in reverse chronological order.
func RevertAll() error {
	err := assertMigrationTable()
	if err != nil {
		return err
	}

	migrations, err := all()
	if err != nil {
		return err
	}

	ordered := SortMigrations(migrations, "desc")

	for _, m := range ordered {
		_, err := m.Revert()
		if err != nil {
			return err
		}
	}

	return nil
}

// Rollback preforms down migrations for `n` active migrations.
func Rollback(n int) error {
	err := assertMigrationTable()
	if err != nil {
		return err
	}

	migrations, err := all()
	if err != nil {
		return err
	}

	ordered := SortMigrations(migrations, "desc")

	count := 1

	for _, m := range ordered {
		// If the count of performed migrations if greater than the number to rollback, we're done.
		if count > n {
			break
		}
		completed, err := m.Revert()
		if err != nil {
			return err
		}
		// Only increment the counter if the migration was completed.
		if completed {
			count++
		}
	}

	return nil
}

// all returns a slice of migrations from the migration directory.
func all() (map[string]*Migration, error) {
	migrations := map[string]*Migration{}

	dir, err := FS.Open(config.MigrationsPath)
	if err != nil {
		return migrations, err
	}

	files, err := dir.Readdir(0)
	if err != nil {
		return migrations, err
	}

	for _, file := range files {
		if valid(file.Name()) {
			id := migrationID(file.Name())
			if _, ok := migrations[id]; !ok {
				migrations[id] = &Migration{
					ID: id,
				}
			}
			m := migrations[id]
			m.AddPath(path.Join(config.MigrationsPath, file.Name()))
		}
	}

	return migrations, nil
}

// id returns the migration ID for a migration file
func migrationID(filename string) string {
	i := strings.LastIndex(filename, "_")
	return filename[0:i]
}

func direction(filename string) string {
	i := strings.LastIndex(filename, "_")
	j := strings.LastIndex(filename, ".")
	return filename[i+1 : j]
}

// valid validates the migration filename
func valid(filename string) bool {
	if upMigrationRegex.MatchString(filename) || downMigrationRegex.MatchString(filename) {
		return true
	}
	return false
}

// assertMigrationTable ensures that the migration table is present in the database.
func assertMigrationTable() error {
	if db.MigrationsTablePresent() {
		return nil
	}

	err := db.CreateMigrationsTable()
	if err != nil {
		return err
	}

	return nil
}
