package main

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/boltdb/bolt"
)

func TestAddOneTask(t *testing.T) {
	// Arrange
	db, err := bolt.Open("test.db", 0600, nil)
	if err != nil {
		panic(err)
	}
	defer os.Remove("test.db")
	defer db.Close()

	// Act
	err = AddTask("test task", db)
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	tasks, err := ListTasks(db)
	if err != nil {
		t.Fatal(err)
	}

	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}

	if tasks[0].Name != "test task" {
		t.Fatalf("expected task name to be 'test task', got %s", tasks[0].Name)
	}
}

func TestAddMultipleTasks(t *testing.T) {
	// Arrange
	db, err := bolt.Open("test.db", 0600, nil)
	if err != nil {
		panic(err)
	}
	defer os.Remove("test.db")
	defer db.Close()

	// Act
	err = AddTask("test task", db)
	if err != nil {
		t.Fatal(err)
	}
	err = AddTask("test task 2", db)
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	tasks, err := ListTasks(db)
	if err != nil {
		t.Fatal(err)
	}

	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}

	if tasks[0].Name != "test task" {
		t.Fatalf("expected task name to be 'test task', got %s", tasks[0].Name)
	}
	if tasks[1].Name != "test task 2" {
		t.Fatalf("expected task name to be 'test task 2', got %s", tasks[0].Name)
	}
}

func TestCompleteTask(t *testing.T) {
	// Arrange
	db, err := bolt.Open("test.db", 0600, nil)
	if err != nil {
		panic(err)
	}
	defer os.Remove("test.db")
	defer db.Close()

	err = AddTask("test task", db)
	if err != nil {
		t.Fatal(err)
	}

	// Act
	_, err = CompleteTask(0, db)
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	tasks, err := ListTasks(db)
	if err != nil {
		t.Fatal(err)
	}

	if len(tasks) != 0 {
		t.Fatalf("expected 0 tasks, got %d", len(tasks))
	}

	completedTasks, err := ListCompletedTasks(db)
	if err != nil {
		t.Fatal(err)
	}

	if len(completedTasks) != 1 {
		t.Fatalf("expected 1 tasks, got %d", len(completedTasks))
	}
}

func TestListCompletedTasksWhenOneTaskIsOlderThanADay(t *testing.T) {
	// Arrange
	db, err := bolt.Open("test.db", 0600, nil)
	if err != nil {
		panic(err)
	}
	defer os.Remove("test.db")
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("completed"))
		if err != nil {
			return err
		}

		task := CompletedTask{Name: "test task", CompletionDate: time.Now()}
		taskBytes, err := json.Marshal(task)
		if err != nil {
			return err
		}
		if err := b.Put([]byte("test task"), taskBytes); err != nil {
			return err
		}

		oldCompletionDate := time.Now().Add(-48 * time.Hour)
		task = CompletedTask{Name: "test task 2", CompletionDate: oldCompletionDate}
		taskBytes, err = json.Marshal(task)
		if err != nil {
			return err
		}
		return b.Put([]byte("test task 2"), taskBytes)
	})
	if err != nil {
		t.Fatal(err)
	}

	// Act
	completedTasks, err := ListCompletedTasks(db)
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	if len(completedTasks) != 1 {
		t.Fatalf("expected 1 tasks, got %d", len(completedTasks))
	}

	if completedTasks[0].Name != "test task" {
		t.Fatalf("expected task name to be 'test task', got %s", completedTasks[0].Name)
	}

}

func TestDeleteTask(t *testing.T) {
	// Arrange
	db, err := bolt.Open("test.db", 0600, nil)
	if err != nil {
		panic(err)
	}
	defer os.Remove("test.db")
	defer db.Close()

	err = AddTask("test task", db)
	if err != nil {
		t.Fatal(err)
	}
	err = AddTask("test task 2", db)
	if err != nil {
		t.Fatal(err)
	}

	// Act
	_, err = DeleteTask(1, db)
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	tasks, err := ListTasks(db)
	if err != nil {
		t.Fatal(err)
	}

	if len(tasks) != 1 {
		t.Fatalf("expected 1 tasks, got %d", len(tasks))
	}

	if tasks[0].Name != "test task" {
		t.Fatalf("expected task name to be 'test task', got %s", tasks[0].Name)
	}
}
