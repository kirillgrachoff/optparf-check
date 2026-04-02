package query

import (
	"context"

	"github.com/kirillgrachoff/optparf-check/internal/types"
	"go.uber.org/zap"
)

func Process(ctx context.Context, logger *zap.Logger, client *QeuryClient, filter *Filter, queries []types.Query) QueryResult {
	result := QueryResult{}

	for _, q := range queries {
		logger := logger.With(
			zap.String("pattern", q.Pattern),
			zap.Int64("peer_id", q.PeerId),
		)

		logger.Info("processing")

		b, err := client.Get(ctx, logger, q.Pattern)
		if err != nil {
			logger.Error("error while querying", zap.Error(err))
		}
		filtered := filter.Filter(b)

		found := result[q.PeerId]
		found = append(found, types.QueryResult{PeerId: q.PeerId, Pattern: q.Pattern, Found: filtered})
		result[q.PeerId] = found
	}

	return result
}

type QueryResult = map[int64][]types.QueryResult
