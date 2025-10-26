package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"auth-service/configs"
)

func main() {
	appCfg := configs.LoadAppConfig()
	dbCfg := configs.LoadDBConfig()
	redisCfg := configs.LoadRedisConfig()
	kafkaCfg := configs.LoadKafkaConfig()

	app, err := InitializeApp(appCfg, dbCfg, redisCfg, kafkaCfg)
	if err != nil {
		log.Fatalf("[FATAL] Failed to initialize app: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	// Create HTTP server with graceful shutdown
	addr := fmt.Sprintf(":%s", appCfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: app.Router,
	}

	// Start HTTP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("[INFO] HTTP server starting in %s mode on port %s", appCfg.Env, appCfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("[ERROR] HTTP server error: %v", err)
			cancel()
		}
	}()

	// Start gRPC server
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("[INFO] gRPC server starting on port %s", appCfg.GRPCPort)
		if err := app.GRPCServer.Start(); err != nil {
			log.Printf("[ERROR] gRPC server error: %v", err)
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

	if app.KafkaProducer != nil {
		app.KafkaProducer.Close()
		log.Println("[INFO] Kafka producer closed successfully")
	}

	if app.GRPCServer != nil {
		log.Println("[INFO] Stopping gRPC server...")
		app.GRPCServer.Stop()
		log.Println("[INFO] gRPC server stopped successfully")
	}

	if app.KafkaConsumer != nil {
		app.KafkaConsumer.Close()
		log.Println("[INFO] Kafka consumer closed successfully")
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
