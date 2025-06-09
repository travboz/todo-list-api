package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/travboz/backend-projects/todo-list-api/internal/db"
	"github.com/travboz/backend-projects/todo-list-api/internal/env"
	"github.com/travboz/backend-projects/todo-list-api/internal/store/mongo"
)

func main() {
	if err := env.LoadEnv(); err != nil {
		log.Fatal("Error loading .env file")
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	mongo_uri := env.GetString(
		"MONGODB_URI",
		"mongodb://travis:secret@localhost:27000/todo-list-api?authSource=admin&readPreference=primary&appname=MongDB%20Compass&directConnection=true&ssl=false",
	)

	mongoClient, err := db.NewMongoDBClient(mongo_uri)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("mongodb successfully connected")

	mongo, err := mongo.NewMongoStore(mongoClient, env.GetString("MONGO_DB_NAME", "todo-list-api"))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := &application{
		Logger:  logger,
		Storage: mongo,
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
