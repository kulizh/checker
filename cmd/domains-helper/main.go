package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"checker/internal/config"
	"checker/internal/domain"
)

const defaultRetries = 3

func usage() {
	fmt.Fprintf(os.Stderr, `Usage:
  domains-helper add [-config path] sitename codes
  domains-helper remove [-config path] sitename

Commands:
  add      Add a site to domains.json
  remove   Remove a site from domains.json

Examples:
  domains-helper add example.com 200,301
  domains-helper remove example.com
`)
}

func loadConfig(path string) (config.Config, error) {
	return config.Load(path)
}

func saveConfig(path string, cfg config.Config) error {
	b, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	return os.WriteFile(path, b, 0o644)
}

func restartService() error {
	scripts := []string{"./stop.sh", "./start.sh"}
	for _, script := range scripts {
		if _, err := os.Stat(script); err != nil {
			return fmt.Errorf("script not found: %s", script)
		}
		cmd := exec.Command("bash", script)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to run %s: %w", script, err)
		}
	}
	return nil
}

func parseCodes(raw string) ([]int, error) {
	parts := strings.Split(raw, ",")
	var codes []int
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		code, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid HTTP status code: %s", part)
		}
		codes = append(codes, code)
	}
	if len(codes) == 0 {
		return nil, fmt.Errorf("at least one status code is required")
	}
	return codes, nil
}

func actionAdd(configPath, domainName, codesRaw string) error {
	cfg, err := loadConfig(configPath)
	if err != nil {
		return err
	}
	if _, ok := cfg[domainName]; ok {
		return fmt.Errorf("domain already exists: %s", domainName)
	}

	codes, err := parseCodes(codesRaw)
	if err != nil {
		return err
	}

	cfg[domainName] = domain.DomainConfig{
		ExpectedCodes: codes,
		Retries:       defaultRetries,
	}

	if err := saveConfig(configPath, cfg); err != nil {
		return err
	}
	fmt.Printf("Added %s with codes: %v\n", domainName, codes)
	return restartService()
}

func actionRemove(configPath, domainName string) error {
	cfg, err := loadConfig(configPath)
	if err != nil {
		return err
	}
	if _, ok := cfg[domainName]; !ok {
		return fmt.Errorf("domain not found: %s", domainName)
	}

	delete(cfg, domainName)
	if err := saveConfig(configPath, cfg); err != nil {
		return err
	}
	fmt.Printf("Removed %s\n", domainName)
	return restartService()
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	subcommand := os.Args[1]
	flagSet := flag.NewFlagSet(subcommand, flag.ExitOnError)
	configPath := flagSet.String("config", "domains.json", "path to domains config")
	flagSet.Parse(os.Args[2:])

	switch subcommand {
	case "add":
		if flagSet.NArg() != 2 {
			usage()
			os.Exit(1)
		}
		domainName := flagSet.Arg(0)
		codesRaw := flagSet.Arg(1)
		err := actionAdd(*configPath, domainName, codesRaw)
		if err != nil {
			fmt.Fprintln(os.Stderr, "ERROR:", err)
			os.Exit(1)
		}
	case "remove":
		if flagSet.NArg() != 1 {
			usage()
			os.Exit(1)
		}
		domainName := flagSet.Arg(0)
		err := actionRemove(*configPath, domainName)
		if err != nil {
			fmt.Fprintln(os.Stderr, "ERROR:", err)
			os.Exit(1)
		}
	default:
		usage()
		os.Exit(1)
	}
}
