// @title Todo Service API
// @version 1.0
// @description A Todo service built with Clean Architecture that manages todo items and publishes them to Redis Stream
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ice/config"
	"ice/internal/adapter/mysql"
	"ice/internal/adapter/redis"
	"ice/internal/handler/http"
	"ice/internal/todo/repository"
	"ice/internal/todo/service"
	"ice/pkg/logger"
	"ice/pkg/migrator"

	"go.uber.org/zap"
)

func main() {
	migrateFlag := flag.Bool("migrate", false, "run DB migrations and exit")
	devFlag := flag.Bool("dev", false, "run in development mode")
	flag.Parse()

	// Initialize logger
	if err := logger.Init(*devFlag); err != nil {
		panic(err)
	}
	defer logger.Sync()

	log := logger.Get()

	cfg := config.Load()

	if *migrateFlag {
		if err := migrator.RunMigrations(cfg.MySQL); err != nil {
			log.Fatal("migration failed", zap.Error(err))
		}
		os.Exit(0)
	}

	// Initialize MySQL adapter
	mysqlAdapter, err := mysql.NewMySQL(cfg.MySQL)
	if err != nil {
		log.Fatal("failed to initialize mysql adapter", zap.Error(err))
	}

	// Initialize Redis adapter
	redisCli, err := redis.NewRedisStreamClient(cfg.Redis)
	if err != nil {
		mysqlAdapter.Close()
		log.Fatal("failed to initialize redis adapter", zap.Error(err))
	}

	// Initialize service
	repo := repository.NewRepository(mysqlAdapter)
	todoService := service.NewTodoService(repo, redisCli)

	// Start HTTP server with dependencies
	server := http.NewServer(http.ServerDependencies{
		TodoService: todoService,
		MySQL:       mysqlAdapter.DB(),
		Redis:       redisCli.Client(),
	}, cfg.HTTP.Port)

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := server.Shutdown(ctx); err != nil {
		log.Error("Error shutting down server", zap.Error(err))
	}

	log.Info("Shutting down components...")

	// Close Redis connection
	if err := redisCli.Close(); err != nil {
		log.Error("Error closing Redis connection", zap.Error(err))
	}

	// Close MySQL connection
	if err := mysqlAdapter.Close(); err != nil {
		log.Error("Error closing MySQL connection", zap.Error(err))
	}

	log.Info("Server exited gracefully")
}
