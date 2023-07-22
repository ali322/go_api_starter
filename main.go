package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"app/handler"
	"app/lib/config"
	"app/lib/logger"
	"app/util"

	"github.com/gorilla/mux"
)

func setup() *mux.Router {
	util.EnsurePath(config.App.LogDir)
	logger := logger.Init(config.App.LogDir)
	defer logger.Sync()
	app := handler.ApplyRoutes()
	return app
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	config.Read()
	app := setup()
	server := &http.Server{
		Handler: app,
		Addr:    fmt.Sprintf(":%d", config.App.Port),
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to listen: %s\n", err)
		}
		log.Printf("server started at %d\n", config.App.Port)
	}()
	<-ctx.Done()
	stop()
	log.Println("shutdown gracefully, press ctrl+c force shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("failed to shutdown server: %s\n", err)
	}
	log.Println("server exiting")
}
