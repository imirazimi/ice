package http

import (
	"database/sql"
	"net/http"
	"time"

	_ "ice/docs" // swagger docs
	"ice/internal/port"
	"ice/pkg/logger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
)

type ServerDependencies struct {
	TodoService port.TodoService
	MySQL       *sql.DB
	Redis       *redis.Client
}

func NewServer(deps ServerDependencies, port string) *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(zapLoggerMiddleware())
	e.Use(middleware.Recover())

	// Routes
	todoHandler := NewTodoHandler(deps.TodoService)
	e.POST("/todo", todoHandler.CreateTodo)

	// Health check
	if deps.MySQL != nil && deps.Redis != nil {
		healthChecker := NewHealthChecker(deps.MySQL, deps.Redis)
		e.GET("/health", healthChecker.HealthCheck)
	} else {
		e.GET("/health", healthCheck)
	}

	// Swagger documentation
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Start server in a goroutine
	addr := ":" + port
	go func() {
		log := logger.Get()
		log.Info("Server starting", zap.String("address", addr))
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server crashed", zap.Error(err))
		}
	}()

	return e
}

func zapLoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()

			logger.Get().Info("HTTP request",
				zap.String("method", req.Method),
				zap.String("path", req.URL.Path),
				zap.Int("status", res.Status),
				zap.Duration("latency", time.Since(start)),
				zap.String("ip", c.RealIP()),
			)

			return err
		}
	}
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
