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

func New(cfgPath string, interval time.Duration) (*Runner, error) {
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
		Config:  cfg,
		Checker: ch,
		State:   st,
		Notify:  tg,
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

		fmt.Printf("cycle finished in %s\n", time.Since(start))

		time.Sleep(r.Interval)
	}
}
