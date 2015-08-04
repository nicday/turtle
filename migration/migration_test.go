package migration_test

import (
	"database/sql"
	"fmt"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/nicday/turtle/db"
	. "github.com/nicday/turtle/migration"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("migration", func() {
	// Setup a mock DB connection
	mockDB, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	db.Conn = mockDB

	// Setup a mock filesystem
	mockFS := NewMockFS()
	mockFS.AddFiles(
		"",
		NewMockFile("migrations", []byte(""),
			NewMockFile("20150703234300001_first_up.sql", []byte("CREATE TABLE first")),
			NewMockFile("20150703234300002_second_up.sql", []byte("CREATE TABLE second")),
			NewMockFile("20150703234300003_third_up.sql", []byte("CREATE TABLE third")),
			NewMockFile("20150703234300001_first_down.sql", []byte("DROP TABLE first")),
			NewMockFile("20150703234300002_second_down.sql", []byte("DROP TABLE second")),
			NewMockFile("20150703234300003_third_down.sql", []byte("DROP TABLE third")),
		),
	)

	FS = mockFS

	db.MigrationsTableName = "migrations"

	Describe("#AddPath", func() {
		Context("with an up migration path", func() {
			It("sets UpPath to the path", func() {
				p := "20150703234300001_first_up.sql"
				m := Migration{}
				m.AddPath(p)

				Expect(m.UpPath).To(Equal(p))
				Expect(m.DownPath).To(Equal(""))
			})
		})

		Context("with an down migration path", func() {
			It("sets DownPath to the path", func() {
				p := "20150703234300001_first_down.sql"
				m := Migration{}
				m.AddPath(p)

				Expect(m.DownPath).To(Equal(p))
				Expect(m.UpPath).To(Equal(""))
			})
		})
	})

	Describe("#Apply", func() {

	})

	Describe("#Revert", func() {

	})

	Describe(".ApplyAll", func() {
		Context("with no active migrations", func() {
			It("applies all migrations", func() {
				ExpectedMigrationActiveQuery("20150703234300001_first", false)
				ExpectedMigration("CREATE TABLE first")
				ExpectedMigrationLogInsert("20150703234300001_first")

				ExpectedMigrationActiveQuery("20150703234300002_second", false)
				ExpectedMigration("CREATE TABLE second")
				ExpectedMigrationLogInsert("20150703234300002_second")

				ExpectedMigrationActiveQuery("20150703234300003_third", false)
				ExpectedMigration("CREATE TABLE third")
				ExpectedMigrationLogInsert("20150703234300003_third")

				err := ApplyAll()

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with some active migrations", func() {
			It("applies all inactive migrations", func() {
				ExpectedMigrationActiveQuery("20150703234300001_first", true)

				ExpectedMigrationActiveQuery("20150703234300002_second", false)
				ExpectedMigration("CREATE TABLE second")
				ExpectedMigrationLogInsert("20150703234300002_second")

				ExpectedMigrationActiveQuery("20150703234300003_third", false)
				ExpectedMigration("CREATE TABLE third")
				ExpectedMigrationLogInsert("20150703234300003_third")

				err := ApplyAll()

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with all active migrations", func() {
			It("doesn't apply any migrations", func() {
				ExpectedMigrationActiveQuery("20150703234300001_first", true)

				ExpectedMigrationActiveQuery("20150703234300002_second", true)

				ExpectedMigrationActiveQuery("20150703234300003_third", true)

				err := ApplyAll()

				Expect(err).NotTo(HaveOccurred())
			})
		})

	})

	Describe(".RevertAll", func() {
		Context("with all active migrations", func() {
			It("reverts all migrations", func() {
				ExpectedMigrationActiveQuery("20150703234300003_third", true)
				ExpectedMigration("DROP TABLE third")
				ExpectedMigrationLogDelete("20150703234300003_third")

				ExpectedMigrationActiveQuery("20150703234300002_second", true)
				ExpectedMigration("DROP TABLE second")
				ExpectedMigrationLogDelete("20150703234300002_second")

				ExpectedMigrationActiveQuery("20150703234300001_first", true)
				ExpectedMigration("DROP TABLE first")
				ExpectedMigrationLogDelete("20150703234300001_first")

				err := RevertAll()

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with some active migrations", func() {
			It("reverts all active migrations", func() {
				ExpectedMigrationActiveQuery("20150703234300003_third", false)

				ExpectedMigrationActiveQuery("20150703234300002_second", false)

				ExpectedMigrationActiveQuery("20150703234300001_first", true)
				ExpectedMigration("DROP TABLE first")
				ExpectedMigrationLogDelete("20150703234300001_first")

				err := RevertAll()

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with all migrations inactive", func() {
			It("doesn't revert any migrations", func() {
				ExpectedMigrationActiveQuery("20150703234300003_third", false)

				ExpectedMigrationActiveQuery("20150703234300002_second", false)

				ExpectedMigrationActiveQuery("20150703234300001_first", false)

				err := RevertAll()

				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})

func ExpectedMigration(sql string) {
	// Migration transaction
	sqlmock.ExpectBegin()
	sqlmock.ExpectExec(regexp.QuoteMeta(sql)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	sqlmock.ExpectCommit()
}

func ExpectedMigrationLogInsert(id string) {
	expectedMigrationLogSQL := fmt.Sprintf(
		"INSERT INTO %s (migration_id) VALUES (?)",
		db.MigrationsTableName,
	)
	sqlmock.ExpectExec(regexp.QuoteMeta(expectedMigrationLogSQL)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func ExpectedMigrationLogDelete(id string) {
	expectedMigrationLogSQL := fmt.Sprintf(
		"DELETE FROM %s WHERE migration_id=?",
		db.MigrationsTableName,
	)
	sqlmock.ExpectExec(regexp.QuoteMeta(expectedMigrationLogSQL)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func ExpectedMigrationActiveQuery(id string, active bool) {
	expectedSQL := fmt.Sprintf(
		"SELECT id FROM %s WHERE migration_id=?",
		db.MigrationsTableName,
	)
	query := sqlmock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs(id)

	if active {
		query.WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	} else {
		query.WillReturnError(sql.ErrNoRows)
	}

}
