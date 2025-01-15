package main

import (
	"arfdev/chat/api/routers"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	appEnv := os.Getenv("APP_ENV")

	if "" == appEnv {
		appEnv = "development"
	}

	env := ".env." + appEnv
	err := godotenv.Load(env)
	if err != nil {
		log.Fatal("[env] - Failed to load env")
	} else {
		log.Println("[env] - Env loaded successfully")
	}

	port := os.Getenv("APP_PORT")
	router := routers.NewRoute()

	server := &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// OS interrupt
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Go routine to make server running in background
	// While still listening request, etc.
	go func() {
		log.Printf("[server] - Server is running on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[server] - Failed to run server: %v", err)
		}
	}()

	// While the other or mainly this function does is to detect signal interrupt like ctrl+c

	// Wait for interrupted signal
	<-stop

	log.Println("[server] - Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("[server] - Server forced to shutdown: %v", err)
	}

	log.Println("[server] - Server gracefully stopped")
}
