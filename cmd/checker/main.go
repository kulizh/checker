package main

import (
	"flag"
	"log"
	"time"

	"checker/internal/runner"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	configPath := flag.String("config", "domains.json", "path to domains config")
	interval := flag.Duration("interval", 30*time.Second, "check interval (e.g. 10s, 1m)")
	flag.Parse()

	log.Printf("checker started: config=%s interval=%s\n", *configPath, *interval)

	r, err := runner.New(*configPath, *interval)
	if err != nil {
		log.Fatal(err)
	}

	r.Start()
}
