package checker

import (
	"checker/internal/domain"
	"checker/internal/probe"
	"checker/internal/state"
)

type Checker struct {
	Probe *probe.HTTPProbe
	State *state.Store
}

func isAllowed(code int, allowed []int) bool {
	for _, c := range allowed {
		if c == code {
			return true
		}
	}
	return false
}

func (c *Checker) Check(domainName string, cfg domain.DomainConfig) domain.CheckResult {
	url := "https://" + domainName

	var code int
	var err error

	for i := 0; i <= cfg.Retries; i++ {
		code, err = c.Probe.Check(url)
		if err == nil {
			break
		}
	}

	if err != nil {
		return domain.CheckResult{
			Domain: domainName,
			Up:     false,
			Error:  err.Error(),
		}
	}

	return domain.CheckResult{
		Domain: domainName,
		Up:     isAllowed(code, cfg.ExpectedCodes),
		Code:   code,
	}
}
