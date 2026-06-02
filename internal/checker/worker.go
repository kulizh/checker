package checker

import (
	"fmt"
	"sync"

	"checker/internal/config"
	"checker/internal/domain"
	"checker/internal/notify"
	"checker/internal/state"
)

type Worker struct {
	Config  config.Config
	Checker *Checker
	State   *state.Store
	Notify  *notify.Telegram

	initialized bool
}

func (w *Worker) Run() {
	var wg sync.WaitGroup

	if !w.initialized {
		fmt.Println("[cold start]")

		w.State.Clear()

		wg.Add(len(w.Config))

		for domainName, cfg := range w.Config {
			go func(d string, c domain.DomainConfig) {
				defer wg.Done()

				res := w.Checker.Check(d, c)

				newState := domain.DomainState{
					Code:  res.Code,
					Error: res.Error,
				}

				if res.Up {
					newState.Status = "UP"
				} else {
					newState.Status = "DOWN"
				}

				w.State.Set(d, newState)

				fmt.Printf("[BASELINE] %s → %s (%d)\n", d, newState.Status, newState.Code)
			}(domainName, cfg)
		}

		wg.Wait()

		w.initialized = true
		return
	}

	wg.Add(len(w.Config))

	for domainName, cfg := range w.Config {
		go func(d string, c domain.DomainConfig) {
			defer wg.Done()

			res := w.Checker.Check(d, c)

			prev, exists := w.State.Get(d)

			newState := domain.DomainState{
				Code:  res.Code,
				Error: res.Error,
			}

			if res.Up {
				newState.Status = "UP"
			} else {
				newState.Status = "DOWN"
			}

			shouldAlert := false

			if exists && prev.Status != newState.Status {
				shouldAlert = true
			}

			if shouldAlert {
				symb := "🚨"
				if newState.Status == "UP" {
					symb = "✅"
				}

				msg := fmt.Sprintf(
					"%s %s is %s (code: %d)",
					symb,
					d,
					newState.Status,
					res.Code,
				)

				_ = w.Notify.Send(msg)
			}

			w.State.Set(d, newState)
		}(domainName, cfg)
	}

	wg.Wait()
}
