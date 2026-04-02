package query

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type QeuryClient struct {
	Prefix string
	Suffix string

	HttpClient *http.Client
}

func NewQueryClient(prefix, suffix string) *QeuryClient {
	return &QeuryClient{
		Prefix:     prefix,
		Suffix:     suffix,
		HttpClient: http.DefaultClient,
	}
}

func (q QeuryClient) Get(ctx context.Context, logger *zap.Logger, payload string) ([]byte, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s%s%s", q.Prefix, payload, q.Suffix),
		nil,
	)
	if err != nil {
		return nil, err
	}
	rsp, err := q.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	logger.Debug("query processed", zap.Int("response_len", len(b)))
	return b, nil
}
