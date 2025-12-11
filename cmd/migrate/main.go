package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"sl-api/config"
)

const (
	dialect     = "pgx"
	fmtDBString = "postgres://%s:%s@%s:%d/%s?sslmode=disable"
)

var (
	flags = flag.NewFlagSet("migrate", flag.ExitOnError)
	dir   = flags.String("dir", "migrations", "directory with migration files")
)

func main() {
	fmt.Printf("Running migrate with the following env vars:\n")
	fmt.Printf("DB_HOST: %q\n", os.Getenv("DB_HOST"))
	fmt.Printf("DB_PORT: %q\n", os.Getenv("DB_PORT"))
	fmt.Printf("DB_USER: %q\n", os.Getenv("DB_USER"))
	fmt.Printf("DB_PASS: %q\n", os.Getenv("DB_PASS"))
	fmt.Printf("DB_NAME: %q\n", os.Getenv("DB_NAME"))
	fmt.Printf("DB_DEBUG: %q\n", os.Getenv("DB_DEBUG"))

	flags.Usage = usage
	flags.Parse(os.Args[1:])

	args := flags.Args()

	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		flags.Usage()

		return
	}

	command := args[0]

	c := config.NewDB()

	dbString := fmt.Sprintf(fmtDBString, c.Username, c.Password, c.Host, c.Port, c.DBName)
	fmt.Printf("Connection string: %s\n", dbString)

	db, err := goose.OpenDBWithDriver(dialect, dbString)
	if err != nil {
		log.Fatalf(err.Error())
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf(err.Error())
		}
	}()

	ctx := context.Background()
	if err := goose.RunContext(ctx, command, db, *dir, args[1:]...); err != nil {
		log.Fatalf("migrate %v: %v", command, err)
	}
}

func usage() {
	fmt.Println(usagePrefix)
	flags.PrintDefaults()
	fmt.Println(usageCommands)
}

var (
	usagePrefix = `Usage: migrate COMMAND
Examples:
    migrate status
`

	usageCommands = `
Commands:
    up                   Migrate the DB to the most recent version available
    up-by-one            Migrate the DB up by 1
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    reset                Roll back all migrations
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with the current timestamp
    fix                  Apply sequential ordering to migrations`
)
