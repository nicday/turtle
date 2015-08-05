package db_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

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
})
