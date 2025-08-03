package main

import (
    "log"
    "github.com/Abdelrahiim/lms/internal/config"
    "github.com/Abdelrahiim/lms/internal/server"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatal("Failed to load config:", err)
    }

    // Create server
    srv, err := server.New(cfg)
    if err != nil {
        log.Fatal("Failed to create server:", err)
    }

    // Start server
    if err := srv.Start(); err != nil {
        log.Fatal("Server failed:", err)
    }
}