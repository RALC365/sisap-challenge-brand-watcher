package main

import (
        "context"
        "net/http"
        "os"
        "os/signal"
        "syscall"
        "time"

        "brand-protection-monitor/internal/config"
        "brand-protection-monitor/internal/db"
        "brand-protection-monitor/internal/features/health"
        "brand-protection-monitor/internal/observability"
        "brand-protection-monitor/internal/scheduler"

        "github.com/gin-gonic/gin"
        "go.uber.org/zap"
)

func main() {
        logger := observability.InitLogger()
        defer observability.SyncLogger()

        logger.Info("starting brand protection monitor")

        cfg, err := config.Load()
        if err != nil {
                logger.Fatal("failed to load config", zap.Error(err))
        }
        logger.Info("config loaded", zap.Int("port", cfg.Port))

        ctx := context.Background()
        if err := db.InitPool(ctx, cfg.DatabaseURL); err != nil {
                logger.Fatal("failed to initialize database pool", zap.Error(err))
        }
        defer db.ClosePool()
        logger.Info("database pool initialized")

        gin.SetMode(gin.ReleaseMode)
        router := gin.New()

        router.Use(observability.RequestIDMiddleware())
        router.Use(observability.RecoveryMiddleware(logger))
        router.Use(observability.AccessLogMiddleware(logger))

        health.RegisterRoutes(router)

        matchesGroup := router.Group("/matches")
        matchesGroup.Use(observability.RateLimitMiddleware(observability.MatchesRateLimiter))

        exportGroup := router.Group("/export")
        exportGroup.Use(observability.RateLimitMiddleware(observability.ExportRateLimiter))

        sched := scheduler.New(logger)
        go sched.Start(ctx)
        defer sched.Stop()

        server := &http.Server{
                Addr:         cfg.GetAddr(),
                Handler:      router,
                ReadTimeout:  15 * time.Second,
                WriteTimeout: 15 * time.Second,
                IdleTimeout:  60 * time.Second,
        }

        go func() {
                logger.Info("HTTP server starting", zap.String("addr", cfg.GetAddr()))
                if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
                        logger.Fatal("HTTP server failed", zap.Error(err))
                }
        }()

        quit := make(chan os.Signal, 1)
        signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
        <-quit

        logger.Info("shutting down server...")

        shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()

        if err := server.Shutdown(shutdownCtx); err != nil {
                logger.Error("server shutdown error", zap.Error(err))
        }

        logger.Info("server exited gracefully")
}
