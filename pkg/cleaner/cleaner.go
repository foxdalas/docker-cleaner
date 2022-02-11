package cleaner

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/shirou/gopsutil/disk"
	"github.com/sirupsen/logrus"
	"strings"
	"syscall"
	"time"
)

type DiskUsageSummary struct {
	Size        int64
	Reclaimable int64
}

type SystemDiskUsage struct {
	Bytes    int64
	Percents float64
}

type DockerDiskUsage struct {
	BuildCache DiskUsageSummary
	Containers DiskUsageSummary
	Volumes    DiskUsageSummary
	Images     DiskUsageSummary
}

type DiskUsage struct {
	System *SystemDiskUsage
	Docker *DockerDiskUsage
}

type Active struct {
	BuildCache int
	Containers int
	Volumes    int
}

type Cleaner struct {
	Docker *client.Client
	Ctx    context.Context
	Log    *logrus.Entry
	Dir    string
	TTL    time.Duration
}

func New(client *client.Client, ctx context.Context, log *logrus.Entry, dir string, ttl time.Duration) *Cleaner {
	return &Cleaner{
		Docker: client,
		Ctx:    ctx,
		Log:    log,
		Dir:    dir,
		TTL:    ttl,
	}
}

func (Cleaner *Cleaner) GetUsageInfo() (*DiskUsage, error) {

	system, err := Cleaner.DeviceSpaceUsage(Cleaner.Dir)
	if err != nil {
		return nil, fmt.Errorf("device %s: %s", Cleaner.Dir, err)
	}

	docker, err := Cleaner.DockerDiskUsage()
	if err != nil {
		return nil, err
	}

	return &DiskUsage{
		System: system,
		Docker: docker,
	}, nil

}

// GB, Percent, Error
func (Cleaner *Cleaner) DeviceSpaceUsage(device string) (*SystemDiskUsage, error) {
	usage, err := disk.Usage(device)
	if err != nil {
		return &SystemDiskUsage{}, err
	}
	return &SystemDiskUsage{
		Bytes:    int64(usage.Used),
		Percents: usage.UsedPercent,
	}, err
}

func (Cleaner *Cleaner) GetDiskUtilization(path string) (uint64, float64, error) {
	var stat syscall.Statfs_t
	err := syscall.Statfs(path, &stat)
	if err != nil {
		return 0, 0, err
	}
	all := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bfree * uint64(stat.Bsize)
	used := all - free
	percentageUtilized := float64(used) / float64(all) * float64(100)
	return used, float64(percentageUtilized), err
}

func (Cleaner *Cleaner) DockerDiskUsage() (*DockerDiskUsage, error) {
	var buildCacheSummary, containersSummary, volumesSummary, imagesSummary DiskUsageSummary

	diskUsage, err := Cleaner.Docker.DiskUsage(Cleaner.Ctx)
	if err != nil {
		return nil, err
	}

	var inUseBytes int64
	for _, buildCache := range diskUsage.BuildCache {
		if !buildCache.Shared {
			buildCacheSummary.Size += buildCache.Size
			if buildCache.InUse && !buildCache.Shared {
				inUseBytes += buildCache.Size
			}
		}
	}
	buildCacheSummary.Reclaimable = buildCacheSummary.Size - inUseBytes

	for _, container := range diskUsage.Containers {
		containersSummary.Size += container.SizeRw
		if !containerIsActive(*container) {
			containersSummary.Reclaimable += container.SizeRw
		}
	}

	for _, volume := range diskUsage.Volumes {
		if volume.UsageData.Size != -1 {
			if volume.UsageData.RefCount == 0 {
				volumesSummary.Reclaimable += volume.UsageData.Size
			}
			volumesSummary.Size += volume.UsageData.Size
		}
	}

	for _, i := range diskUsage.Images {
		if i.Containers != 0 {
			if i.VirtualSize == -1 || i.SharedSize == -1 {
				continue
			}
			imagesSummary.Reclaimable += i.VirtualSize - i.SharedSize
		}
	}
	imagesSummary.Size = diskUsage.LayersSize

	return &DockerDiskUsage{
		BuildCache: buildCacheSummary,
		Containers: containersSummary,
		Volumes:    volumesSummary,
		Images:     imagesSummary,
	}, err
}

func (Cleaner *Cleaner) BuildCachePrune() (uint64, error) {
	filter := filters.NewArgs(filters.KeyValuePair{
		Key:   "until",
		Value: Cleaner.TTL.String(),
	})
	opts := types.BuildCachePruneOptions{
		All:     true,
		Filters: filter,
	}

	report, err := Cleaner.Docker.BuildCachePrune(Cleaner.Ctx, opts)
	if err != nil {
		return 0, err
	}
	return report.SpaceReclaimed, nil
}

func (Cleaner *Cleaner) ContainersPrune() (uint64, error) {
	filter := filters.NewArgs(filters.KeyValuePair{
		Key:   "until",
		Value: Cleaner.TTL.String(),
	})

	report, err := Cleaner.Docker.ContainersPrune(Cleaner.Ctx, filter)
	if err != nil {
		return 0, err
	}
	return report.SpaceReclaimed, nil
}

func (Cleaner *Cleaner) VolumesPrune() (uint64, error) {
	filter := filters.NewArgs()

	report, err := Cleaner.Docker.VolumesPrune(Cleaner.Ctx, filter)
	if err != nil {
		return 0, err
	}
	return report.SpaceReclaimed, nil
}

func (Cleaner *Cleaner) ImagesPrune() (uint64, error) {
	filter := filters.NewArgs(filters.KeyValuePair{
		Key:   "until",
		Value: Cleaner.TTL.String(),
	})
	report, err := Cleaner.Docker.ImagesPrune(Cleaner.Ctx, filter)
	if err != nil {
		return 0, err
	}
	return report.SpaceReclaimed, nil
}

func containerIsActive(container types.Container) bool {
	return strings.Contains(container.State, "running") ||
		strings.Contains(container.State, "paused") ||
		strings.Contains(container.State, "restarting")
}
