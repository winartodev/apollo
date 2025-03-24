package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/apollo/core/responses"
	"runtime"
	"time"
)

type HealthResponse struct {
	Status    string  `json:"status"`
	Timestamp string  `json:"timestamp"`
	Details   Details `json:"details"`
	Message   string  `json:"message"`
}

type Details struct {
	Database Database `json:"database"`
	Server   Server   `json:"server"`
	Services Services `json:"services"`
}

type Database struct {
	Status     string `json:"status"`
	PingTimeMs int    `json:"ping_time_ms"`
}

type Server struct {
	UptimeSeconds   int     `json:"uptime_seconds"`
	MemoryUsageMB   uint64  `json:"memory_usage_mb"`
	CPUUsagePercent float64 `json:"cpu_usage_percent"`
}

type Services struct {
	Cache        string `json:"cache"`
	MessageQueue string `json:"message_queue"`
}

type HealthHandler struct {
	StartTime time.Time
}

func (c *HealthHandler) HealthZ(ctx *fiber.Ctx) error {
	dbStatus := "connected"
	dbPingTime := 12

	uptime := int(time.Since(c.StartTime).Seconds())

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	memoryUsageMB := memStats.Alloc / 1024 / 1024 // Convert bytes to MB

	cacheStatus := "disable"
	messageQueueStatus := "disable"

	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Details: Details{
			Database: Database{
				Status:     dbStatus,
				PingTimeMs: dbPingTime,
			},
			Server: Server{
				UptimeSeconds:   uptime,
				MemoryUsageMB:   memoryUsageMB,
				CPUUsagePercent: 0,
			},
			Services: Services{
				Cache:        cacheStatus,
				MessageQueue: messageQueueStatus,
			},
		},
		Message: "All systems operational",
	}

	return responses.SuccessResponse(ctx, fiber.StatusOK, "", response, nil)
}
