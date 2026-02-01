package export

import (
        "fmt"
        "io"
        "net/http"
        "time"

        "brand-protection-monitor/internal/features/matches"
        "brand-protection-monitor/internal/observability"

        "github.com/gin-gonic/gin"
        "github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
        service *Service
        limiter *observability.IPRateLimiter
}

func NewHandler(pool *pgxpool.Pool, matchService *matches.Service, limiter *observability.IPRateLimiter) *Handler {
        repo := NewRepository(pool)
        service := NewService(repo, matchService)
        return &Handler{
                service: service,
                limiter: limiter,
        }
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
        router.GET("/export.csv", h.ExportCSV)
}

func (h *Handler) ExportCSV(c *gin.Context) {
        ip := c.ClientIP()
        if !h.limiter.GetLimiter(ip).Allow() {
                c.JSON(http.StatusTooManyRequests, ErrorResponse{
                        Error:   ErrorCodeRateLimited,
                        Message: "export rate limit exceeded, please try again later",
                })
                return
        }

        query, err := parseExportQuery(c)
        if err != nil {
                c.JSON(http.StatusBadRequest, ErrorResponse{
                        Error:   ErrorCodeExportError,
                        Message: err.Error(),
                })
                return
        }

        filename := fmt.Sprintf("matches_export_%s.csv", time.Now().Format("20060102_150405"))
        c.Header("Content-Type", "text/csv; charset=utf-8")
        c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
        c.Header("Cache-Control", "no-cache, no-store, must-revalidate")

        c.Stream(func(w io.Writer) bool {
                _, err := h.service.ExportCSV(c.Request.Context(), query, c.Writer)
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
                t, err := parseDate(dateFrom)
                if err != nil {
                        return query, matches.ErrInvalidDateFrom
                }
                query.DateFrom = &t
        }

        if dateTo := c.Query("date_to"); dateTo != "" {
                t, err := parseDate(dateTo)
                if err != nil {
                        return query, matches.ErrInvalidDateTo
                }
                endOfDay := t.Add(24*time.Hour - time.Nanosecond)
                query.DateTo = &endOfDay
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

func parseDate(s string) (time.Time, error) {
        t, err := time.Parse(time.RFC3339, s)
        if err == nil {
                return t, nil
        }

        t, err = time.Parse("2006-01-02", s)
        if err == nil {
                return t, nil
        }

        return time.Time{}, err
}
