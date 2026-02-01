package monitor

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	service *Service
}

func NewHandler(pool *pgxpool.Pool) *Handler {
	repo := NewRepository(pool)
	service := NewService(repo)
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/monitor/status", h.GetStatus)
}

func (h *Handler) GetStatus(c *gin.Context) {
	status, err := h.service.GetStatus(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "DB_UNAVAILABLE",
		})
		return
	}

	c.JSON(http.StatusOK, status)
}
