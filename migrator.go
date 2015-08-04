package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/nicday/turtle/migration"
)

func main() {
	app := cli.NewApp()
	app.Name = "turtle"
	app.Usage = "for incredible (SQL) migrations, just the sea turtle!"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		cli.Command{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "Generates a new set of migration files",
			Action: func(c *cli.Context) {
				if len(c.Args()) == 0 {
					fmt.Println("Please call with a migration name, e.g. `turtle generate users`")
					return
				}
				if len(c.Args()) != 0 {
					migrationName := c.Args()[0]
					migration.Generate(migrationName)
				}

			},
		},
		cli.Command{
			Name:    "up",
			Aliases: []string{"u"},
			Usage:   "Processes all outstanding migrations",
			Action: func(c *cli.Context) {
				// db.CreateMigrationsTable()
				migration.ApplyAll()
			},
		},
		cli.Command{
			Name:    "down",
			Aliases: []string{"d"},
			Usage:   "Processes all outstanding migrations",
			Action: func(c *cli.Context) {
				migration.RevertAll()
			},
		},
	}

	app.Run(os.Args)
}
