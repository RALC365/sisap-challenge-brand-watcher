package matches

import (
	"net/http"
	"strconv"
	"time"

	"brand-protection-monitor/internal/observability"

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

func (h *Handler) RegisterRoutes(matchesGroup *gin.RouterGroup) {
	matchesGroup.GET("", h.List)
}

func (h *Handler) List(c *gin.Context) {
	query, err := parseListQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_QUERY", "message": err.Error()})
		return
	}

	response, err := h.service.List(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB_ERROR"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func parseListQuery(c *gin.Context) (ListQuery, error) {
	query := ListQuery{
		Keyword: c.Query("keyword"),
		Q:       c.Query("q"),
		Issuer:  c.Query("issuer"),
		NewOnly: c.Query("new_only") == "true",
	}

	if pageStr := c.Query("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			return query, err
		}
		query.Page = page
	} else {
		query.Page = 1
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil || (pageSize != 10 && pageSize != 25 && pageSize != 50) {
			return query, err
		}
		query.PageSize = pageSize
	} else {
		query.PageSize = 10
	}

	if dateFrom := c.Query("date_from"); dateFrom != "" {
		t, err := time.Parse(time.RFC3339, dateFrom)
		if err != nil {
			t, err = time.Parse("2006-01-02", dateFrom)
			if err != nil {
				return query, err
			}
		}
		query.DateFrom = &t
	}

	if dateTo := c.Query("date_to"); dateTo != "" {
		t, err := time.Parse(time.RFC3339, dateTo)
		if err != nil {
			t, err = time.Parse("2006-01-02", dateTo)
			if err != nil {
				return query, err
			}
		}
		query.DateTo = &t
	}

	sort := c.Query("sort")
	switch sort {
	case "last_seen_desc":
		query.Sort = SortLastSeenDesc
	case "domain_asc":
		query.Sort = SortDomainAsc
	default:
		query.Sort = SortFirstSeenDesc
	}

	return query, nil
}

var _ = observability.MatchesRateLimiter
