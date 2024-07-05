package checker

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/unicoooorn/docker-monitoring-tool/internal/config"
	"github.com/unicoooorn/docker-monitoring-tool/internal/model"
)

type BlackListChecker struct {
	cfg   config.FileMonitoringConfig
	files map[string]time.Time
}

func NewBlackListFileMonitoringChecker(cfg config.FileMonitoringConfig, stats map[string]model.Stat) (*BlackListChecker, error) {
	checker := FileMonitoringChecker{cfg, make(map[string]time.Time)}

	for _, stat := range stats {
		for _, mount := range stat.Volumes {
			if err := filepath.Walk(mount, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					checker.files[path] = info.ModTime()
				}
				return err
			}); err != nil {
				return nil, err
			}
		}
	}
	return &checker, nil
}

func (d *BlackListChecker) Check(_ context.Context, stat model.Stat) bool {
	checkpassed := true

	for _, mount := range stat.Volumes {
		err := filepath.Walk(mount, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				for _, el := range d.cfg.Files {
					if el == path {
						if val, ok := d.Files[path]; ok {
							if val != info.ModTime() {
								checkpassed = false
								return err
							}
						}
					}
				}
			}
			return err
		})
		if err != nil {
			return false
		}
	}
	return checkpassed
}