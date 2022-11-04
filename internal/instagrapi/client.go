package instagrapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/inst-api/parser/internal/transport"
	"github.com/inst-api/parser/pkg/logger"
)

// ErrBotIsBlocked аккаунт заблокирован
var ErrBotIsBlocked = errors.New("bot account is blocked")

type Client struct {
	cli              *http.Client
	saveResponseFunc func(ctx context.Context, sessionID string, response *http.Response, opts ...SaveResponseOption) error
	host             string
}

func NewClient(host string) *Client {
	return &Client{cli: transport.InitHTTPClient(), saveResponseFunc: saveResponse, host: host}
}

// CheckBot проверяет, что бот
func (c *Client) CheckBot(ctx context.Context, sessionID string) error {
	startedAt := time.Now()
	val := map[string][]string{"sessionid": {sessionID}}

	resp, err := c.cli.PostForm(c.host+"/auth/get_settings", val)
	if err != nil {
		return err
	}

	err = c.saveResponseFunc(ctx, sessionID, resp, WithElapsedTime(time.Since(startedAt)))
	if err != nil {
		logger.Errorf(ctx, "failed to save response: %v", err)
	}

	if resp.StatusCode == http.StatusBadRequest {
		return ErrBotIsBlocked
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got %d response code, expected 200", resp.StatusCode)
	}

	return nil
}
