package checker

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/unicoooorn/docker-monitoring-tool/internal/config"
	"github.com/unicoooorn/docker-monitoring-tool/internal/model"
)

type FileMonitoringChecker struct {
	cfg   config.FileMonitoringConfig
	Files map[string]time.Time
}

func NewFileMonitoringChecker(cfg config.FileMonitoringConfig, stat model.Stat) (*FileMonitoringChecker, error) {
	checker := FileMonitoringChecker{cfg, make(map[string]time.Time)}
	for _, mount := range stat.Volumes {
		err := filepath.Walk(mount, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				checker.Files[path] = info.ModTime()
			}
			return err
		})
		if err != nil {
			return nil, err
		}
	}
	return &checker, nil
}

func (d *FileMonitoringChecker) Check(_ context.Context, stat model.Stat) (bool, error) {
	flag := true
	if d.cfg.Type == config.WhitelistMode {
		for _, mount := range stat.Volumes {
			err := filepath.Walk(mount, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					for _, el := range d.cfg.Files {
						if el == path {
							return err
						}
					}

					if val, ok := d.Files[path]; ok {
						if val != info.ModTime() {
							flag = false
							return err
						}
					}

					//mb save modtime
				}
				return err
			})
			if err != nil {
				return false, err
			}
		}
	} else {
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
									flag = false
									return err
								}
							}
						}
					}
				}
				return err
			})
			if err != nil {
				return false, err
			}
		}
	}
	return flag, nil
}
