package checker

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/unicoooorn/kpp/internal/config"
	"github.com/unicoooorn/kpp/internal/model"
)

type WhiteListChecker struct {
	cfg   config.FileMonitoringConfig
	files map[string]time.Time
}

func NewWhiteListFileMonitoringChecker(cfg config.FileMonitoringConfig, stats map[string]model.Stat) (*WhiteListChecker, error) {
	checker := WhiteListChecker{cfg, make(map[string]time.Time)}

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

func (d *WhiteListChecker) Check(_ context.Context, stat model.Stat) bool {
	checkpassed := true

	for _, mount := range stat.Volumes {
		if err := filepath.Walk(mount, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				for _, filename := range d.cfg.Files {
					if filename == path {
						return err
					}
				}

				if modtime, ok := d.files[path]; ok {
					if modtime != info.ModTime() {
						checkpassed = false
						return err
					}
				} else {
					d.files[path] = info.ModTime()
				}
			}
			return err
		}); err != nil {
			return false
		}
	}
	return checkpassed
}
