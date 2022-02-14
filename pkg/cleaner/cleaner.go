package cleaner

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/shirou/gopsutil/disk"
	"github.com/sirupsen/logrus"
)

type DockerDiskUsage struct {
	BuildCache float64
	Containers float64
	Volumes    float64
}

type Usage struct {
	DiskUsageGb       float64
	DiskUsagePercents float64
	BuildCacheUsage   float64
	ContainerUsage    float64
	VolumesUsage      float64
}

type Cleaner struct {
	Docker *client.Client
	Ctx    context.Context
	Log    *logrus.Entry
	Device string
}

func New(client *client.Client, ctx context.Context, log *logrus.Entry, device string) *Cleaner {
	return &Cleaner{
		Docker: client,
		Ctx:    ctx,
		Log:    log,
		Device: device,
	}
}

func (Cleaner *Cleaner) GetUsageInfo() (*Usage, error) {
	diskSpaceUsageGB, diskSpaceUsagePercents, err := Cleaner.DeviceSpaceUsage(Cleaner.Device)
	if err != nil {
		return nil, fmt.Errorf("device %s: %s", Cleaner.Device, err)
	}

	dockerSpaceUsage, err := Cleaner.DockerDiskUsage()
	if err != nil {
		return nil, err
	}

	return &Usage{
		DiskUsageGb:       diskSpaceUsageGB,
		DiskUsagePercents: diskSpaceUsagePercents,
		BuildCacheUsage:   dockerSpaceUsage.BuildCache,
		ContainerUsage:    dockerSpaceUsage.Containers,
		VolumesUsage:      dockerSpaceUsage.Volumes,
	}, nil

}

// GB, Percent, Error
func (Cleaner *Cleaner) DeviceSpaceUsage(device string) (float64, float64, error) {
	usage, err := disk.Usage(device)
	if err != nil {
		return 0, 0, err
	}
	return float64(usage.Used) / (1 << 30), float64(usage.UsedPercent), err
}

func (Cleaner *Cleaner) DockerDiskUsage() (*DockerDiskUsage, error) {
	var buildCacheUsage float64
	var containersUsage float64
	var volumesUsage float64

	diskUsage, err := Cleaner.Docker.DiskUsage(Cleaner.Ctx)
	if err != nil {
		return nil, err
	}

	for _, buildCache := range diskUsage.BuildCache {
		buildCacheUsage += float64(buildCache.Size)
	}
	for _, container := range diskUsage.Containers {
		containersUsage += float64(container.SizeRootFs)
	}
	for _, volume := range diskUsage.Volumes {
		volumesUsage += float64(volume.UsageData.Size)
	}

	return &DockerDiskUsage{
		BuildCache: buildCacheUsage / (1 << 30),
		Containers: containersUsage / (1 << 30),
		Volumes:    volumesUsage / (1 << 30),
	}, err
}

func (Cleaner *Cleaner) BuildCachePrune() (float64, error) {
	opts := types.BuildCachePruneOptions{
		All: true,
	}

	report, err := Cleaner.Docker.BuildCachePrune(Cleaner.Ctx, opts)
	if err != nil {
		return 0, err
	}
	return float64(report.SpaceReclaimed) / (1 << 30), nil
}

func (Cleaner *Cleaner) ContainersPrune() (float64, error) {
	filter := filters.NewArgs(filters.KeyValuePair{
		Key:   "until",
		Value: "48h",
	})

	report, err := Cleaner.Docker.ContainersPrune(Cleaner.Ctx, filter)
	if err != nil {
		return 0, err
	}
	return float64(report.SpaceReclaimed) / (1 << 30), nil
}

func (Cleaner *Cleaner) VolumesPrune() (float64, error) {
	filter := filters.NewArgs()

	report, err := Cleaner.Docker.VolumesPrune(Cleaner.Ctx, filter)
	if err != nil {
		return 0, err
	}
	return float64(report.SpaceReclaimed) / (1 << 30), nil
}
