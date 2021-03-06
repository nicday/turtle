package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/nicday/turtle/db"
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
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "Creates the database on the host",
			Action: func(c *cli.Context) {
				migration.CreateDB()
			},
		},
		cli.Command{
			Name:    "drop",
			Aliases: []string{"c"},
			Usage:   "Drops the database on the host",
			Action: func(c *cli.Context) {
				migration.DropDB()
			},
		},
		cli.Command{
			Name:    "up",
			Aliases: []string{"u"},
			Usage:   "Processes all outstanding migrations",
			Action: func(c *cli.Context) {
				db.UseDB()
				migration.ApplyAll()
			},
		},
		cli.Command{
			Name:    "down",
			Aliases: []string{"d"},
			Usage:   "Reverts all applied migrations",
			Action: func(c *cli.Context) {
				db.UseDB()
				migration.RevertAll()
			},
		},
		cli.Command{
			Name:    "rollback",
			Aliases: []string{"r"},
			Usage:   "Rollback n active migrations",
			Action: func(c *cli.Context) {
				if len(c.Args()) == 0 {
					fmt.Println("Please call with a number of migrations to rollback, e.g. `turtle rollback 3`")
					return
				}
				if len(c.Args()) != 0 {
					n, err := strconv.Atoi(c.Args()[0])
					if err != nil {
						log.Fatal("[Error] Rollback parameter is not an integer")
					}
					db.UseDB()
					migration.Rollback(n)
				}
			},
		},
	}

	app.Run(os.Args)
}
