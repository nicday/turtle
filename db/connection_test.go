package db_test

import (
	"os"

	"github.com/nicday/turtle/config"
	. "github.com/nicday/turtle/db"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var oldEnv = map[string]string{}

var _ = Describe("db", func() {
	BeforeEach(func() {
		OverwriteEnv("ENV", "other")
		OverwriteEnv("DB_HOST", "host")
		OverwriteEnv("DB_PORT", "port")
		OverwriteEnv("DB_USER", "user")
		OverwriteEnv("DB_NAME", "test")
	})

	AfterEach(func() {
		ResetEnv()
	})

	Describe(".ConnString", func() {
		Context("when DB_DRIVER=mysql", func() {
			It("returns a mysql connection string", func() {
				os.Setenv("DB_DRIVER", "mysql")
				config.InitEnv()

				actual := ConnString()
				expected := "user@tcp(host:port)/"
				Expect(actual).To(Equal(expected))
			})
		})

		Context("when DB_DRIVER=postgres", func() {
			It("returns a postgres connection string", func() {
				os.Setenv("DB_DRIVER", "postgres")
				config.InitEnv()

				actual := ConnString()
				expected := "postgres://user@host:port/test?sslmode=disable"
				Expect(actual).To(Equal(expected))
			})
		})
	})
})

func OverwriteEnv(envVar, val string) {
	oldEnv[envVar] = os.Getenv(envVar)
	os.Setenv(envVar, val)
}

func ResetEnv() {
	for k, v := range oldEnv {
		os.Setenv(k, v)
	}
}
