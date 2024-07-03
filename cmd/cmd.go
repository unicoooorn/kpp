package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/unicoooorn/docker-monitoring-tool/internal/checker"
	"github.com/unicoooorn/docker-monitoring-tool/internal/config"
	"github.com/unicoooorn/docker-monitoring-tool/internal/docker"
	"github.com/unicoooorn/docker-monitoring-tool/internal/tool"
)

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:  "dmt",
		Long: "A tool for monitoring and management of docker containers",
		RunE: run,
	}

	rootCmd.PersistentFlags().StringP("config", "c", "config/config.yaml", "specify a config file")

	return rootCmd
}

func run(cmd *cobra.Command, _ []string) error {
	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()

	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		return fmt.Errorf("failed to get <config> flag: %w", err)
	}

	cfg, err := config.LoadApp(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client, err := docker.NewClient()
	if err != nil {
		return fmt.Errorf("new docker client: %w", err)
	}
	defer client.Close()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	checkerList := make([]tool.Checker, 0)

	checkerList = append(checkerList, checker.NewDiskUsageChecker(cfg.DiskUsage))

	if err := tool.New(client, checkerList, *cfg, *logger).Run(ctx); errors.Is(err, context.Canceled) {
		logger.Info("shutdown")
		return nil
	} else if err != nil {
		return fmt.Errorf("running tool: %w", err)
	}

	return nil
}

func Execute() {
	err := newRootCmd().Execute()

	if err != nil {
		os.Exit(1)
	}
}
