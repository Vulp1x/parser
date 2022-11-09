package instagrapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/inst-api/parser/internal/domain"
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

func (c Client) FindSimilarBloggers(ctx context.Context, sessionID, bloggerUserName string) (domain.InstUsers, error) {
	startedAt := time.Now()
	val := map[string][]string{"sessionid": {sessionID}, "user_name": {bloggerUserName}}

	resp, err := c.cli.PostForm(c.host+"/auth/get_settings", val)
	if err != nil {
		return nil, err
	}

	err = c.saveResponseFunc(ctx, sessionID, resp, WithElapsedTime(time.Since(startedAt)), WithReuseResponseBody(true))
	if err != nil {
		logger.Errorf(ctx, "failed to save response: %v", err)
	}

	if resp.StatusCode == http.StatusBadRequest {
		return nil, ErrBotIsBlocked
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got %d response code, expected 200", resp.StatusCode)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %v", err)
	}

	var users []domain.InstUser
	err = json.Unmarshal(respBytes, &users)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal users: %v", err)
	}

	return users, nil
}
