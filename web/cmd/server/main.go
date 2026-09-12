package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"agigame/web/server"
)

func main() {
	configPath := flag.String("config", "web/config.yaml", "path to config.yaml")
	flag.Parse()

	cfg, err := server.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	if _, err := os.Stat(cfg.ROM); err != nil {
		log.Fatalf("ROM not found at %q: %v\nProvide your own legally owned GameBoy (GameBoy / Super Mario Land) ROM dump.", cfg.ROM, err)
	}

	srv, err := server.New(cfg)
	if err != nil {
		log.Fatalf("init server: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	addr := cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port)
	httpServer := &http.Server{Addr: addr, Handler: srv.Handler()}

	go func() {
		<-sig
		log.Println("shutting down...")
		cancel()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer shutdownCancel()
		_ = httpServer.Shutdown(shutdownCtx)
	}()

	go func() {
		if err := srv.Start(ctx); err != nil {
			log.Printf("emulator error: %v", err)
			cancel()
		}
	}()

	log.Printf("GoBoy-LLM listening on http://%s (game=%s)", addr, cfg.Game)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
	log.Println("server stopped")
}