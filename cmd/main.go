package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/kirillgrachoff/optparf-check/internal/bootstrap"
	"github.com/kirillgrachoff/optparf-check/internal/config"
	"github.com/kirillgrachoff/optparf-check/internal/notify"
	"github.com/kirillgrachoff/optparf-check/internal/query"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

var (
	configPath     string
	queryPath      string
	outFile        string
	concurrencyMax int
	serve          bool

	notifyMode string
)

func init() {
	pflag.StringVar(&configPath, "config", "./config.yml", "path to config")
	pflag.StringVar(&queryPath, "query", "./config-query.yml", "path to query config")
	pflag.IntVar(&concurrencyMax, "concurrency", 1, "max concurrent requests")
	pflag.StringVarP(&outFile, "out", "o", "", "out file")
	pflag.BoolVar(&serve, "serve", false, "start a server which will be querying every period")

	pflag.StringVar(&notifyMode, "mode", "cli", "one of: [cli, tg]")
}

func main() {
	logger := bootstrap.InitLogger()
	defer logger.Sync()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGKILL, syscall.SIGINT)
	defer cancel()

	pflag.Parse()
	logger.Debug("parse flags", zap.String("config_path", configPath))

	var config config.Config
	bootstrap.LoadConfig(logger, configPath, queryPath, &config)
	if err := config.Validate(); err != nil {
		logger.Fatal("invalid config", zap.Error(err))
	}

	var notifier notify.Notifier
	var err error
	switch notifyMode {
	case "tg":
		notifier, err = notify.NewTgNotifier(config.Telegram.Secret)
	case "cli":
		notifier, err = notify.NewTableNotifier(logger, outFile)
	default:
		logger.Fatal("invalid parameter: use one of [cli, tg]", zap.String("mode", notifyMode))
	}
	if err != nil {
		logger.Error("error while creating notifier", zap.Error(err))
	}

	filter := query.NewFilter(".*svc-price.*", "")
	queryClient := query.NewQueryClient(config.Http.QueryPrefix, config.Http.QuerySuffix)

	var ticker *time.Ticker
	now := time.Now()

	if config.Period > 0 {
		ticker = time.NewTicker(config.Period)
	}

	for {
		logger := logger.With(zap.Time("start", now))
		logger.Info("start processing")
		result := query.Process(ctx, logger, queryClient, filter, config.Queries)

		for peerId, r := range result {
			notifier.Notify(ctx, logger, peerId, r)
		}
		notifier.Flush(ctx, logger)

		if !serve {
			break
		}

		select {
		case t := <-ticker.C:
			now = t
		case <-ctx.Done():
			logger.Error("context cancelled", zap.Error(ctx.Err()))
			return
		}
	}
}
