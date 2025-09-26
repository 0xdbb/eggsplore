package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/redis"
)

const (
	AuthRateLimit    = "20-M" // 5 requests per minute for auth routes
	AccountRateLimit = "10-M" // 20 requests per minute for account routes
)

// @BasePath /api/v1
func (s *Server) RegisterRoutes() {
	// Redirect "/" to Swagger
	s.engine.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusPermanentRedirect, "/api/v1/swagger/index.html")
	})

	// Redirect "/api" to Swagger
	s.engine.GET("/api", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusPermanentRedirect, "/api/v1/swagger/index.html")
	})

	// Redirect "/api/v1" to Swagger
	s.engine.GET("/api/v1", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusPermanentRedirect, "/api/v1/swagger/index.html")
	})

	// Register versioned API routes
	api := s.engine.Group("/api/v1")
	{
		s.swaggerRoute(api)
		s.authRoutes(api)
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "up"})
		})
	}
}

func (s *Server) Cors() {
	s.engine.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"https://citizens-reports-portal.netlify.app",
			"http://localhost:3001",
			"https://github.com/0xdbb/eggsplore.blvcksapphire.com",
			"https://github.com/0xdbb/eggsplore-v2.blvcksapphire.com",
			"https://citizensreport.github.com/0xdbb/eggsplore.blvcksapphire.com",
			"https://citizen-report.blvcksapphire.com",
			"https://mdukq3hdsh.us-east-1.awsapprunner.com",
			"https://api.github.com/0xdbb/eggsplore.blvcksapphire.com",
			"https://citizensreport-v2.github.com/0xdbb/eggsplore.blvcksapphire.com",
			"https://hoppscotch.io", // Add Hoppscotch origin for testing
		},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
			"Upgrade",
			"Connection",
			"Sec-WebSocket-Key",
			"Sec-WebSocket-Version",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}

// createLimiter creates a rate limiter from string rate and Redis store
func createLimiter(formattedRate string, store limiter.Store) *limiter.Limiter {
	rate, err := limiter.NewRateFromFormatted(formattedRate)
	if err != nil {
		log.Fatalf("Failed to parse rate %s: %v", formattedRate, err)
	}
	return limiter.New(store, rate, limiter.WithTrustForwardHeader(true))
}

// settingsRoute configures settings routes

func (s *Server) authRoutes(group *gin.RouterGroup) {
	store, err := redis.NewStoreWithOptions(s.redisClient, limiter.StoreOptions{
		Prefix:   "limiter_auth",
		MaxRetry: 3,
	})
	if err != nil {
		log.Fatalf("Failed to create Redis store for auth: %v", err)
	}

	authLimiter := createLimiter(AuthRateLimit, store)

	auth := group.Group("/auth")
	auth.Use(mgin.NewMiddleware(authLimiter))
	{
		auth.POST("/login", s.Login)
		auth.POST("/verify-otp", s.VerifyOTP)
		auth.POST("/send-otp", s.SendOTP)
		auth.POST("/logout", s.Logout)
		auth.POST("/renew", s.RenewAccessToken)
	}
}

func (s *Server) swaggerRoute(group *gin.RouterGroup) {
	group.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	group.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusPermanentRedirect, "/api/v1/swagger/index.html")
	})
}
