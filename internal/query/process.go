package query

import (
	"context"

	"github.com/kirillgrachoff/optparf-check/internal/types"
	"go.uber.org/zap"
)

func Process(ctx context.Context, logger *zap.Logger, client *QeuryClient, parser *Parser, queries []types.Query) QueryResult {
	result := QueryResult{}

	for _, q := range queries {
		logger := logger.With(
			zap.String("pattern", q.Pattern),
			zap.Int64("peer_id", q.PeerId),
		)

		logger.Info("processing")

		r := types.QueryResult{PeerId: q.PeerId, Pattern: q.Pattern}
		b, err := client.Get(ctx, logger, q.Pattern)
		if err != nil {
			logger.Error("error while querying", zap.Error(err))
			r.Err = err
		} else {
			r.Items = parser.Parse(b)
		}

		result[q.PeerId] = append(result[q.PeerId], r)
	}

	return result
}

type QueryResult = map[int64][]types.QueryResult
