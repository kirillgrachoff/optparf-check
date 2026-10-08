package notify

import (
	"context"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
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

	otable.AppendHeader(table.Row{"peer_id", "pattern", "article", "name", "price"})
	otable.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true},
		{Number: 2, AutoMerge: true},
		{Number: 5, Align: text.AlignRight},
	})
	otable.SetStyle(table.StyleRounded)
	otable.Style().Options.SeparateRows = false

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
		pattern := r.Pattern
		switch {
		case r.Err != nil:
			t.otable.AppendRow(table.Row{r.PeerId, pattern, "", "error: " + r.Err.Error(), ""})
		case len(r.Items) == 0:
			t.otable.AppendRow(table.Row{r.PeerId, pattern, "", "nothing found", ""})
		}
		for _, it := range r.Items {
			t.otable.AppendRow(table.Row{r.PeerId, pattern, it.Article, it.Name, strings.ReplaceAll(formatPrice(it.Price), nbsp, " ")})
		}
		t.otable.AppendSeparator()
	}
	return nil
}

var _ Notifier = (*TableNotifier)(nil)
