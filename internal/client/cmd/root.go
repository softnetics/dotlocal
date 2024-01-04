package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"time"

	api "github.com/softnetics/dotlocal/internal/api/proto"
	"github.com/softnetics/dotlocal/internal/client"
	"github.com/softnetics/dotlocal/internal/util"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var logger *zap.Logger

var (
	hostname     string
	pathPrefix   string
	targetArg    string
	overridePort string

	rootCmd = &cobra.Command{
		Use: "dotlocal",
		Run: func(cmd *cobra.Command, args []string) {
			apiClient, err := client.NewApiClient()
			if err != nil {
				log.Fatal(err)
			}
			target := getTarget()

			loopCtx, cancel := context.WithCancel(context.Background())
			defer cancel()
			go func() {
				wasSuccessful := false
				for {
					_, err := apiClient.CreateMapping(loopCtx, &api.CreateMappingRequest{
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
					case <-loopCtx.Done():
						return
					}
				}
			}()

			if len(args) > 0 {
				ch := make(chan os.Signal, 1)
				signal.Notify(ch, os.Interrupt, os.Kill)

				cmd := exec.Command(args[0], args[1:]...)
				cmd.Env = os.Environ()
				if overridePort != "" {
					cmd.Env = append(cmd.Env, fmt.Sprintf("PORT=%s", overridePort))
				}
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err = cmd.Start()
				if err != nil {
					log.Fatal(err)
				}

				go func() {
					for {
						select {
						case <-ch:
							cancel()
							cmd.Process.Kill()
							return
						case <-loopCtx.Done():
							return
						}
					}
				}()

				err = cmd.Wait()
				if err != nil {
					if exiterr, ok := err.(*exec.ExitError); ok {
						os.Exit(exiterr.ExitCode())
					} else {
						log.Fatal(err)
					}
				}
			} else {
				<-context.Background().Done()
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
	rootCmd.Flags().StringVarP(&targetArg, "target", "t", "", "Target URL")
}

func getTarget() string {
	if targetArg != "" {
		return targetArg
	}
	portString := os.Getenv("PORT")
	port, err := strconv.Atoi(portString)
	if err != nil {
		port, err = util.FindAvailablePort()
		if err != nil {
			log.Fatal(err)
		}
	}
	overridePort = strconv.Itoa(port)
	return fmt.Sprintf("http://localhost:%d", port)
}
