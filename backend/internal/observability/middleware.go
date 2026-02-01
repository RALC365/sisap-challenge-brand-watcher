package observability

import (
        "net/http"
        "sync"
        "time"

        "github.com/gin-gonic/gin"
        "github.com/google/uuid"
        "go.uber.org/zap"
        "golang.org/x/time/rate"
)

const (
        RequestIDHeader = "X-Request-ID"
        RequestIDKey    = "request_id"
)

func RequestIDMiddleware() gin.HandlerFunc {
        return func(c *gin.Context) {
                requestID := c.GetHeader(RequestIDHeader)
                if requestID == "" {
                        requestID = uuid.New().String()
                }
                c.Set(RequestIDKey, requestID)
                c.Header(RequestIDHeader, requestID)
                c.Next()
        }
}

func AccessLogMiddleware(logger *zap.Logger) gin.HandlerFunc {
        return func(c *gin.Context) {
                start := time.Now()
                path := c.Request.URL.Path
                query := c.Request.URL.RawQuery

                c.Next()

                latency := time.Since(start)
                statusCode := c.Writer.Status()

                reqID := GetRequestID(c)

                logger.Info("request",
                        zap.String("request_id", reqID),
                        zap.String("method", c.Request.Method),
                        zap.String("path", path),
                        zap.String("query", query),
                        zap.Int("status", statusCode),
                        zap.Duration("latency", latency),
                        zap.String("client_ip", c.ClientIP()),
                        zap.String("user_agent", c.Request.UserAgent()),
                )
        }
}

func GetRequestID(c *gin.Context) string {
        if requestID, exists := c.Get(RequestIDKey); exists {
                if reqID, ok := requestID.(string); ok {
                        return reqID
                }
        }
        return ""
}

func RecoveryMiddleware(logger *zap.Logger) gin.HandlerFunc {
        return func(c *gin.Context) {
                defer func() {
                        if err := recover(); err != nil {
                                reqID := GetRequestID(c)

                                logger.Error("panic recovered",
                                        zap.String("request_id", reqID),
                                        zap.Any("error", err),
                                        zap.String("path", c.Request.URL.Path),
                                )

                                c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
                                        "error":      "Internal Server Error",
                                        "request_id": reqID,
                                })
                        }
                }()
                c.Next()
        }
}

type IPRateLimiter struct {
        limiters map[string]*rate.Limiter
        mu       sync.RWMutex
        rate     rate.Limit
        burst    int
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
        return &IPRateLimiter{
                limiters: make(map[string]*rate.Limiter),
                rate:     r,
                burst:    b,
        }
}

func (rl *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
        rl.mu.Lock()
        defer rl.mu.Unlock()

        limiter, exists := rl.limiters[ip]
        if !exists {
                limiter = rate.NewLimiter(rl.rate, rl.burst)
                rl.limiters[ip] = limiter
        }

        return limiter
}

func RateLimitMiddleware(limiter *IPRateLimiter) gin.HandlerFunc {
        return func(c *gin.Context) {
                ip := c.ClientIP()
                if !limiter.GetLimiter(ip).Allow() {
                        reqID := GetRequestID(c)

                        c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
                                "error":      "Rate limit exceeded",
                                "request_id": reqID,
                        })
                        return
                }
                c.Next()
        }
}

var (
        MatchesRateLimiter = NewIPRateLimiter(rate.Limit(10), 20)
        ExportRateLimiter  = NewIPRateLimiter(rate.Limit(1), 3)
)
