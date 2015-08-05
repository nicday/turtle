package db_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/nicday/turtle/db"

	"testing"
	"time"
)

func TestDb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Db Suite")
}

var _ = BeforeSuite(func() {
	// Reduce the backoff timeout to help avoid long running tests
	BackoffTimeout = time.Duration(1) * time.Second

	mockDBConn()
})

// TODO: Perhaps extract this to a shared test_helpers package
func mockDBConn() {
	mockDB, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	Conn = mockDB
}
