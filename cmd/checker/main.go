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
	renotify := flag.Duration("renotify", 1*time.Hour, "re-notify interval for DOWN sites (e.g. 30m, 1h)")
	flag.Parse()

	log.Printf("checker started: config=%s interval=%s renotify=%s\n", *configPath, *interval, *renotify)

	r, err := runner.New(*configPath, *interval, *renotify)
	if err != nil {
		log.Fatal(err)
	}

	r.Start()
}
