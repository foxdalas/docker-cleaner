package main

import (
	"context"
	"docker-cleaner/pkg/cleaner"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	cleaner := cleaner.New(cli, ctx)

	diskUsagedGb, diskUsagePercent, _ := cleaner.DeviceSpaceUsage("/dev/disk1s5s1")
	logrus.Infof("GB Usage: %f", diskUsagedGb)
	logrus.Infof("Percent Usage: %f", diskUsagePercent)

	dDiskUsage, err := cleaner.DockerDiskUsage()
	if err != nil {
		logrus.Error(err)
	}
	logrus.Infof("BuildCache Usage: %f", dDiskUsage.BuildCache)
	logrus.Infof("Containers Usage: %f", dDiskUsage.Containers)
	logrus.Infof("Volumes Usage: %f", dDiskUsage.Volumes)

	if diskUsagePercent > 50.0 {
		logrus.Infof("Disk usage > 50 percents (%f usage), starting cleanup", diskUsagePercent)
		if dDiskUsage.BuildCache > 0 {
			reclamed, err := cleaner.BuildCachePrune()
			if err != nil {
				logrus.Errorf("%s", err)
			}
			logrus.Infof("Build cache prune reclamed Gb: %f", reclamed)
		}
		if dDiskUsage.Containers > 0 {
			reclamed, err := cleaner.ContainersPrune()
			if err != nil {
				logrus.Errorf("%s", err)
			}
			logrus.Infof("Containers prune reclamed Gb: %f", reclamed)
		}
		if dDiskUsage.Volumes > 0 {
			reclamed, err := cleaner.VolumesPrune()
			if err != nil {
				logrus.Errorf("%s", err)
			}
			logrus.Infof("Volumes prune reclamed Gb: %f", reclamed)
		}
	}
}
