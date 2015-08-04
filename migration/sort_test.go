package migration_test

import (
	. "github.com/nicday/turtle/migration"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sort", func() {
	Describe(".SortMigrations", func() {
		first := &Migration{
			ID: "20150703234300001_first",
		}

		second := &Migration{
			ID: "20150703234300002_second",
		}

		third := &Migration{
			ID: "20150703234300003_third",
		}

		migrations := map[string]*Migration{
			"20150703234300002_second": second,
			"20150703234300001_first":  first,
			"20150703234300003_third":  third,
		}

		Context("with `asc` direction", func() {
			It("sorts ascending", func() {
				expected := []*Migration{
					first,
					second,
					third,
				}
				Expect(SortMigrations(migrations, "asc")).To(Equal(expected))
			})
		})

		Context("with `desc` direction", func() {
			It("sorts descending", func() {
				expected := []*Migration{
					third,
					second,
					first,
				}
				Expect(SortMigrations(migrations, "desc")).To(Equal(expected))
			})
		})

		Context("with an unexpected direction", func() {
			It("sorts ascending", func() {
				expected := []*Migration{
					first,
					second,
					third,
				}
				Expect(SortMigrations(migrations, "unexpected")).To(Equal(expected))
			})
		})
	})
})
