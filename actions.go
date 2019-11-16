package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type App struct {
	DB *gorm.DB
	R  *mux.Router
}

type taskParams struct {
	Title  string `json:"title"`
	Status int    `json:"status"`
}

func (a App) BuildRoutes() {
	a.R.StrictSlash(false)
	a.R.HandleFunc("/tasks", app.tasksIndexGet).Methods("GET")
	a.R.HandleFunc("/tasks", app.tasksCreatePost).Methods("POST")
	a.R.HandleFunc("/tasks/{id:[0-9]+}", app.tasksUpdatePatch).Methods("PATCH")
}

func (a App) tasksIndexGet(w http.ResponseWriter, r *http.Request) {
	tasks := []Task{}
	if err := a.DB.Find(&tasks).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "fatal error getting tasks: %s", err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&tasks)
}

func (a App) tasksCreatePost(w http.ResponseWriter, r *http.Request) {
	p := taskParams{}
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid body: %s", err)
		return
	}

	task := Task{Title: p.Title, Status: TaskStatus(p.Status)}
	task.Validate()
	if len(task.Errors) > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "invalid parameters: %s", strings.Join(task.Errors, ", "))
		return
	}

	if err := a.DB.Create(&task).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "couldn't create a new task: %s", err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&task)
}

func (a App) tasksUpdatePatch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	task := Task{}
	if a.DB.First(&task, id).RecordNotFound() {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "couldn't find Task with ID %d", id)
		return
	}

	p := taskParams{}
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid body: %s", err)
		return
	}

	if p.Title != "" {
		task.Title = p.Title
	}
	task.Status = TaskStatus(p.Status)

	task.Validate()
	if len(task.Errors) > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "invalid parameters: %s", strings.Join(task.Errors, ", "))
		return
	}

	if err := a.DB.Save(&task).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "couldn't updated the task: %s", err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&task)
}
