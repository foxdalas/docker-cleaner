package cmd

import (
	"context"
	"flag"
	"fmt"
	"github.com/docker/docker/client"
	"github.com/foxdalas/docker-cleaner/pkg/cleaner"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
	"time"
)

type flags struct {
	dockerDevice         string
	dockerSpaceThreshold float64
	interval             time.Duration
	exporterHost         string
	exporterPort         int
}

func params() (*flags, error) {
	flags := &flags{}

	flag.StringVar(&flags.dockerDevice, "docker.device", "/dev/sda1", "Docker block device.")
	flag.Float64Var(&flags.dockerSpaceThreshold, "docker.threshold", 50.0, "Docker used space threshold in percents")
	flag.DurationVar(&flags.interval, "cleaner.interval", 15*time.Second, "Cleaner interval")
	flag.StringVar(&flags.exporterHost, "exporter.host", "localhost", "Exporter host")
	flag.IntVar(&flags.exporterPort, "exporter.port", 9203, "Exporter port")

	flag.Parse()
	return flags, nil
}

func makeLog() *log.Entry {
	logtype := strings.ToLower(os.Getenv("LOG_TYPE"))
	if logtype == "" {
		logtype = "text"
	}
	if logtype == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	} else if logtype == "text" {
		log.SetFormatter(&log.TextFormatter{
			ForceColors: true,
		})
	} else {
		log.WithField("logtype", logtype).Fatal("Given logtype was not valid, check LOG_TYPE configuration")
		os.Exit(1)
	}

	loglevel := strings.ToLower(os.Getenv("LOG_LEVEL"))
	if len(loglevel) == 0 {
		log.SetLevel(log.InfoLevel)
	} else if loglevel == "debug" {
		log.SetLevel(log.DebugLevel)
	} else if loglevel == "info" {
		log.SetLevel(log.InfoLevel)
	} else if loglevel == "warn" {
		log.SetLevel(log.WarnLevel)
	} else if loglevel == "error" {
		log.SetLevel(log.ErrorLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	return log.WithField("context", "nodeup")
}

func cleanup(threshold float64, device string, interval time.Duration) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	cleaner := cleaner.New(cli, ctx, makeLog(), device)

	for {
		_, err = cli.Ping(ctx)
		if err != nil {
			return err
		}

		usage, err := cleaner.GetUsageInfo()
		if err != nil {
			return err
		}

		log.Infof("Current disk usage is %f GB and %f", usage.DiskUsageGb, usage.DiskUsagePercents)
		if usage.DiskUsagePercents > threshold {
			log.Infof("Disk usage > %f percents (%f usage), starting cleanup", threshold, usage.DiskUsagePercents)
			if usage.BuildCacheUsage > 0 {
				reclamed, err := cleaner.BuildCachePrune()
				if err != nil {
					return err
				}
				log.Infof("Build cache prune reclamed Gb: %f", reclamed)
			}
			if usage.ContainerUsage > 0 {
				reclamed, err := cleaner.ContainersPrune()
				if err != nil {
					return err
				}
				log.Infof("Containers prune reclamed Gb: %f", reclamed)
			}
			if usage.VolumesUsage > 0 {
				reclamed, err := cleaner.VolumesPrune()
				if err != nil {
					return err
				}
				log.Infof("Volumes prune reclaimed Gb: %f", reclamed)
			}
		}
		log.Infof("Waiting for interval: %s", interval)
		time.Sleep(interval)
	}
}

func Run() {
	flags, err := params()
	if err != nil {
		log.Fatal(err)
	}

	// Cleanup
	go func(threshold float64, device string, interval time.Duration) {
		err = cleanup(threshold, device, interval)
		if err != nil {
			log.Error(err)
		}
	}(flags.dockerSpaceThreshold, flags.dockerDevice, flags.interval)

	http.Handle("/metrics", promhttp.Handler())
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", flags.exporterHost, flags.exporterPort), nil)
	if err != nil {
		log.Fatal(err)
	}
}
