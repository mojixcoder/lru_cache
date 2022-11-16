package main

import (
	"log"

	"github.com/MojtabaArezoomand/lru_cache/internal/server"
)

func main() {
	log.Println("Running server...")
	server.RunServer()
}
