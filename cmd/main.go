package main

import (
	"fmt"
	"log"

	"github.com/Ramiro-9/api-gateway/internal/config"
	"github.com/Ramiro-9/api-gateway/internal/logger"
	"github.com/Ramiro-9/api-gateway/internal/router"
)

func main() {
	config.Load()
	logger.Init()

	r := router.Setup()

	addr := fmt.Sprintf(":%s", config.Cfg.GatewayPort)
	log.Printf("Gateway corriendo en http://localhost%s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatal("Error arrancando el gateway:", err)
	}
}
