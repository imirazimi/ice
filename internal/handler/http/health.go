package http

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

// HealthCheck checks the health of the service
// @Summary Health check endpoint
// @Description Check the health status of the service, MySQL, and Redis
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Success 503 {object} map[string]string
// @Router /health [get]

type HealthChecker struct {
	mysql *sql.DB
	redis *redis.Client
}

func NewHealthChecker(mysql *sql.DB, redis *redis.Client) *HealthChecker {
	return &HealthChecker{
		mysql: mysql,
		redis: redis,
	}
}

func (h *HealthChecker) HealthCheck(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	status := map[string]string{
		"status": "ok",
	}

	// Check MySQL
	if h.mysql != nil {
		if err := h.mysql.PingContext(ctx); err != nil {
			status["mysql"] = "unhealthy"
			status["status"] = "degraded"
		} else {
			status["mysql"] = "healthy"
		}
	}

	// Check Redis
	if h.redis != nil {
		if err := h.redis.Ping(ctx).Err(); err != nil {
			status["redis"] = "unhealthy"
			status["status"] = "degraded"
		} else {
			status["redis"] = "healthy"
		}
	}

	statusCode := http.StatusOK
	if status["status"] == "degraded" {
		statusCode = http.StatusServiceUnavailable
	}

	return c.JSON(statusCode, status)
}

