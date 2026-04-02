package notify

import (
	"context"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/kirillgrachoff/optparf-check/internal/tabwriter"
	"github.com/kirillgrachoff/optparf-check/internal/types"
	"go.uber.org/zap"
)

type TableNotifier struct {
	logger *zap.Logger
	out    string

	otable table.Writer
}

func NewTableNotifier(logger *zap.Logger, out string) (Notifier, error) {
	n := &TableNotifier{
		logger: logger,
		out:    out,
	}
	err := n.init()
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (t *TableNotifier) init() error {
	otable, err := tabwriter.CreateTable(t.logger, t.out)
	if err != nil {
		t.logger.Error("table not created", zap.Error(err))
		return err
	}

	otable.AppendHeader(table.Row{"peer_id", "pattern", "result"})

	t.otable = otable

	return nil
}

// Flush implements [Notifier].
func (t *TableNotifier) Flush(ctx context.Context, logger *zap.Logger) error {
	t.otable.Render()
	return t.init()
}

// Notify implements [Notifier].
func (t *TableNotifier) Notify(ctx context.Context, logger *zap.Logger, peerId int64, result []types.QueryResult) error {
	for _, r := range result {
		t.otable.AppendRow(table.Row{r.PeerId, r.Pattern, string(r.Found)})
	}
	return nil
}

var _ Notifier = (*TableNotifier)(nil)
