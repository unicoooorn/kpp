package tool

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/unicoooorn/kpp/internal/config"
)

type Tool struct {
	cfg              config.Config
	logger           slog.Logger
	containerManager ContainerManager
	checkers         []Checker
	action           func(ContainerManager, context.Context, string) error
}

func New(containerManager ContainerManager, checkers []Checker, cfg config.Config, logger slog.Logger) *Tool {

	var action func(ContainerManager, context.Context, string) error
	switch cfg.Strat {
	case config.StratKill:
		action = ContainerManager.Kill
	case config.StratPause:
		action = ContainerManager.Pause
	case config.StratStop:
		action = ContainerManager.Stop
	case config.StratRestart:
		action = ContainerManager.Restart
	}
	return &Tool{
		checkers:         checkers,
		containerManager: containerManager,
		cfg:              cfg,
		logger:           logger,
		action:           action,
	}
}

func (t *Tool) Run(ctx context.Context) error {
	ticker := time.NewTicker(t.cfg.MonitoringPeriod)
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context canceled: %w", ctx.Err())
		case <-ticker.C:
			statuses, err := t.check(ctx)
			if err != nil {
				slog.Error("check statuses", slog.Any("err", err))
			}

			for container, status := range statuses {
				if !status {
					t.logger.Info("killing container", slog.String("container", container))
					if err := t.action(t.containerManager, ctx, container); err != nil {
						t.logger.Error("kill container",
							slog.String("container", container),
							slog.Any("err", err),
						)
					}
				}
			}
		}
	}
}

func (t *Tool) check(ctx context.Context) (map[string]bool, error) {
	statuses := make(map[string]bool)

	containersStats, err := t.containerManager.ContainersStats(ctx)
	if err != nil {
		return nil, err
	}

	for id, stat := range containersStats {
		ok := true
		for _, c := range t.checkers {
			if !c.Check(ctx, stat) {
				ok = false
				break
			}
		}
		statuses[id] = ok
	}

	return statuses, nil
}
