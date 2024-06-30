package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/unicoooorn/docker-monitoring-tool/internal/config"
)

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:  "dmt",
		Long: "A tool for monitoring and management of docker containers",
		RunE: run,
	}

	rootCmd.PersistentFlags().StringP("config", "c", "configs/config.yaml", "specify a config file")

	return rootCmd
}

func run(cmd *cobra.Command, _ []string) error {
	_, cancel := context.WithCancel(cmd.Context())
	defer cancel()

	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		return fmt.Errorf("failed to get <config> flag: %w", err)
	}

	cfg, err := config.LoadApp(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Println(cfg.DiskLimit) // fixme: stub
	return nil
}

func Execute() {
	err := newRootCmd().Execute()

	if err != nil {
		os.Exit(1)
	}
}
