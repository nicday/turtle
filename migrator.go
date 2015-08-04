package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

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
					migration.Rollback(n)
				}
			},
		},
	}

	app.Run(os.Args)
}
