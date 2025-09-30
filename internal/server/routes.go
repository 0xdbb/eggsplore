package server

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
		s.gameRoutes(api)
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "up"})
		})
	}
}

func (s *Server) Cors() {
	s.engine.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://192.168.1.201:3000",
			"http://localhost:3001",
			"https://eggsplore.netlify.app", // Add Hoppscotch origin for testing
		},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}

// settingsRoute configures settings routes

func (s *Server) authRoutes(group *gin.RouterGroup) {
	auth := group.Group("/auth")
	{
		auth.POST("/login", s.Login)
		auth.POST("/register", s.Register)
		auth.POST("/logout", s.Logout)
		auth.POST("/renew", s.RenewAccessToken)
	}
}

func (s *Server) gameRoutes(group *gin.RouterGroup) {
	game := group.Group("/game")
	{
		game.GET("/eggs", s.GetPlayerEggs)
		game.POST("/eggs", s.DropEgg)
		game.GET("/inventory", s.GetPlayerInventory)
		game.GET("/player", s.GetPlayerStats)
		game.GET("/tools", s.GetPlayerTools)
	}
}

func (s *Server) swaggerRoute(group *gin.RouterGroup) {
	group.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	group.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusPermanentRedirect, "/api/v1/swagger/index.html")
	})
}
