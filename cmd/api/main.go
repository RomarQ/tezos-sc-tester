package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	API "github.com/romarq/visualtez-testing/internal/api"
	LOG "github.com/romarq/visualtez-testing/internal/logger"

	_ "github.com/romarq/visualtez-testing/docs"
)

// InitializeAPI - Initialize REST API
// @title Visualtez Testing API
// @version 1.0
// @description API documentation
// @BasePath /
func main() {
	configuration := GetConfig()
	LOG.SetupLogger(configuration.Log.Location, configuration.Log.Level)

	LOG.Info("Initializing API...")

	e := echo.New()

	e.Use(middleware.CORS())
	// Limit body size to 2MB
	e.Use(middleware.BodyLimit("2M"))
	// Rate limit
	rateLimit := middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{
				Rate:      3,
				Burst:     5,
				ExpiresIn: time.Minute,
			},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			return ctx.RealIP(), nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, nil)
		},
	})

	testingAPI := API.InitTestingAPI()

	// API Documentation
	e.GET("/doc/*", echoSwagger.WrapHandler)

	// API Endpoints
	e.POST("/testing", testingAPI.RunTest, rateLimit)

	// Start REST API Service
	go func() {
		if err := e.Start(":" + GetConfig().Port); err != nil && err != http.ErrServerClosed {
			LOG.Fatal("Shutting down REST API service: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Using a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	// Wait for the signal
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		LOG.Fatal("Error during shutdown: %v", err)
	}
}
