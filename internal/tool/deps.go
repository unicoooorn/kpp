package tool

import (
	"context"

	"github.com/unicoooorn/docker-monitoring-tool/internal/model"
)

//go:generate mockery --name ContainerManager
type ContainerManager interface {
	ContainersStats(context.Context) (map[string]model.Stat, error)
	Kill(context.Context, string) error
}
