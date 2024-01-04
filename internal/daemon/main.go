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
	err = dotlocal.Start()
	if err != nil {
		return err
	}

	// err = dotlocal.SetMappings([]internal.Mapping{
	// 	{
	// 		Host:       "app1.local",
	// 		PathPrefix: "/",
	// 		Target:     "http://localhost:3000",
	// 	},
	// 	{
	// 		Host:       "app1.local",
	// 		PathPrefix: "/_api",
	// 		Target:     "http://localhost:4000",
	// 	},
	// 	{
	// 		Host:       "app2.local",
	// 		PathPrefix: "/",
	// 		Target:     "http://localhost:5555",
	// 	},
	// })
	// if err != nil {
	// 	panic(err)
	// }

	apiServer, err := NewAPIServer(logger.Named("api"), dotlocal)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var t tomb.Tomb
	t.Go(func() error {
		return apiServer.Start()
	})

	<-ctx.Done()
	logger.Info("Shutting down")

	err = apiServer.Stop()
	if err != nil {
		return err
	}
	err = dotlocal.Stop()
	if err != nil {
		return err
	}

	return t.Wait()
}
