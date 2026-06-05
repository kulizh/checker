package runner

import (
	"fmt"
	"time"

	"checker/internal/checker"
	"checker/internal/config"
	"checker/internal/notify"
	"checker/internal/probe"
	"checker/internal/state"
)

type Runner struct {
	Worker   *checker.Worker
	Interval time.Duration
}

func New(cfgPath string, interval, renotifyAfter time.Duration) (*Runner, error) {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return nil, err
	}

	st := state.New()
	tg := notify.NewTelegram()

	p := probe.NewHTTPProbe(5 * time.Second)

	ch := &checker.Checker{
		Probe: p,
		State: st,
	}

	worker := &checker.Worker{
		Config:        cfg,
		Checker:       ch,
		State:         st,
		Notify:        tg,
		RenotifyAfter: renotifyAfter,
	}

	return &Runner{
		Worker:   worker,
		Interval: interval,
	}, nil
}

func (r *Runner) Start() {
	for {
		start := time.Now()

		r.Worker.Run()

		elapsed := time.Since(start)
		if elapsed > r.Interval {
			fmt.Printf("[runner] cycle took %s (longer than interval %s)\n", elapsed, r.Interval)
		} else {
			fmt.Printf("[runner] cycle finished in %s\n", elapsed)
		}

		time.Sleep(r.Interval)
	}
}
