package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/travboz/backend-projects/todo-list-api/internal/env"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
		log.Fatal("Error loading .env file")
	}
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	app := &application{
		Logger: logger,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", env.GetInt("SERVER_PORT", 8000)),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.Logger.Handler(), slog.LevelError),
	}

	app.Logger.Info("server started and running on", "addr", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		app.Logger.Error(err.Error())
		os.Exit(1)
	}
}
