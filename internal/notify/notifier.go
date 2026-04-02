package notify

import (
	"context"

	"github.com/kirillgrachoff/optparf-check/internal/types"
	"go.uber.org/zap"
)

type Notifier interface {
	Notify(ctx context.Context, logger *zap.Logger, peerId int64, result []types.QueryResult) error
	Flush(ctx context.Context, logger *zap.Logger) error
}
