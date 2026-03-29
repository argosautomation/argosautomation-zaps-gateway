package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/redis/go-redis/v9"
)

const Port = ":3000"

var ctx = context.Background()
var rdb *redis.Client
var database *sql.DB

func main() {
	// Connect to Redis
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	rdb = redis.NewClient(&redis.Options{
		Addr: redisURL,
	})

	// Test Redis connection
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("⚠️  Redis not available: %v. Running without caching.", err)
	} else {
		log.Println("✓ Connected to Redis")
	}

	// Initialize PostgreSQL
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	var err error
	database, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer database.Close()

	if err = database.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("✓ Connected to PostgreSQL")

	// Setup Fiber
	app := fiber.New(fiber.Config{
		AppName:                 "Zaps Gateway",
		EnableTrustedProxyCheck: true,
		TrustedProxies:          strings.Split(func() string {
			if p := os.Getenv("TRUSTED_PROXIES"); p != "" {
				return p
			}
			return "127.0.0.1"
		}(), ","),
		ProxyHeader:             fiber.HeaderXForwardedFor,
		BodyLimit:               10 * 1024 * 1024, // 10 MB max request body
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{"error": err.Error()})
		},
	})

	// Middleware
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:3000,http://localhost:3001"
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, x-client-id",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	// Security Headers
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		c.Set("X-XSS-Protection", "0")
		c.Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; object-src 'none'")
		return c.Next()
	})

	// Health Check
	app.Get("/health", func(c *fiber.Ctx) error {
		dbHealthy := database.Ping() == nil
		redisHealthy := rdb.Ping(ctx).Err() == nil

		status := "healthy"
		code := 200
		if !dbHealthy || !redisHealthy {
			status = "degraded"
			code = 503
		}

		return c.Status(code).JSON(fiber.Map{
			"status":   status,
			"database": dbHealthy,
			"redis":    redisHealthy,
			"version":  "2.0.0",
		})
	})

	// ==============================================
	// GATEWAY PROXY (PII Redaction)
	// ==============================================

	// For self-hosted deployment, implement API key auth middleware
	// and register the proxy routes here. See docs/api.md for details.

	// Example:
	// apiProxy := app.Group("/v1")
	// apiProxy.Use(APIKeyAuthMiddleware(database))
	// apiProxy.Post("/chat/completions", HandleChatCompletion(rdb))
	// apiProxy.Get("/models", HandleListModels(rdb))

	log.Printf("🚀 Zaps Gateway starting on %s", Port)
	log.Fatal(app.Listen(Port))
}
