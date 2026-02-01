package health

import (
	"context"
	"net/http"
	"time"

	"brand-protection-monitor/internal/db"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/healthz", HealthzHandler)
	r.GET("/readyz", ReadyzHandler)
}

func HealthzHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func ReadyzHandler(c *gin.Context) {
	pool := db.GetPool()
	if pool == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "DB_UNAVAILABLE",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var id int
	var state string
	err := pool.QueryRow(ctx, "SELECT id, state FROM monitor_state WHERE id = 1").Scan(&id, &state)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "DB_UNAVAILABLE",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":        "ok",
		"monitor_state": state,
	})
}
