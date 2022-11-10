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

// ErrBloggerNotFound не нашли блогера по username
var ErrBloggerNotFound = errors.New("blogger not found")

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
	resp, err := c.cli.Get(fmt.Sprintf("%s:/auth/settings/get?sessionid=%s", c.host, sessionID))
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

// FindSimilarBloggers находит похожих блогеров, первый аккаунт - блогер, для которого искали похожих
func (c Client) FindSimilarBloggers(ctx context.Context, sessionID, bloggerUserName string) (domain.InstUsers, error) {
	startedAt := time.Now()
	val := map[string][]string{"sessionid": {sessionID}, "username": {bloggerUserName}}

	resp, err := c.cli.PostForm(c.host+"/user/similar/full", val)
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

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrBloggerNotFound
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

	if len(users) == 0 {
		return nil, ErrBloggerNotFound
	}

	return users, nil
}

// FindSimilarBloggersShort находит похожих блогеров, первый аккаунт - блогер, для которого искали похожих
func (c Client) FindSimilarBloggersShort(ctx context.Context, sessionID, bloggerUserName string) (domain.ShortInstUsers, error) {
	startedAt := time.Now()
	val := map[string][]string{"sessionid": {sessionID}, "username": {bloggerUserName}}

	resp, err := c.cli.PostForm(c.host+"/user/similar", val)
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

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrBloggerNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got %d response code, expected 200", resp.StatusCode)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %v", err)
	}

	var users []domain.InstUserShort
	err = json.Unmarshal(respBytes, &users)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal users: %v", err)
	}

	if len(users) == 0 {
		return nil, ErrBloggerNotFound
	}

	return users, nil
}
