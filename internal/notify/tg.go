package notify

import (
	"context"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/kirillgrachoff/optparf-check/internal/types"
	"go.uber.org/zap"
)

type TgNotifier struct {
	bot *bot.Bot

	data map[int64][]types.QueryResult
}

func NewTgNotifier(token string) (Notifier, error) {
	n := &TgNotifier{
		data: make(map[int64][]types.QueryResult),
	}
	err := n.init(token)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (t *TgNotifier) init(token string) error {
	tgbot, err := bot.New(
		token,
		bot.WithMiddlewares(
		),
	)
	if err != nil {
		return err
	}
	t.bot = tgbot
	return nil
}

func (t *TgNotifier) sendMessage(ctx context.Context, logger *zap.Logger, peerId int64, text string) error {
	msg, err := t.bot.SendMessage(
		ctx,
		&bot.SendMessageParams{
			Text:   text,
			ChatID: peerId,
		},
	)
	logger.Info("message sent", zap.Any("message", msg))
	if err != nil {
		logger.Error("msg not sent", zap.Error(err))
		return err
	}
	return nil
}

// Flush implements [Notifier].
func (t *TgNotifier) Flush(ctx context.Context, logger *zap.Logger) error {
	preparedMessage := map[int64]*strings.Builder{}

	for peerId, res := range t.data {
		b, ok := preparedMessage[peerId]
		if !ok {
			b = new(strings.Builder)
			preparedMessage[peerId] = b
		}
		for _, r := range res {
			b.WriteString("Pattern: ")
			b.WriteString(r.Pattern)
			b.WriteRune('\n')
			b.Write(r.Found)
			b.WriteRune('\n')
		}
	}

	for peerId, b := range preparedMessage {
		err := t.sendMessage(ctx, logger, peerId, b.String())
		if err != nil {
			return err
		}
	}

	clear(t.data)

	return nil
}

// Notify implements [Notifier].
func (t *TgNotifier) Notify(ctx context.Context, logger *zap.Logger, peerId int64, result []types.QueryResult) error {
	s := t.data[peerId]
	s = append(s, result...)
	t.data[peerId] = s
	return nil
}

var _ Notifier = (*TgNotifier)(nil)
