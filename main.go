package main

import (
	"kv-cache/api"
	"kv-cache/cache"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"time"
)

func main() {
	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Initialize cache with configuration
	cacheConfig := cache.CacheConfig{
		MaxMemoryMB:     1500,         // 1.8GB max (leaving some headroom)
		MaxItemAge:      5 * time.Minute,    // Items expire after 5 minutes
		CleanupInterval: time.Minute,  // Check for eviction every minute
	}
	
	cache := cache.NewCache(cacheConfig)
	defer cache.Stop()

	// Initialize handler
	h := api.NewHandler(cache)

	// Routes
	e.POST("/put", h.Put)
	e.GET("/get", h.Get)

	// Start server
	e.Logger.Fatal(e.Start(":7171"))
} 