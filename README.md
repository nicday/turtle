# Turtle
Sea turtles have migrating vast oceans for centuries and now you can be be power and grace of a turtle migration for your mySQL database.

## Commands
The `generate` command generates a new set of migration files with your chosen migration name. Once the files have been
generated you will need to populate them with your migration SQL.

```sh
turtle generate [name]
```

The `up` command applies all inactive migrations. Migrations that have already been applied are ignored.

```sh
turtle up
```

The `down` command reverts all active migrations. Migrations that are haven't been applied are ignroned.

```sh
turtle down
```

### TODO
- Setup Travis CI
- Ability to revert _n_ migrations
- Ability to revert to migration _x_
- Add PostgreSQL support
- Provide information output on performed migrations
- Allow custom migration path
- Example of using turtle migrations within Go
- Create and update schema file after each performed migration

### Author
Nic Day
