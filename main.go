package main

import (
	"fmt"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/unicoooorn/docker-monitoring-tool/internal/docker"
)

func main() {
	dc, _ := docker.NewDockerClient()
	cs, _ := dc.GetContainers(container.ListOptions{})
	defer dc.Close()
	for i, c := range cs {
		go func() {
			ch := make(chan uint64)
			go dc.MonitorContainerStats(c.ID, 1*time.Second, ch)
			for {
				fmt.Printf("%v %v\n", i, <-ch)
			}
		}()
	}
	a := ""
	fmt.Scan(&a)
}
