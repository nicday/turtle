package config_test

import (
	"os"

	. "github.com/nicday/turtle/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("config", func() {
	BeforeEach(func() {
		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_NAME", "test")
		os.Setenv("ENV", "")
	})

	Describe(".InitEnv", func() {
		Context("with all required env vars", func() {
			It("returns without error", func() {
				os.Setenv("DB_HOST", "localhost")
				os.Setenv("DB_NAME", "test")

				err := InitEnv()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with missing env vars", func() {
			It("returns an error", func() {
				os.Setenv("DB_HOST", "")

				err := InitEnv()
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(ErrNoDBHost))
			})
			It("returns an error", func() {
				os.Setenv("DB_NAME", "")

				err := InitEnv()
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(ErrNoDBName))
			})
		})

		Context("with an unknown database driver", func() {
			It("returns an error", func() {
				os.Setenv("DB_DRIVER", "somethingelse")

				err := InitEnv()
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(ErrUnknownDBDriver))
			})
		})
	})

	Describe(".IsTestEnv", func() {
		It("returns true with ENV=test", func() {
			os.Setenv("ENV", "test")
			actual := IsTestEnv()

			Expect(actual).To(BeTrue())
		})
		It("returns false when ENV is not `test`", func() {
			actual := IsTestEnv()

			Expect(actual).To(BeFalse())
		})
	})
})
