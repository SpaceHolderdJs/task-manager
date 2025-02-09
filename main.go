package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
)

const fileName string = "tasks.json"

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	Priority  string `json:"priority"`
}

type TaskManager struct {
	Tasks  []Task
	nextID int
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		Tasks:  make([]Task, 0),
		nextID: 1,
	}
}

func (tm *TaskManager) AddTask(title string, priority string) {
	task := Task{
		ID:        tm.nextID,
		Title:     title,
		Completed: false,
		Priority:  priority,
	}
	tm.Tasks = append(tm.Tasks, task)
	tm.nextID++
}

func (tm *TaskManager) DeleteTask(title string) {
	tm.Tasks = slices.DeleteFunc(tm.Tasks, func(e Task) bool {
		return e.Title == title
	})
	tm.nextID--
}

func (tm *TaskManager) ListTasks() {
	if len(tm.Tasks) == 0 {
		fmt.Println("No tasks found")
		return
	}

	for _, task := range tm.Tasks {
		status := "☐"
		if task.Completed {
			status = "✓"
		}
		fmt.Printf("%d. [%s] %s (%s)\n", task.ID, status, task.Title, task.Priority)
	}
}

func (tm *TaskManager) SaveToFile() error {
	data, err := json.MarshalIndent(tm, "", "   ")
	if err != nil {
		return fmt.Errorf("error marshaling tasks: %v", err)
	}

	err = os.WriteFile(fileName, data, 0644)
	if err != nil {
		return fmt.Errorf("error saving the task %v", err)
	}

	return nil
}

// Func with multi return
func LoadTaskManager() (*TaskManager, error) {
	data, err := os.ReadFile(fileName)
	if os.IsNotExist(err) {
		// If file doesn't exist, return new TaskManager
		return NewTaskManager(), nil
	}
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	var tm TaskManager
	if err := json.Unmarshal(data, &tm); err != nil {
		return nil, fmt.Errorf("error unmarshaling tasks: %v", err)
	}

	return &tm, nil
}

func main() {
	tm, err := LoadTaskManager()
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		return
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  add <task-name> <priority>  - Add a new task")
		fmt.Println("  list                        - List all tasks")
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 4 {
			fmt.Println("Please provide a task name and priority")
			return
		}

		tm.AddTask(os.Args[2], os.Args[3])

		if err := tm.SaveToFile(); err != nil {
			fmt.Printf("Error saving tasks: %v\n", err)
			return
		}
		fmt.Println("Task added successfully!")

	case "delete":
		tm.DeleteTask(os.Args[2])

		if err := tm.SaveToFile(); err != nil {
			fmt.Printf("Error deleting task(s): %v\n", err)
			return
		}
		fmt.Println("Task deleted successfully!")

	case "list":
		tm.ListTasks()

	default:
		fmt.Println("Unknown command")
	}
}
