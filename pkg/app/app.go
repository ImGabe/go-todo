package app

import (
	"context"
	"log"
	"os"
	"os/user"

	commandLine "github.com/imgabe/todo/pkg/cli"
	"github.com/imgabe/todo/pkg/store"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/urfave/cli/v2"
)

func DBPath() string {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	return user.HomeDir + string(os.PathSeparator) + "database.todo"
}

func OpenDatabase(path string) *sqlx.DB {
	db, err := sqlx.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("error opening database: %s", err)
	}

	db.MustExec(store.CreateDatabaseStatement)

	return db
}

func openCliApp() *cli.App {
	return &cli.App{
		Name:  "todo",
		Usage: "keeps track of todos",
		Authors: []*cli.Author{
			{
				Name:  "Victor Freire",
				Email: "victor@freire.dev.br",
			},
			{
				Name:  "Gabe",
				Email: "",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  string(commandLine.FileFlagKey),
				Value: DBPath(),
				Usage: "database file",
			},
			&cli.BoolFlag{
				Name:  string(commandLine.DoneFlagKey),
				Value: false,
				Usage: "show done tasks",
			},
		},
		Before: func(c *cli.Context) error {
			file := c.String(string(commandLine.FileFlagKey))
			db := OpenDatabase(file)
			if err := db.Ping(); err != nil {
				return err
			}
			c.Context = context.WithValue(c.Context, commandLine.DatabaseContextKey, db)

			ts := store.TaskStore{DB: db}
			c.Context = context.WithValue(c.Context, commandLine.TaskStoreContextKey, ts)
			return nil
		},
		After: func(c *cli.Context) error {
			db := c.Context.Value(commandLine.DatabaseContextKey).(*sqlx.DB)
			err := db.Close()
			return err
		},
		Commands: []*cli.Command{
			{
				Name:   "add",
				Usage:  "adds a new task",
				Action: commandLine.AddTask,
			},
			{
				Name:   "check",
				Usage:  "check a task by ID",
				Action: commandLine.CheckTask,
			},
			{
				Name:   "list",
				Usage:  "lists all tasks",
				Action: commandLine.ListTasks,
			},
			{
				Name:   "edit",
				Usage:  "edit a task",
				Action: commandLine.EditTask,
			},
			{
				Name:   "remove",
				Usage:  "remove a task by ID",
				Action: commandLine.RemoveTask,
			},
			{
				Name:   "show",
				Usage:  "show a task by ID",
				Action: commandLine.ShowTask,
			},
			{
				Name:   "web",
				Usage:  "starts a web server",
				Action: commandLine.Webserver,
			},
		},
	}
}

func Run() {
	err := openCliApp().Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
