package db_test

import (
	"database/sql"
	"fmt"
	"regexp"

	"github.com/nicday/turtle/config"
	. "github.com/nicday/turtle/db"

	sqlmock "github.com/DATA-DOG/go-sqlmock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("db", func() {
	// Setup a mock DB connection
	mockDB, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	Conn = mockDB

	// TODO: default is not actually using the default
	tableNames := map[string]string{
		"the default": "migrations",
		"a custom":    "custom_name",
	}

	for desc, tableName := range tableNames {
		config.MigrationsTableName = tableName

		Context(fmt.Sprintf("with %s table name", desc), func() {
			Describe(".CreateMigrationsTable", func() {
				It("creates the migration table in the database", func() {
					expectedSQL := fmt.Sprintf(
						"CREATE TABLE %s (id INT NOT NULL AUTO_INCREMENT, migration_id VARCHAR(255) NOT NULL UNIQUE, PRIMARY KEY(id))",
						config.MigrationsTableName,
					)
					sqlmock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
						WillReturnResult(sqlmock.NewResult(0, 0))

					err := CreateMigrationsTable()
					Expect(err).NotTo(HaveOccurred())
				})

			})

			Describe(".DropMigrationsTable", func() {
				It("drops the migration table in the database", func() {
					expectedSQL := fmt.Sprintf(
						"DROP TABLE %s",
						config.MigrationsTableName,
					)
					sqlmock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
						WillReturnResult(sqlmock.NewResult(0, 0))

					err := DropMigrationsTable()
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Describe(".InsertMigration", func() {
				It("inserts a migration in the migration table", func() {
					ID := "123"
					expectedSQL := fmt.Sprintf(
						"INSERT INTO %s (migration_id) VALUES (?)",
						config.MigrationsTableName,
					)
					sqlmock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
						WithArgs(ID).
						WillReturnResult(sqlmock.NewResult(0, 0))

					err := InsertMigration(ID)
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Describe(".DeleteMigration", func() {
				It("deletes a migration from the migration table", func() {
					ID := "123"
					expectedSQL := fmt.Sprintf(
						"DELETE FROM %s WHERE migration_id=?",
						config.MigrationsTableName,
					)
					sqlmock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
						WithArgs(ID).
						WillReturnResult(sqlmock.NewResult(0, 1))

					err := DeleteMigration(ID)
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Describe(".MigrationActive", func() {
				Context("when the migration is active", func() {
					It("returns true", func() {
						ID := "123"
						expectedSQL := fmt.Sprintf(
							"SELECT id FROM %s WHERE migration_id=?",
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
						ID := "123"
						expectedSQL := fmt.Sprintf(
							"SELECT id FROM %s WHERE migration_id=?",
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
			})
		})
	}
})
