package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"golang.org/x/net/proxy"
)

type Telegram struct {
	Token  string
	ChatID string
	Client *http.Client
}

func NewTelegram() *Telegram {
	proxyAddr := os.Getenv("PROXY")

	transport := &http.Transport{}

	if proxyAddr != "" {
		dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
		if err == nil {
			transport.DialContext = dialer.(proxy.ContextDialer).DialContext
		}
	}

	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: transport,
	}

	return &Telegram{
		Token:  os.Getenv("TELEGRAM_BOT_TOKEN"),
		ChatID: os.Getenv("TELEGRAM_CHAT_ID"),
		Client: client,
	}
}

type tgResponse struct {
	Ok          bool   `json:"ok"`
	Description string `json:"description"`
}

func (t *Telegram) Send(msg string) error {
	api := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.Token)

	form := url.Values{}
	form.Add("chat_id", t.ChatID)
	form.Add("text", msg)

	req, err := http.NewRequest("POST", api, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return fmt.Errorf("create request error: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := t.Client.Do(req)
	if err != nil {
		if os.IsTimeout(err) {
			return fmt.Errorf("telegram request timeout: %w", err)
		}
		return fmt.Errorf("telegram request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"telegram http error: status=%d body=%s",
			resp.StatusCode,
			string(body),
		)
	}

	var result tgResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("telegram response decode error: %w body=%s", err, string(body))
	}

	if !result.Ok {
		return fmt.Errorf("telegram api error: %s", result.Description)
	}

	return nil
}
