package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

func Help() {
	fmt.Print(`task is a CLI for managing your TODOs.

Usage:
  task [command]

Available Commands:
  add         Add a new task to your TODO list
  rm          Remove a task from your TODO list.
  do          Mark a task on your TODO list as complete
  list        List all of your incomplete tasks
  completed   List all of your completed tasks

Use "task [command] --help" for more information about a command.

`,
	)
}

type Task struct {
	Name string
}

type CompletedTask struct {
	Name           string
	CompletionDate time.Time
}

func AddTask(name string, db *bolt.DB) error {
	if strings.Trim(name, " ") == "" {
		return errors.New("task name cannot be empty")
	}

	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("todo"))
		if err != nil {
			return err
		}
		task := Task{Name: name}
		taskBytes, err := json.Marshal(task)
		if err != nil {
			return err
		}
		return b.Put([]byte(name), taskBytes)
	})
}

func ListTasks(db *bolt.DB) ([]Task, error) {
	tasks := []Task{}
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("todo"))
		if b == nil {
			return nil
		}

		return b.ForEach(func(k, v []byte) error {
			var task Task
			if err := json.Unmarshal(v, &task); err != nil {
				return err
			}
			tasks = append(tasks, task)
			return nil
		})
	})

	if err != nil {
		return []Task{}, err
	}

	return tasks, nil
}

func ListCompletedTasks(db *bolt.DB) ([]CompletedTask, error) {
	tasks := []CompletedTask{}
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("completed"))
		if b == nil {
			return nil
		}

		return b.ForEach(func(k, v []byte) error {
			var task CompletedTask
			if err := json.Unmarshal(v, &task); err != nil {
				return err
			}
			if task.CompletionDate.Compare(time.Now().Add(-24*time.Hour)) == 1 {
				tasks = append(tasks, task)
			}
			return nil
		})
	})

	if err != nil {
		return []CompletedTask{}, err
	}

	return tasks, nil
}

func CompleteTask(index int, db *bolt.DB) (Task, error) {
	var task Task

	err := db.Update(func(tx *bolt.Tx) error {
		todoBucket, err := tx.CreateBucketIfNotExists([]byte("todo"))
		if err != nil {
			return err
		}
		completedBucket, err := tx.CreateBucketIfNotExists([]byte("completed"))
		if err != nil {
			return err
		}

		tasks, err := ListTasks(db)
		if err != nil {
			return err
		}

		task = tasks[index]

		if err := todoBucket.Delete([]byte(task.Name)); err != nil {
			return err
		}

		completedTask := CompletedTask{Name: task.Name, CompletionDate: time.Now()}
		completedTaskBytes, err := json.Marshal(completedTask)
		if err != nil {
			return err
		}
		return completedBucket.Put([]byte(task.Name), completedTaskBytes)
	})

	return task, err
}

func DeleteTask(index int, db *bolt.DB) (Task, error) {
	var task Task

	err := db.Update(func(tx *bolt.Tx) error {
		todoBucket, err := tx.CreateBucketIfNotExists([]byte("todo"))
		if err != nil {
			return err
		}
		tasks, err := ListTasks(db)
		task = tasks[index]
		if err != nil {
			return err
		}
		return todoBucket.Delete([]byte(task.Name))
	})

	return task, err
}

func main() {
	db, err := bolt.Open("tasks.db", 0600, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	flag.Usage = Help
	flag.Parse()
	if len(os.Args) == 1 {
		Help()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "add":
		if flag.NArg() < 2 {
			fmt.Println("Please provide a task to add")
			os.Exit(1)
		}

		taskName := strings.Join(os.Args[2:], " ")
		if err := AddTask(taskName, db); err != nil {
			fmt.Println("could not add task:", err)
			os.Exit(1)
		}
		fmt.Println("Added \"" + taskName + "\" to your task list.")

	case "rm":
		if flag.NArg() < 2 {
			fmt.Println("Please provide a task index")
			os.Exit(1)
		}

		taskIndex, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("cannot remove task:", err)
			os.Exit(1)
		}

		task, err := DeleteTask(taskIndex-1, db)
		if err != nil {
			fmt.Println("cannot remove task:", err)
			os.Exit(1)
		}
		fmt.Println("Removed \"" + task.Name + "\" from your task list.")

	case "do":
		if flag.NArg() < 2 {
			fmt.Println("Please provide a task to add")
			os.Exit(1)
		}

		taskIndex, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("cannot complete task:", err)
			os.Exit(1)
		}

		task, err := CompleteTask(taskIndex-1, db)
		if err != nil {
			fmt.Println("cannot complete task:", err)
		}

		fmt.Printf("You have completed \"%s\" task.\n", task.Name)

	case "list":
		tasks, err := ListTasks(db)
		if err != nil {
			fmt.Println("cannot list tasks:", err)
			os.Exit(1)
		}
		fmt.Println("You have the following tasks:")
		for idx, task := range tasks {
			fmt.Printf("%d. %s\n", idx+1, task.Name)
		}

	case "completed":
		tasks, err := ListCompletedTasks(db)
		if err != nil {
			fmt.Println("cannot list tasks:", err)
			os.Exit(1)
		}
		fmt.Println("You have the following tasks:")
		for _, task := range tasks {
			fmt.Printf("- %s\n", task.Name)
		}

	default:
		fmt.Printf("Unknown subcommand: %s\n", os.Args[1])
		os.Exit(1)
	}
}
