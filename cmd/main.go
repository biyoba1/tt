package main

import (
	_ "github.com/lib/pq"
	"log"
	"os"
	mdk "valera"
	"valera/initializer"
	"valera/internal/handler"
	"valera/internal/repository"
	"valera/internal/service"
)

func init() {
	initializer.LoadEnvVariables()
	initializer.PingDatabase()
}

func main() {
	postgres, err := repository.NewPostgresDB(repository.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		log.Fatalf("Failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(postgres)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	srv := new(mdk.Server)
	port := os.Getenv("PORT")
	log.Printf("Server is running on port %s", port)
	if err := srv.Run(port, handlers.InitRoutes()); err != nil {
		log.Fatalf("error running a server: %s", err.Error())
	}
}
