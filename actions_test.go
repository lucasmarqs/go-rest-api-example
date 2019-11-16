package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

var app App
var dbNotConnected = false

func TestMain(m *testing.M) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s\n", err))
	}
	cfg := viper.Sub("test")

	db, err := gorm.Open("postgres", cfg.Get("database_url"))
	if err != nil {
		log.Println("couldn't connect to testing database", err)
		dbNotConnected = true
	}
	if !dbNotConnected {
		defer db.Close()
		db = db.Begin()
	}

	app = App{DB: db, R: mux.NewRouter()}
	app.BuildRoutes()
	code := m.Run()
	db.RollbackUnlessCommitted()
	os.Exit(code)
}

func TestTasksIndexGet(t *testing.T) {
	if dbNotConnected {
		t.Skip("Testing database must be connected")
		return
	}

	task := Task{Title: "Testing", Status: TaskStatus(1)}
	app.DB.Create(&task)

	req, _ := http.NewRequest("GET", "/tasks", nil)
	rr := httptest.NewRecorder()
	app.R.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Errorf("Expected 200 status code but got %d\n", rr.Code)
	}

	tasks := []Task{}
	if err := json.NewDecoder(rr.Body).Decode(&tasks); err != nil {
		t.Errorf("Expected body to have a valid slice of Task but got %s", rr.Body.String())
	}

	if len(tasks) != 1 {
		t.Errorf("Expected body to have 1 task but got %d", len(tasks))
	}
}

func TestTasksCreatePost(t *testing.T) {
	if dbNotConnected {
		t.Skip("Testing database must be connected")
		return
	}
}
