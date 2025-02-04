package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/DmitriyGiryntsev/TODO-API/cmd/docs"
	"github.com/DmitriyGiryntsev/TODO-API/internal/config"
	"github.com/DmitriyGiryntsev/TODO-API/internal/db"
	"github.com/DmitriyGiryntsev/TODO-API/internal/handlers"
	"github.com/DmitriyGiryntsev/TODO-API/internal/repository"
	"github.com/DmitriyGiryntsev/TODO-API/internal/routes"
	"github.com/DmitriyGiryntsev/TODO-API/migrations"
	"github.com/gin-gonic/gin"
)

// @title TODO API
// @version 1.0
// @description API для управления задачами
// @host localhost:8080
// @BasePath /api/v1
func main() {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatal("cannot load env variables:", err)
	}

	migrations.ApplyMigrations(cfg.DBURL)

	database, err := db.ConnectPostgres(cfg.DBURL)
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}
	defer database.Close()

	//init repositories
	userRepo := repository.NewUserRepository(database)
	taskRepo := repository.NewTaskRepository(database)

	//init handlers
	authHandler := handlers.NewAuthHandler(userRepo)
	taskHendler := handlers.NewTaskHandler(taskRepo)

	//init server
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	//setup routes
	routes.SetupRoutes(router, authHandler, taskHendler)

	//start server
	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: router,
	}

	go func() {
		log.Printf("server is running on %s", cfg.ServerAddress)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown:", err)
	}

	log.Println("server exited")
}
