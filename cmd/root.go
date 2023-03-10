/*
Package cmd contains all the commands and subcommands of this CLI tool
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/controlplaneio/netassertv2-l4-client/pkg/config"
	"github.com/controlplaneio/netassertv2-l4-client/pkg/conntester"
	"github.com/controlplaneio/netassertv2-l4-client/pkg/log"
	"github.com/spf13/cobra"
)

var clientConfig = &config.Config{
	LogLevel:    log.GetDefaultLevel(),
	LogEncoding: log.ConsoleEncoding,
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "client",
	Short:         "A simple TCP / UDP client",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return clientConfig.Init(cmd)
	},
	RunE: start,
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func start(cmd *cobra.Command, args []string) error {
	logger, err := log.DefaultLogger(clientConfig.LogLevel, clientConfig.LogEncoding)
	if err != nil {
		return err
	}
	log.SetLogger(logger)
	logger.Info(fmt.Sprintf("%+v", clientConfig))

	ct, err := conntester.New(clientConfig, logger)
	if err != nil {
		return err
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	c := make(chan os.Signal, 1)
	wait := make(chan error)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	go ct.Start(ctx, wait)

	for {
		select {
		case <-c:
			logger.Info("stop signal received, stopping...")
			cancel()
		case err = <-wait:
			logger.Info("jobs stopped")
			switch {
			case errors.Is(err, context.Canceled):
				logger.Info("exiting because of context cancellation")
				return err
			case errors.Is(err, conntester.ErrTestFailed):
				logger.Info("test failed")
				return err
			case errors.Is(err, nil):
				logger.Info("test passed")
				return nil
			default:
				logger.Info("unknown error")
				return err
			}
		}
	}
}

func init() {
	rootCmd.PersistentFlags().VarP(&clientConfig.LogLevel, "log-level", "l", "set log level")
	rootCmd.PersistentFlags().VarP(&clientConfig.LogEncoding, "log-encoding", "e", "set log encoding")

	rootCmd.Flags().StringVarP(&clientConfig.Protocol, "protocol", "P", "tcp", "either tcp or udp")
	rootCmd.Flags().StringVar(&clientConfig.TargetHost, "target-host", "", "target host (either its name or IP address)")
	rootCmd.Flags().Uint16VarP(&clientConfig.TargetPort, "target-port", "p", 0, "target port")
	rootCmd.Flags().StringVarP(&clientConfig.Message, "message", "m", "defaultmessage", "message to send")
	rootCmd.Flags().UintVarP(&clientConfig.Timeout, "timeout", "t", 2000, "timeout in ms")
	rootCmd.Flags().UintVarP(&clientConfig.Attempts, "attempts", "r", 1, "number of attempts, successful or not")
	rootCmd.Flags().UintVar(&clientConfig.Period, "period", 5000, "send a new message every <period> ms")
	rootCmd.Flags().UintVar(&clientConfig.SuccThrPec, "success-threshold", 80, "percentage of successful attempts needed")
}
