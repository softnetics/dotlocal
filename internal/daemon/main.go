package daemon

import (
	"context"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"gopkg.in/tomb.v2"
)

func Start(logger *zap.Logger) error {
	dotlocal, err := NewDotLocal(logger)
	if err != nil {
		return err
	}

	apiServer, err := NewAPIServer(logger.Named("api"), dotlocal)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var t tomb.Tomb
	t.Go(func() error {
		return dotlocal.Start()
	})
	t.Go(func() error {
		return apiServer.Start()
	})

	<-ctx.Done()
	logger.Info("Shutting down")

	err = apiServer.Stop()
	if err != nil {
		return err
	}

	return t.Wait()
}
