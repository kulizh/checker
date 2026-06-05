package checker

import (
	"fmt"
	"sync"
	"time"

	"checker/internal/config"
	"checker/internal/domain"
	"checker/internal/notify"
	"checker/internal/state"
)

type Worker struct {
	Config        config.Config
	Checker       *Checker
	State         *state.Store
	Notify        *notify.Telegram
	RenotifyAfter time.Duration

	initialized bool
}

func (w *Worker) Run() {
	var wg sync.WaitGroup

	if !w.initialized {
		w.baseline()
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

			if exists {
				newState.LastNotifiedAt = prev.LastNotifiedAt
			}

			now := time.Now()

			switch {
			case !exists:
				// First check after baseline — notify immediately if DOWN
				if !res.Up {
					msg := fmt.Sprintf("[DOWN] %s — status code %d (%s)", d, res.Code, res.Error)
					w.tryNotify(d, msg, now, &newState)
				}

			case prev.Status != newState.Status:
				// Status changed — notify immediately
				label := "DOWN"
				if res.Up {
					label = "UP"
				}
				msg := fmt.Sprintf("[%s] %s — status code %d", label, d, res.Code)
				if res.Error != "" {
					msg += fmt.Sprintf(" (%s)", res.Error)
				}
				w.tryNotify(d, msg, now, &newState)

			case !res.Up && now.Sub(prev.LastNotifiedAt) >= w.RenotifyAfter:
				// Still down and re-notify interval passed
				msg := fmt.Sprintf("[DOWN REMINDER] %s — still down after %s (code %d)", d, w.RenotifyAfter, res.Code)
				if res.Error != "" {
					msg += fmt.Sprintf(" (%s)", res.Error)
				}
				w.tryNotify(d, msg, now, &newState)
			}

			w.State.Set(d, newState)
		}(domainName, cfg)
	}

	wg.Wait()
}

func (w *Worker) baseline() {
	fmt.Println("[baseline] initial check of all sites")

	var (
		wg    sync.WaitGroup
		mu    sync.Mutex
		lines []string
	)

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

			statusLine := fmt.Sprintf("  %s — %s", d, newState.Status)
			if !res.Up {
				statusLine += fmt.Sprintf(" (code %d)", res.Code)
				if res.Error != "" {
					statusLine += fmt.Sprintf(" [%s]", res.Error)
				}
			}
			fmt.Println("[baseline]", statusLine)

			mu.Lock()
			lines = append(lines, statusLine)
			mu.Unlock()
		}(domainName, cfg)
	}

	wg.Wait()

	// Send one summary notification
	msg := "Sites added to monitoring:\n"
	for _, l := range lines {
		msg += l + "\n"
	}
	if err := w.Notify.Send(msg); err != nil {
		fmt.Printf("[notify] failed to send baseline: %v\n", err)
	}

	w.initialized = true
}

func (w *Worker) tryNotify(domainName, msg string, now time.Time, state *domain.DomainState) {
	if err := w.Notify.Send(msg); err != nil {
		fmt.Printf("[notify] failed for %s: %v\n", domainName, err)
		return
	}
	state.LastNotifiedAt = now
}
