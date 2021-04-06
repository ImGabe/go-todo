package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/imgabe/todo/database"
)

var (
	addPrt    = flag.String("add", "", "add new task")
	listPrt   = flag.Bool("list", false, "list tasks")
	clearPrt  = flag.Bool("clear", false, "clear tasks")
	deletePrt = flag.Int("delete", 0, "delete task")
	checkPrt  = flag.Int("check", 0, "check task")
)

func main() {
	flag.Parse()

	db, err := database.StartDB("todo.sql")
	if err != nil {
		log.Fatal(err)
	}

	if *addPrt != "" {
		todo, err := db.Add(*addPrt)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s was added (%d)\n", todo.Task, todo.ID)
	}

	if *listPrt {
		todos, err := db.List()
		if err != nil {
			log.Fatal(err)
		}

		for _, todo := range todos {
			check := ""

			if todo.Done {
				check = "X"
			}

			fmt.Printf("(%d) - [%s] %s\n", todo.ID, check, todo.Task)
		}
	}

	if *clearPrt {
		err := db.Clear()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("All TODO's have been removed")
	}

	if *deletePrt != 0 && *deletePrt > 0 {
		todoID, err := db.Delete(*deletePrt)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("To-Do %d was removed\n", todoID)
	}

	if *checkPrt != 0 && *checkPrt > 0 {
		todoID, err := db.Check(*checkPrt)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("To-Do %d was checked\n", todoID)
	}
}
