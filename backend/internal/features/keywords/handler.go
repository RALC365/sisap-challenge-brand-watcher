package keywords

import (
	"errors"
	"net/http"
	"strconv"

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
	r.GET("/keywords", h.List)
	r.POST("/keywords", h.Create)
	r.DELETE("/keywords/:keyword_id", h.Delete)
}

func (h *Handler) List(c *gin.Context) {
	query := ListQuery{
		Q: c.Query("q"),
	}

	if pageStr := c.Query("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_QUERY", "message": "page must be a positive integer"})
			return
		}
		query.Page = page
	} else {
		query.Page = 1
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil || (pageSize != 10 && pageSize != 25 && pageSize != 50) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_QUERY", "message": "page_size must be 10, 25, or 50"})
			return
		}
		query.PageSize = pageSize
	} else {
		query.PageSize = 10
	}

	response, err := h.service.List(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB_ERROR"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "VALIDATION_ERROR", "message": "value is required and must be a string"})
		return
	}

	keyword, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrEmptyValue) || errors.Is(err, ErrValueTooLong) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "VALIDATION_ERROR", "message": err.Error()})
			return
		}
		if errors.Is(err, ErrDuplicateKeyword) {
			c.JSON(http.StatusConflict, gin.H{"error": "DUPLICATE_KEYWORD"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB_ERROR"})
		return
	}

	c.JSON(http.StatusCreated, keyword)
}

func (h *Handler) Delete(c *gin.Context) {
	keywordID := c.Param("keyword_id")
	if keywordID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_PATH_PARAM", "message": "keyword_id is required"})
		return
	}

	err := h.service.Delete(c.Request.Context(), keywordID)
	if err != nil {
		if errors.Is(err, ErrKeywordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "NOT_FOUND"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB_ERROR"})
		return
	}

	c.JSON(http.StatusOK, DeleteResponse{OK: true})
}
