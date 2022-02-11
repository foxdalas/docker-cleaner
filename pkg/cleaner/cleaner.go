package cleaner

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/shirou/gopsutil/disk"
	"github.com/sirupsen/logrus"
	"syscall"
	"time"
)

type DockerDiskUsage struct {
	BuildCache int64
	Containers int64
	Volumes    int64
	Images     int64
}

type Active struct {
	BuildCache int
	Containers int
	Volumes    int
}

type Usage struct {
	DiskUsage         int64
	DiskUsagePercents float64
	BuildCacheUsage   int64
	ContainerUsage    int64
	VolumesUsage      int64
	ImagesUsage       int64
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

func (Cleaner *Cleaner) GetUsageInfo() (*Usage, error) {
	diskSpaceUsage, diskSpaceUsagePercents, err := Cleaner.DeviceSpaceUsage(Cleaner.Dir)
	if err != nil {
		return nil, fmt.Errorf("device %s: %s", Cleaner.Dir, err)
	}

	dockerSpaceUsage, err := Cleaner.DockerDiskUsage()
	if err != nil {
		return nil, err
	}

	return &Usage{
		DiskUsage:         diskSpaceUsage,
		DiskUsagePercents: diskSpaceUsagePercents,
		BuildCacheUsage:   dockerSpaceUsage.BuildCache,
		ContainerUsage:    dockerSpaceUsage.Containers,
		VolumesUsage:      dockerSpaceUsage.Volumes,
		ImagesUsage:       dockerSpaceUsage.Images,
	}, nil

}

// GB, Percent, Error
func (Cleaner *Cleaner) DeviceSpaceUsage(device string) (int64, float64, error) {
	usage, err := disk.Usage(device)
	if err != nil {
		return 0, 0, err
	}
	return int64(usage.Used), float64(usage.UsedPercent), err
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
	var buildCacheUsage, containersUsage, volumesUsage, imagesUsage int64

	diskUsage, err := Cleaner.Docker.DiskUsage(Cleaner.Ctx)
	if err != nil {
		return nil, err
	}

	for _, buildCache := range diskUsage.BuildCache {
		if !buildCache.Shared {
			buildCacheUsage += buildCache.Size
		}
	}
	for _, container := range diskUsage.Containers {
		containersUsage += container.SizeRootFs
	}
	for _, volume := range diskUsage.Volumes {
		volumesUsage += volume.UsageData.Size
	}

	//for _, image := range diskUsage.Images {
	//
	//	if image.Containers > 0 {
	//		fmt.Println(image.ID)
	//	}
	//	imagesUsage += image. - image.SharedSize
	//}
	imagesUsage = diskUsage.LayersSize

	return &DockerDiskUsage{
		BuildCache: buildCacheUsage,
		Containers: containersUsage,
		Volumes:    volumesUsage,
		Images:     imagesUsage,
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
