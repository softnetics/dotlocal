package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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

			go func() {
				wasSuccessful := false
				for {
					_, err := apiClient.CreateMapping(ctx, &api.CreateMappingRequest{
						Host:       &hostname,
						PathPrefix: &pathPrefix,
						Target:     &target,
					})
					duration := 1 * time.Minute
					if err != nil {
						logger.Error("Failed to add mapping. Retrying in 5 seconds.", zap.Error(err))
						duration = 5 * time.Second
						wasSuccessful = false
					} else if !wasSuccessful {
						logger.Info(fmt.Sprintf("Forwarding %s%s to %s", hostname, pathPrefix, target))
						wasSuccessful = true
					}
					timer := time.NewTimer(duration)
					select {
					case <-timer.C:
						continue
					case <-ctx.Done():
						return
					}
				}
			}()

			<-ctx.Done()
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
