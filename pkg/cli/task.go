package cli

import (
	"fmt"
	"math"
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
	// DatabaseContextKey is the context key used to store the database
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

	max := 1
	for _, task := range tasks {
		if taskID := int(task.ID); taskID > max {
			max = taskID
		}
	}
	width := 1 + int(math.Log10(float64(max)))

	for _, task := range tasks {
		fmt.Printf("%-*d %s %s\n", width, task.ID, check(task.Done), task.Description)
	}

	return nil
}

// Checktask is responsible for the 'check' command on the CLI
func CheckTask(c *cli.Context) error {
	ts := c.Context.Value(TaskStoreContextKey).(store.TaskStore)

	tmp := c.Args().First()
	taskID, err := strconv.Atoi(tmp)
	if err != nil {
		return err
	}

	err = ts.Check(int64(taskID))
	if err != nil {
		return err
	}

	fmt.Printf("task (%d) successfully check\n", taskID)
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
	ts := c.Context.Value(TaskStoreContextKey).(store.TaskStore)

	taskID, err := strconv.ParseInt(c.Args().Get(0), 10, 64)
	if err != nil {
		fmt.Println("Missing ID field or is not a number")
		return err
	}

	taskDescription := c.Args().Get(1)
	if len(taskDescription) == 0 {
		fmt.Println("Missing description field")
		return nil
	}

	tmp := c.Args().Get(2)
	taskDone, err := strconv.ParseBool(tmp)
	if err != nil {
		fmt.Println("Missing done field")
		return err
	}

	newTask, err := ts.Update(models.Task{ID: int64(taskID), Description: taskDescription, Done: taskDone})
	if err != nil {
		return err
	}

	fmt.Printf("task (%d) successfully edited\n", newTask.ID)
	return nil
}

// ShowTask is responsible for the 'show' command on the CLI
func ShowTask(c *cli.Context) error {
	ts := c.Context.Value(TaskStoreContextKey).(store.TaskStore)
	taskID, err := strconv.ParseInt(c.Args().First(), 10, 64)

	tmp := c.Args().First()
	taskID, err := strconv.Atoi(tmp)
	if err != nil {
		return err
	}

	task, err := ts.Select(models.Task{ID: taskID})
	if err != nil {
		return err
	}

	fmt.Printf("%d %s %s\n", task.ID, check(task.Done), task.Description)
	return nil
}

func check(done bool) string {
	if done {
		return "*"
	}
	return " "
}

// Webserver is responsible for the 'web' command on the CLI
func Webserver(c *cli.Context) error {
	return nil
}
