package notify

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

var token string

func init() {
	token = os.Getenv("TG_SECRET")
}

func TestExampleSend(t *testing.T) {
	logger := zaptest.NewLogger(t)

	logger.Debug("check token", zap.String("token", token))

	tg := &TgNotifier{}
	err := tg.init(os.Getenv("TG_SECRET"))
	require.NoError(t, err)
	err = tg.sendMessage(t.Context(), logger, int64(0), "Hello from local machine")
	assert.NoError(t, err, "message not sent")
}
