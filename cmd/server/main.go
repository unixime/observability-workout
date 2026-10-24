package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"condominium-manager/internal/handler"
	"condominium-manager/internal/model"
	"condominium-manager/internal/repository"
	"condominium-manager/internal/telemetry"
)

const serviceName = "condominium-manager"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	shutdownTelemetry, err := telemetry.Setup(ctx, serviceName)
	if err != nil {
		fatal(slog.Default(), "setup telemetry", err)
	}
	logger := telemetry.NewLogger(serviceName)
	slog.SetDefault(logger)

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=condominium port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		fatal(logger, "connect database", err)
	}
	if err := db.AutoMigrate(&model.Tenant{}); err != nil {
		fatal(logger, "migrate", err)
	}

	r := gin.New()
	// otelgin must precede the request logger so log records carry the trace ID
	r.Use(gin.Recovery(), otelgin.Middleware(serviceName), handler.RequestLogger(logger))
	handler.Register(r, repository.NewTenantRepository(db))

	specPath := os.Getenv("SPEC_PATH")
	if specPath == "" {
		specPath = "api/spec.yml"
	}
	handler.RegisterSpec(r, specPath)

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	// Serve until SIGINT/SIGTERM, then flush telemetry before exiting.
	srv := &http.Server{Addr: addr, Handler: r, ReadHeaderTimeout: 10 * time.Second}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fatal(logger, "server stopped", err)
		}
	}()
	logger.Info("server started", slog.String("addr", addr))

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown", slog.Any("error", err))
	}
	if err := shutdownTelemetry(shutdownCtx); err != nil {
		logger.Error("telemetry shutdown", slog.Any("error", err))
	}
}

func fatal(l *slog.Logger, msg string, err error) {
	l.Error(msg, slog.Any("error", err))
	os.Exit(1)
}
