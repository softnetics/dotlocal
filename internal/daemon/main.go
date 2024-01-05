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
	err = t.Wait()
	if err != nil {
		return err
	}

	var t2 tomb.Tomb
	t2.Go(func() error {
		return apiServer.Serve()
	})

	<-ctx.Done()
	logger.Info("Shutting down")

	t2.Go(func() error {
		return apiServer.Stop()
	})
	t2.Go(func() error {
		return dotlocal.Stop()
	})

	return t.Wait()
}
