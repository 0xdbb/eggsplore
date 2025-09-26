package server

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/0xdbb/eggsplore/internal/config"
	"github.com/0xdbb/eggsplore/token"

	db "github.com/0xdbb/eggsplore/internal/database/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	engine *gin.Engine

	tokenMaker  token.Maker
	redisClient *redis.Client
	db          *db.Service
	config      *config.Config
}

func NewServer(appConfig *config.Config) (*Server, *http.Server, error) {
	// Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:      appConfig.RedisAddr,
		Password:  appConfig.RedisPassword,
		DB:        appConfig.RedisDB,
		TLSConfig: &tls.Config{},
	})

	// Set Gin mode
	mode := gin.DebugMode
	if appConfig.Production == "1" {
		mode = gin.ReleaseMode
	}
	gin.SetMode(mode)

	// JWT token maker
	tokenMaker, err := token.NewJWTMaker(appConfig.TokenSecret)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating token maker: %w", err)
	}

	// Database service
	newService := db.NewService(appConfig.DbUrl)

	// Build our Server struct
	appServer := &Server{
		engine:      gin.Default(),
		config:      appConfig,
		tokenMaker:  tokenMaker,
		redisClient: redisClient,
		db:          newService,
	}

	// Register custom validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("StrongPassword", StrongPassword)
		v.RegisterValidation("ValidUsername", ValidUsername)
	}

	appServer.Cors()
	appServer.RegisterRoutes()

	port := fmt.Sprintf(":%s", appConfig.Port)

	// Standard http.Server
	httpSrv := &http.Server{
		Addr:         port,
		Handler:      appServer.engine,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return appServer, httpSrv, nil
}
