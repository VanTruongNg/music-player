package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"notification/configs"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	appCfg := configs.LoadAppConfig()

	app, err := InitializeApp(appCfg)
	if err != nil {
		log.Fatalf("[FATAL] Failed to initialize app: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	addr := fmt.Sprintf(":%s", appCfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: app.Router,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("[INFO] HTTP server starting in %s mode on port %s", appCfg.Env, appCfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("[ERROR] HTTP server error: %v", err)
			cancel()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		log.Println("[INFO] Shutting down gracefully...")
	case <-ctx.Done():
		log.Println("[INFO] Context cancelled, shutting down...")
	}

	cancel()

	log.Println("[INFO] Cleaning up resources...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("[ERROR] HTTP server forced to shutdown: %v", err)
	} else {
		log.Println("[INFO] HTTP server stopped successfully")
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("[INFO] All services stopped successfully")
	case <-time.After(5 * time.Second):
		log.Println("[WARN] Timeout waiting for services to stop")
	}

	log.Println("[INFO] Application shutdown complete")
}
