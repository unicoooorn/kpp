package checker

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/unicoooorn/kpp/internal/config"
	"github.com/unicoooorn/kpp/internal/model"
)

type BlackListChecker struct {
	cfg   config.FileMonitoringConfig
	files map[string]time.Time
}

func NewBlackListFileMonitoringChecker(cfg config.FileMonitoringConfig, stats map[string]model.Stat) (*BlackListChecker, error) {
	checker := BlackListChecker{cfg, make(map[string]time.Time)}

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
						if val, ok := d.files[path]; ok {
							if val != info.ModTime() {
								checkpassed = false
								return err
							}
						}
					} else if strings.HasPrefix(path, el) {
						checkpassed = false
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
