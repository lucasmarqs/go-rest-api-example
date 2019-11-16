package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s\n", err))
	}

	db, err := gorm.Open("postgres", viper.Get("database_url"))
	if err != nil {
		panic(fmt.Errorf("Fatal database connection: %s\n", err))
	}
	defer db.Close()

	r := mux.NewRouter()
	app := App{DB: db, R: r}
	app.BuildRoutes()

	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", viper.Get("port")),
		Handler: r,
	}
	go func() {
		fmt.Println("Starting server on port", viper.Get("port"))
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}
