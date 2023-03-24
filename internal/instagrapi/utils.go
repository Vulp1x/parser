package instagrapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/inst-api/parser/pkg/logger"
	"go.uber.org/zap"
)

var once = &sync.Once{}
var accessLogFile *os.File

func init() {
	once.Do(func() {
		var err error
		accessLogFile, err = os.OpenFile("tmp/instagrapi-rest.access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			logger.Fatalf(context.Background(), "failed to open instagrapi access file log: %v, err")
		}
	})
}

func saveResponse(ctx context.Context, sessionID string, resp *http.Response, opts ...SaveResponseOption) ([]byte, error) {
	var cfg saveResponseConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	startedAt := time.Now()
	if resp == nil {
		return []byte{}, fmt.Errorf("empty resp")
	}

	headerBytes, err := json.Marshal(resp.Header)
	if err != nil {
		logger.Errorf(ctx, "failed to marshal resp headers: %v", err)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf(ctx, "failed to read response body: %v", err)
	}

	err = resp.Body.Close()
	if err != nil {
		logger.Errorf(ctx, "failed to close response body: %v", err)
	}

	var fields []zap.Field

	if cfg.elapsed != nil {
		fields = append(fields, zap.String("elapsed_time", cfg.elapsed.String()))
	}

	fields = append(fields,
		zap.Int("response_code", resp.StatusCode),
		zap.Int64("response_len", resp.ContentLength),
		zap.String("session_id", sessionID),
		zap.ByteString("headers", headerBytes),
		zap.ByteString("body", bodyBytes),
	)

	log := logger.NewWithSink(
		zap.DebugLevel,
		accessLogFile,
		zap.WithCaller(true),
		zap.AddCallerSkip(1),
		zap.Fields(fields...),
	).Sugar()

	log.Infof("saving response from instagrapi, saving took %s", time.Since(startedAt))
	return bodyBytes, nil
}

type saveResponseConfig struct {
	elapsed *time.Duration
}

type SaveResponseOption func(config *saveResponseConfig)

func WithElapsedTime(duration time.Duration) SaveResponseOption {
	return func(config *saveResponseConfig) {
		config.elapsed = &duration
	}
}
