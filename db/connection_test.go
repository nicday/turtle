package db_test

import (
	"os"

	. "github.com/nicday/turtle/db"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("connection", func() {
	Describe("InitConn", func() {
		It("sets up and opens the database connection", func() {
			// InitConn will initializes a new connection that will work in test. Restore the mock DB connection when this
			// test completes
			defer mockDBConn()

			// Remove `test` from ENV
			setEnvVar("ENV", "")
			setDBEnvVars()

			err := InitConn()

			Expect(err).NotTo(HaveOccurred())
		})

		// TODO: Remove pending state from test once there is support for multiple database drivers.
		PIt("returns an error", func() {
			// InitConn will initializes a new connection that will work in test. Restore the mock DB connection when this
			// test completes
			defer mockDBConn()

			// Remove `test` from ENV
			setEnvVar("ENV", "")

			err := InitConn()

			Expect(err).To(HaveOccurred())
		})
	})

	Describe("VerifyConnection", func() {
		It("doesn't panic when a connection can be established", func() {
			Expect(VerifyConnection).NotTo(Panic())
		})

		It("returns an error", func() {
			// Restore the DB connection when this test completes
			defer mockDBConn()

			// Ensure we have no DB connection
			Conn.Close()

			Expect(VerifyConnection).To(Panic())
		})
	})
})

func setEnvVar(key, value string) {
	err := os.Setenv(key, value)
	if err != nil {
		panic(err)
	}
}

func setDBEnvVars() {
	setEnvVar("DB_HOST", "127.0.0.1")
	setEnvVar("DB_NAME", "test")
	setEnvVar("DB_PASSWORD", "password")
}
