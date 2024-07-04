package tool

import (
	"context"

	"github.com/unicoooorn/docker-monitoring-tool/internal/model"
)

//go:generate mockery --name ContainerManager
type ContainerManager interface {
	ContainersStats(context.Context) (map[string]model.Stat, error)
	Kill(context.Context, string) error
	Stop(context.Context, string) error
	Pause(context.Context, string) error
	Restart(context.Context, string) error
}

//go:generate mockery --name Checker
type Checker interface {
	Check(context.Context, model.Stat) bool
}
