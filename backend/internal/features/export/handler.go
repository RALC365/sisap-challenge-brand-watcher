package export

import (
        "fmt"
        "io"
        "net/http"
        "strconv"
        "time"

        "brand-protection-monitor/internal/features/matches"
        "brand-protection-monitor/internal/observability"

        "github.com/gin-gonic/gin"
        "github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
        service *Service
}

func NewHandler(pool *pgxpool.Pool, matchService *matches.Service) *Handler {
        repo := NewRepository(pool)
        service := NewService(repo, matchService)
        return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router *gin.Engine, limiter *observability.IPRateLimiter) {
        router.GET("/export.csv", observability.RateLimitMiddleware(limiter), h.ExportCSV)
}

func (h *Handler) ExportCSV(c *gin.Context) {
        query, err := parseExportQuery(c)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "EXPORT_ERROR", "message": err.Error()})
                return
        }

        filename := fmt.Sprintf("matches_export_%s.csv", time.Now().Format("20060102_150405"))
        c.Header("Content-Type", "text/csv")
        c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

        c.Stream(func(w io.Writer) bool {
                _, err := h.service.ExportCSV(c.Request.Context(), query, w)
                if err != nil {
                        return false
                }
                return false
        })
}

func parseExportQuery(c *gin.Context) (matches.ListQuery, error) {
        query := matches.ListQuery{
                Keyword: c.Query("keyword"),
                Q:       c.Query("q"),
                Issuer:  c.Query("issuer"),
                NewOnly: c.Query("new_only") == "true",
        }

        query.Page = 1
        query.PageSize = 1000000

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
                query.Sort = matches.SortLastSeenDesc
        case "domain_asc":
                query.Sort = matches.SortDomainAsc
        default:
                query.Sort = matches.SortFirstSeenDesc
        }

        return query, nil
}

var _ = strconv.Atoi
var _ = observability.ExportRateLimiter
