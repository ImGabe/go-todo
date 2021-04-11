package cli

import (
	"fmt"
	"strconv"

	"github.com/imgabe/todo/pkg/models"
	"github.com/imgabe/todo/pkg/store"
	"github.com/urfave/cli/v2"
)

// FlagKey is used as a key type on CLI flags
type FlagKey string

// ContextKey is used as a key type on 'context.Context'
type ContextKey string

var (
	// TaskStoreContextKey is the context key used to store the task store
	TaskStoreContextKey ContextKey = "taskstore"
	// TaskStoreContextKey is the context key used to store the database
	DatabaseContextKey ContextKey = "db"
	// FileFlagKey is the flag key used to store the database file path
	FileFlagKey FlagKey = "file"
	// DoneFlagKey is the flag key to decide the visualization of done tasks
	DoneFlagKey FlagKey = "done"
)

// AddTask is responsible for the 'add' command on the CLI
func AddTask(c *cli.Context) error {
	ts := c.Context.Value(TaskStoreContextKey).(store.TaskStore)

	description := c.Args().First()
	task, err := ts.Insert(models.Task{Description: description})
	if err != nil {
		return err
	}

	fmt.Printf("'%s' was added as (%d)\n", task.Description, task.ID)
	return nil
}

// ListTasks is responsible for the 'list' command on the CLI
func ListTasks(c *cli.Context) error {
	ts := c.Context.Value(TaskStoreContextKey).(store.TaskStore)

	tasks, err := ts.SelectAll(c.Bool(string(DoneFlagKey)))
	if err != nil {
		return err
	}

	for _, task := range tasks {
		check := ""
		if task.Done {
			check = "X"
		}

		fmt.Printf("(%d) - [%s] %s\n", task.ID, check, task.Description)
	}

	return nil
}

// RemoveTask is responsible for the 'remove' command on the CLI
func RemoveTask(c *cli.Context) error {
	ts := c.Context.Value(TaskStoreContextKey).(store.TaskStore)

	tmp := c.Args().First()
	taskID, err := strconv.Atoi(tmp)
	if err != nil {
		return err
	}

	err = ts.Delete(int64(taskID))
	if err != nil {
		return err
	}

	fmt.Printf("task (%d) successfully removed\n", taskID)
	return nil
}

// EditTask is responsible for the 'edit' command on the CLI
func EditTask(c *cli.Context) error {
	return nil
}

// ShowTask is responsible for the 'show' command on the CLI
func ShowTask(c *cli.Context) error {
	return nil
}

// Webserver is responsible for the 'web' command on the CLI
func Webserver(c *cli.Context) error {
	return nil
}
