// @title           github.com/0xdbb/eggsplore API
// @version         1.0
// @description     API documentation for the github.com/0xdbb/eggsplore service
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/0xdbb/eggsplore/docs/swagger" // Import generated Swagger docs
	"github.com/0xdbb/eggsplore/internal/config"
	"github.com/0xdbb/eggsplore/internal/server"
)

func main() {
	// ------- Load Config -------
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Configuration loaded. Running in production=%s, port=%s", cfg.Production, cfg.Port)

	// ------- Initialize Server -------
	_, httpServer, err := server.NewServer(cfg)
	if err != nil {
		panic(fmt.Sprintf("server initialization error: %v", err))
	}

	// ------- Context & Background tasks -------
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start report notifier
	// appServer.StartReportNotifier(ctx)

	// ------- Graceful Shutdown -------

	done := make(chan bool, 1)
	go gracefulShutdown(httpServer, done, cancel)

	log.Printf("------ Server listening on Port %s ------\n", cfg.Port)

	err = httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("HTTP server error: %s", err))
	}

	<-done
	log.Println("Graceful shutdown complete.")
}

func gracefulShutdown(apiServer *http.Server, done chan bool, cancel context.CancelFunc) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	log.Println("🔻 Shutdown signal received.")

	cancel() // stop background tasks

	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelTimeout()
	if err := apiServer.Shutdown(ctxTimeout); err != nil {
		log.Printf("⚠️ Forced to shutdown: %v", err)
	}

	done <- true
}
