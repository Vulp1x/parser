package instagrapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/inst-api/parser/internal/dbmodel"
	"github.com/inst-api/parser/internal/domain"
	"github.com/inst-api/parser/internal/transport"
	"github.com/inst-api/parser/pkg/logger"
)

// ErrBotIsBlocked аккаунт заблокирован
var ErrBotIsBlocked = errors.New("bot account is blocked")

// ErrBloggerNotFound не нашли блогера по username
var ErrBloggerNotFound = errors.New("blogger not found")

var ErrToManyRequests = errors.New("wait a few minutes before next request")

var ErrBloggerIsPrivate = errors.New("blogger is private couldn't fetch")

type Client struct {
	cli              *http.Client
	saveResponseFunc func(ctx context.Context, sessionID string, response *http.Response, opts ...SaveResponseOption) ([]byte, error)
	host             string
}

func NewClient(host string) *Client {
	return &Client{cli: transport.InitHTTPClient(), saveResponseFunc: saveResponse, host: host}
}

// CheckBot проверяет, что бот
func (c *Client) CheckBot(ctx context.Context, sessionID string) error {
	startedAt := time.Now()
	resp, err := c.getWithCtx(ctx, fmt.Sprintf("%s:/auth/settings/get?sessionid=%s", c.host, sessionID))
	if err != nil {
		return err
	}

	_, err = c.saveResponseFunc(ctx, sessionID, resp, WithElapsedTime(time.Since(startedAt)))
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

	resp, err := c.postFormWithCtx(ctx, c.host+"/user/similar/full", val)
	if err != nil {
		return nil, err
	}

	var respBytes []byte
	respBytes, err = c.saveResponseFunc(ctx, sessionID, resp, WithElapsedTime(time.Since(startedAt)))
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

	var respBytes []byte
	respBytes, err = c.saveResponseFunc(ctx, sessionID, resp, WithElapsedTime(time.Since(startedAt)))
	if err != nil {
		logger.Errorf(ctx, "failed to save response: %v", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusForbidden:
		return nil, ErrBloggerIsPrivate
	case http.StatusTooManyRequests:
		return nil, ErrToManyRequests
	case http.StatusNotFound:
		return nil, ErrBloggerNotFound
	case http.StatusBadRequest:
		return nil, ErrBotIsBlocked

	case http.StatusInternalServerError:
		logger.Errorf(ctx, "got internal error, assuming bot is blocked: %s", string(respBytes))
		return nil, ErrBotIsBlocked
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got %d response code, expected 200", resp.StatusCode)
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

// ParseUsers берет postsToParse последних постов блогера с user_id = bloggerUserID
// Для каждого поста:
// 1. Выбирает commentsToParse пользователей, который комментировали этот пост
// 2. Выбирает likesToParse пользователей, который поставили лайк этому посту
func (c Client) ParseUsers(
	ctx context.Context,
	sessionID string,
	bloggerUserID int64,
	dataset dbmodel.Dataset,
) (domain.ShortInstUsers, error) {
	startedAt := time.Now()
	val := map[string][]string{
		"sessionid":      {sessionID},
		"user_id":        {strconv.FormatInt(bloggerUserID, 10)},
		"posts_count":    {strconv.FormatInt(int64(dataset.PostsPerBlogger), 10)},
		"comments_count": {strconv.FormatInt(int64(dataset.CommentedPerPost), 10)},
		"likes_count":    {strconv.FormatInt(int64(dataset.LikedPerPost), 10)},
	}

	resp, err := c.postFormWithCtx(ctx, c.host+"/user/parse", val)
	if err != nil {
		return nil, err
	}

	var respBytes []byte
	respBytes, err = c.saveResponseFunc(ctx, sessionID, resp, WithElapsedTime(time.Since(startedAt)))
	if err != nil {
		logger.Errorf(ctx, "failed to save response: %v", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusForbidden:
		return nil, ErrBloggerIsPrivate
	case http.StatusTooManyRequests:
		return nil, ErrToManyRequests
	case http.StatusNotFound:
		return nil, ErrBloggerNotFound
	case http.StatusBadRequest:
		return nil, ErrBotIsBlocked

	case http.StatusInternalServerError:
		logger.Errorf(ctx, "got internal error, assuming bot is blocked: %s", string(respBytes))
		return nil, ErrBotIsBlocked
	}

	var users []domain.InstUserShort
	err = json.Unmarshal(respBytes, &users)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal users: %v", err)
	}

	return users, nil
}

func (c Client) postFormWithCtx(ctx context.Context, url string, data url.Values) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.cli.Do(req)
}

func (c *Client) getWithCtx(ctx context.Context, url string) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	return c.cli.Do(req)
}
