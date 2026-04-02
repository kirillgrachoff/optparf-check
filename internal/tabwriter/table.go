package tabwriter

import (
	"io"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"go.uber.org/zap"
)

func CreateTable(logger *zap.Logger, out string) (table.Writer, error) {
	var outFd io.Writer
	if out != "" {
		var err error
		outFd, err = os.OpenFile(out, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
		if err != nil {
			logger.Error("unable to create file", zap.Error(err))
		}
	}

	if out == "" || outFd == nil {
		outFd = os.Stdout
	}

	otable := table.NewWriter()
	otable.SetOutputMirror(outFd)

	return otable, nil
}
