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
                c.JSON(http.StatusBadRequest, ErrorResponse{
                        Error:   ErrorCodeInvalidQuery,
                        Message: err.Error(),
                })
                return
        }

        response, err := h.service.List(c.Request.Context(), query)
        if err != nil {
                c.JSON(http.StatusInternalServerError, ErrorResponse{Error: ErrorCodeDBError})
                return
        }

        c.JSON(http.StatusOK, response)
}

func parseListQuery(c *gin.Context) (ListQuery, error) {
        query := ListQuery{
                Keyword:    c.Query("keyword"),
                KeywordIDs: c.Query("keyword_ids"),
                Q:          c.Query("q"),
                Search:     c.Query("search"),
                Issuer:     c.Query("issuer"),
                NewOnly:    c.Query("new_only") == "true",
        }

        if pageStr := c.Query("page"); pageStr != "" {
                page, err := strconv.Atoi(pageStr)
                if err != nil || page < 1 {
                        return query, ErrInvalidPage
                }
                query.Page = page
        } else {
                query.Page = 1
        }

        if limitStr := c.Query("limit"); limitStr != "" {
                limit, err := strconv.Atoi(limitStr)
                if err != nil {
                        return query, ErrInvalidPageSize
                }
                if limit < 1 || limit > 100 {
                        limit = 20
                }
                query.PageSize = limit
        } else if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
                pageSize, err := strconv.Atoi(pageSizeStr)
                if err != nil {
                        return query, ErrInvalidPageSize
                }
                if pageSize != 10 && pageSize != 25 && pageSize != 50 {
                        return query, ErrInvalidPageSize
                }
                query.PageSize = pageSize
        } else {
                query.PageSize = 20
        }

        if dateFrom := c.Query("start_date"); dateFrom != "" {
                t, err := parseDate(dateFrom)
                if err != nil {
                        return query, ErrInvalidDateFrom
                }
                query.DateFrom = &t
        } else if dateFrom := c.Query("date_from"); dateFrom != "" {
                t, err := parseDate(dateFrom)
                if err != nil {
                        return query, ErrInvalidDateFrom
                }
                query.DateFrom = &t
        }

        if dateTo := c.Query("end_date"); dateTo != "" {
                t, err := parseDate(dateTo)
                if err != nil {
                        return query, ErrInvalidDateTo
                }
                endOfDay := t.Add(24*time.Hour - time.Nanosecond)
                query.DateTo = &endOfDay
        } else if dateTo := c.Query("date_to"); dateTo != "" {
                t, err := parseDate(dateTo)
                if err != nil {
                        return query, ErrInvalidDateTo
                }
                endOfDay := t.Add(24*time.Hour - time.Nanosecond)
                query.DateTo = &endOfDay
        }

        sortBy := c.Query("sort_by")
        sortOrder := c.Query("sort_order")
        sortStr := c.Query("sort")

        if sortBy != "" {
                switch sortBy {
                case "last_seen_at":
                        if sortOrder == "asc" {
                                query.Sort = SortLastSeenDesc
                        } else {
                                query.Sort = SortLastSeenDesc
                        }
                case "domain_name":
                        query.Sort = SortDomainAsc
                default:
                        query.Sort = SortFirstSeenDesc
                }
        } else {
                switch sortStr {
                case "last_seen_desc":
                        query.Sort = SortLastSeenDesc
                case "domain_asc":
                        query.Sort = SortDomainAsc
                case "first_seen_desc", "":
                        query.Sort = SortFirstSeenDesc
                default:
                        return query, ErrInvalidSort
                }
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

var _ = observability.MatchesRateLimiter
