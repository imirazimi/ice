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
	outboxrepo "ice/internal/outbox/repository"
	outboxservice "ice/internal/outbox/service"
	"ice/internal/todo/repository"
	"ice/internal/todo/service"
	"ice/pkg/logger"
	"ice/pkg/migrator"

	"go.uber.org/zap"
)

func main() {

	// Flags
	migrateFlag := flag.Bool("migrate", false, "run DB migrations and exit")
	devFlag := flag.Bool("dev", false, "run in development mode")
	flag.Parse()

	// Logger
	if err := logger.Init(*devFlag); err != nil {
		panic(err)
	}
	defer logger.Sync()

	log := logger.Get()
	cfg := config.Load()

	// Run migrations
	if *migrateFlag {
		if err := migrator.RunMigrations(cfg.MySQL); err != nil {
			log.Fatal("migration failed", zap.Error(err))
		}
		os.Exit(0)
	}

	// Initialize MySQL
	mysqlAdapter, err := mysql.NewMySQL(cfg.MySQL)
	if err != nil {
		log.Fatal("failed to initialize mysql adapter", zap.Error(err))
	}

	// Initialize Redis
	redisCli, err := redis.NewRedisStreamClient(cfg.Redis)
	if err != nil {
		mysqlAdapter.Close()
		log.Fatal("failed to initialize redis adapter", zap.Error(err))
	}
	outboxRepo := outboxrepo.NewRepository(mysqlAdapter)
	outboxService := outboxservice.NewService(outboxRepo, redisCli)
	// Initialize Repository + Service
	TodoRepository := repository.NewRepository(mysqlAdapter)
	todoService := service.NewService(TodoRepository, outboxService)

	// ---------------------------------------
	// NEW: Outbox Processor Context + Goroutine
	// ---------------------------------------
	outboxCtx, outboxCancel := context.WithCancel(context.Background())
	outboxService.StartProcessor(outboxCtx)
	log.Info("Outbox processor started")

	// ---------------------------------------
	// HTTP Server
	// ---------------------------------------
	server := http.NewServer(http.ServerDependencies{
		TodoService: todoService,
		MySQL:       mysqlAdapter.DB(),
		Redis:       redisCli.Client(),
	}, cfg.HTTP.Port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// GLOBAL shutdown context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// ---------------------------------------
	// NEW: Stop Outbox Processor
	// ---------------------------------------
	log.Info("Stopping Outbox processor...")
	outboxCancel()
	// Optional: small sleep to give goroutine time to stop cleanly
	time.Sleep(200 * time.Millisecond)

	// Shutdown HTTP server
	if err := server.Shutdown(ctx); err != nil {
		log.Error("Error shutting down HTTP server", zap.Error(err))
	}

	log.Info("Closing components...")

	// Close Redis
	if err := redisCli.Close(); err != nil {
		log.Error("Error closing Redis connection", zap.Error(err))
	}

	// Close MySQL
	if err := mysqlAdapter.Close(); err != nil {
		log.Error("Error closing MySQL connection", zap.Error(err))
	}

	log.Info("Server exited gracefully")
}
