package main

import "time"

type TaskStatus int

const (
	StatusTodo TaskStatus = iota + 1
	StatusInProgress
	StatusCompleted
)

// Task represents a record on `tasks` table. It is a GORM model
type Task struct {
	ID        uint       `json:"id"`
	Title     string     `json:"title"`
	Status    TaskStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`

	Errors []string `json:"-" sql:"-"`
}

// Validate build Errors slice and allow you to check its len to check if something went wrong
func (t *Task) Validate() {
	t.Errors = []string{}

	if t.Title == "" {
		t.Errors = append(t.Errors, "title must be present")
	}

	if t.Status < StatusTodo || t.Status > StatusCompleted {
		t.Errors = append(t.Errors, "status is invalid")
	}
}
