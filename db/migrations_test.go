package db_test

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"

	"github.com/nicday/turtle/config"
	. "github.com/nicday/turtle/db"

	sqlmock "github.com/DATA-DOG/go-sqlmock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("db", func() {
	// TODO: default is not actually using the default
	tableNames := map[string]string{
		"the default": "migrations",
		"a custom":    "custom_name",
	}

	for desc, tableName := range tableNames {
		config.MigrationsTableName = tableName

		Context(fmt.Sprintf("with %s table name", desc), func() {
			Describe(".CreateMigrationsTable", func() {
				sqlFormat := "CREATE TABLE %s (id INT NOT NULL AUTO_INCREMENT, migration_id VARCHAR(255) NOT NULL UNIQUE, PRIMARY KEY(id))"

				It("creates the migration table in the database", func() {
					expectedSQL := fmt.Sprintf(
						sqlFormat,
						config.MigrationsTableName,
					)
					sqlmock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
						WillReturnResult(sqlmock.NewResult(0, 0))

					err := CreateMigrationsTable()
					Expect(err).NotTo(HaveOccurred())
				})

				It("returns an error if query prepare fails", func() {
					// Throw an error when prepare is run
					sqlmock.ExpectPrepare().
						WillReturnError(errors.New("sql: database is closed"))

					err := CreateMigrationsTable()
					Expect(err).To(HaveOccurred())
				})

				It("returns an error if query exec fails", func() {
					// Throw an error when exec is run
					expectedSQL := fmt.Sprintf(
						sqlFormat,
						config.MigrationsTableName,
					)
					sqlmock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
						WillReturnError(errors.New("sql: database is closed"))

					err := CreateMigrationsTable()
					Expect(err).To(HaveOccurred())
				})

			})

			Describe(".DropMigrationsTable", func() {
				sqlFormat := "DROP TABLE %s"

				It("drops the migration table in the database", func() {
					expectedSQL := fmt.Sprintf(
						sqlFormat,
						config.MigrationsTableName,
					)
					sqlmock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
						WillReturnResult(sqlmock.NewResult(0, 0))

					err := DropMigrationsTable()
					Expect(err).NotTo(HaveOccurred())
				})

				It("returns an error if query prepare fails", func() {
					// Throw an error when prepare is run
					sqlmock.ExpectPrepare().
						WillReturnError(errors.New("sql: database is closed"))

					err := DropMigrationsTable()
					Expect(err).To(HaveOccurred())
				})

				It("returns an error if query exec fails", func() {
					// Throw an error when exec is run
					expectedSQL := fmt.Sprintf(
						sqlFormat,
						config.MigrationsTableName,
					)
					sqlmock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
						WillReturnError(errors.New("sql: database is closed"))

					err := DropMigrationsTable()
					Expect(err).To(HaveOccurred())
				})
			})

			Describe(".InsertMigration", func() {
				sqlFormat := "INSERT INTO %s (migration_id) VALUES (?)"
				ID := "123"

				It("inserts a migration in the migration table", func() {
					expectedSQL := fmt.Sprintf(
						sqlFormat,
						config.MigrationsTableName,
					)
					sqlmock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
						WithArgs(ID).
						WillReturnResult(sqlmock.NewResult(0, 0))

					err := InsertMigration(ID)
					Expect(err).NotTo(HaveOccurred())
				})

				It("returns an error if query prepare fails", func() {
					// Throw an error when prepare is run
					sqlmock.ExpectPrepare().
						WillReturnError(errors.New("sql: database is closed"))

					err := InsertMigration(ID)
					Expect(err).To(HaveOccurred())
				})

				It("returns an error if query exec fails", func() {
					// Throw an error when exec is run
					expectedSQL := fmt.Sprintf(
						sqlFormat,
						config.MigrationsTableName,
					)
					sqlmock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
						WithArgs(ID).
						WillReturnError(errors.New("sql: database is closed"))

					err := InsertMigration(ID)
					Expect(err).To(HaveOccurred())
				})
			})

			Describe(".DeleteMigration", func() {
				sqlFormat := "DELETE FROM %s WHERE migration_id=?"
				ID := "123"

				It("deletes a migration from the migration table", func() {
					expectedSQL := fmt.Sprintf(
						sqlFormat,
						config.MigrationsTableName,
					)
					sqlmock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
						WithArgs(ID).
						WillReturnResult(sqlmock.NewResult(0, 1))

					err := DeleteMigration(ID)
					Expect(err).NotTo(HaveOccurred())
				})

				It("returns an error if query prepare fails", func() {
					// Throw an error when prepare is run
					sqlmock.ExpectPrepare().
						WillReturnError(errors.New("sql: database is closed"))

					err := DeleteMigration(ID)
					Expect(err).To(HaveOccurred())
				})

				It("returns an error if query exec fails", func() {
					// Throw an error when exec is run
					expectedSQL := fmt.Sprintf(
						sqlFormat,
						config.MigrationsTableName,
					)
					sqlmock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
						WithArgs(ID).
						WillReturnError(errors.New("sql: database is closed"))

					err := DeleteMigration(ID)
					Expect(err).To(HaveOccurred())
				})
			})

			Describe(".MigrationActive", func() {
				sqlFormat := "SELECT id FROM %s WHERE migration_id=?"
				ID := "123"

				Context("when the migration is active", func() {
					It("returns true", func() {
						expectedSQL := fmt.Sprintf(
							sqlFormat,
							config.MigrationsTableName,
						)
						sqlmock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
							WithArgs(ID).
							WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

						active, err := MigrationActive(ID)

						Expect(active).To(BeTrue())
						Expect(err).NotTo(HaveOccurred())
					})
				})

				Context("when the migration is inactive", func() {
					It("returns false", func() {
						expectedSQL := fmt.Sprintf(
							sqlFormat,
							config.MigrationsTableName,
						)
						sqlmock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
							WithArgs(ID).
							WillReturnError(sql.ErrNoRows)

						active, err := MigrationActive(ID)

						Expect(active).To(BeFalse())
						Expect(err).NotTo(HaveOccurred())
					})
				})

				It("returns an error if query prepare fails", func() {
					// Throw an error when prepare is run
					sqlmock.ExpectPrepare().
						WillReturnError(errors.New("sql: database is closed"))

					active, err := MigrationActive(ID)

					Expect(active).To(BeFalse())
					Expect(err).To(HaveOccurred())
				})

				It("returns an error if query fails", func() {
					// Throw an error when exec is run
					expectedSQL := fmt.Sprintf(
						sqlFormat,
						config.MigrationsTableName,
					)
					sqlmock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
						WithArgs(ID).
						WillReturnError(errors.New("sql: database is closed"))

					active, err := MigrationActive(ID)

					Expect(active).To(BeFalse())
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Describe(".MigrationsTablePresent", func() {
			Context("when migration table is present", func() {
				It("returns true", func() {
					expectMigrationsTablePresenceQuery(true)

					present := MigrationsTablePresent()

					Expect(present).To(BeTrue())
				})
			})

			Context("when migration table is not present", func() {
				It("returns false", func() {
					expectMigrationsTablePresenceQuery(false)

					present := MigrationsTablePresent()

					Expect(present).To(BeFalse())
				})
			})
		})
	}
})

// TODO: This func is used in migration and db tests - extract to a test helpers package.
func expectMigrationsTablePresenceQuery(present bool) {
	expectedSQL := fmt.Sprintf(
		"SELECT 1 FROM %s LIMIT 1",
		config.MigrationsTableName,
	)
	query := sqlmock.ExpectExec(regexp.QuoteMeta(expectedSQL))

	if present {
		query.WillReturnResult(sqlmock.NewResult(0, 0))
	} else {
		query.WillReturnError(sql.ErrNoRows)
	}
}
