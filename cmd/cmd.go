package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/unicoooorn/kpp/internal/checker"
	"github.com/unicoooorn/kpp/internal/config"
	"github.com/unicoooorn/kpp/internal/docker"
	"github.com/unicoooorn/kpp/internal/tool"
)

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:  "kpp",
		Long: "A lightweight tool for imposing limits on docker volumes usage",
		RunE: run,
	}

	rootCmd.PersistentFlags().StringP("config", "c", "config/config.yaml", "specify a config file")
	rootCmd.Flags().Bool("dry-run", false, "run in a dry-run")
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

	dryRun, _ := cmd.Flags().GetBool("dry-run")

	var client interface {
		tool.ContainerManager
		io.Closer
	}
	if dryRun {
		fmt.Println("[dry-run] running a dry-run mode")
		client, err = docker.NewDryRunClient()
	} else {
		client, err = docker.NewClient()
	}
	if err != nil {
		return fmt.Errorf("new docker client: %w", err)
	}
	defer client.Close()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	checkerList := make([]tool.Checker, 0)

	checkerList = append(checkerList, checker.NewDiskUsageChecker(cfg.DiskUsage))

	stats, err := client.ContainersStats(ctx)
	if err != nil {
		return fmt.Errorf("failed to receive container stats: %w", err)
	}
	var fmchecker tool.Checker
	switch cfg.FileMonitoring.Type {
	case "whitelist":
		fmchecker, err = checker.NewWhiteListFileMonitoringChecker(cfg.FileMonitoring, stats)
	case "blacklist":
		fmchecker, err = checker.NewBlackListFileMonitoringChecker(cfg.FileMonitoring, stats)
	}
	if err != nil {
		return fmt.Errorf("failed to init file monitoring checker %w", err)
	}
	checkerList = append(checkerList, fmchecker)

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
