package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	api "github.com/softnetics/dotlocal/internal/api/proto"
	"github.com/softnetics/dotlocal/internal/client"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var logger *zap.Logger

var (
	hostname   string
	pathPrefix string
	target     string

	rootCmd = &cobra.Command{
		Use: "dotlocal",
		Run: func(cmd *cobra.Command, args []string) {
			apiClient, err := client.NewApiClient()
			if err != nil {
				log.Fatal(err)
			}

			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()

			stream, err := apiClient.CreateMappingWhileConnected(ctx)
			if err != nil {
				log.Fatal(err)
			}

			logger.Info(fmt.Sprintf("Forwarding %s%s to %s", hostname, pathPrefix, target))
			err = stream.Send(&api.CreateMappingRequest{
				Host:       &hostname,
				PathPrefix: &pathPrefix,
				Target:     &target,
			})
			if err != nil {
				log.Fatal(err)
			}

			<-ctx.Done()
			logger.Info("Shutting down")
			_, err = stream.CloseAndRecv()
			if err != nil {
				log.Fatal(err)
			}
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	var err error
	logger, err = zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}

	rootCmd.Flags().StringVarP(&hostname, "host", "", "", "Hostname to map")
	rootCmd.MarkFlagRequired("host")
	rootCmd.Flags().StringVarP(&pathPrefix, "path-prefix", "p", "", "Path prefix")
	rootCmd.Flags().StringVarP(&target, "target", "t", "", "Target URL")
	rootCmd.MarkFlagRequired("target")
}
