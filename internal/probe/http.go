package probe

import (
	"net/http"
	"time"
)

type HTTPProbe struct {
	Client *http.Client
}

func NewHTTPProbe(timeout time.Duration) *HTTPProbe {
	return &HTTPProbe{
		Client: &http.Client{Timeout: timeout},
	}
}

func (p *HTTPProbe) Check(url string) (int, error) {
	resp, err := p.Client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}
