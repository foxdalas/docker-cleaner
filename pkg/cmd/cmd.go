package cmd

import (
	"context"
	"flag"
	"fmt"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	units "github.com/docker/go-units"
	"github.com/foxdalas/docker-cleaner/pkg/cleaner"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var usage = &cleaner.DiskUsage{
	Docker: &cleaner.DockerDiskUsage{},
	System: &cleaner.SystemDiskUsage{},
}
var timeStamp int64

var mtx sync.Mutex

type flags struct {
	dockerDir            string
	dockerSpaceThreshold float64
	dockerTTL            time.Duration
	interval             time.Duration
	exporterHost         string
	exporterPort         int
	exporterTimeout      time.Duration
	exporterMetricsPath  string
}

func params() (*flags, error) {
	flags := &flags{}

	flag.StringVar(&flags.dockerDir, "docker.dir", "/var/lib/docker", "Docker storage directory")
	flag.Float64Var(&flags.dockerSpaceThreshold, "docker.threshold", 50.0, "Docker volume usage threshold")
	flag.DurationVar(&flags.dockerTTL, "docker.ttl", 48*time.Hour, "Docker volumes TTL. Same until=48h")
	flag.DurationVar(&flags.interval, "cleaner.interval", 15*time.Second, "Cleaner check interval")
	flag.StringVar(&flags.exporterHost, "exporter.host", "0.0.0.0", "Docker cleaner exporter listen host")
	flag.IntVar(&flags.exporterPort, "exporter.port", 9203, "Docker cleaner exporter listen port")
	flag.DurationVar(&flags.exporterTimeout, "exporter.timeout", 15*time.Second, "Docker cleaner exporter timeout")
	flag.StringVar(&flags.exporterMetricsPath, "exporter.telemetry-path", "/metrics", "Docker cleaner exporter path under which to expose metrics.")

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
	return log.WithField("context", "docker-cleaner")
}

func cleanup(threshold float64, dir string, ttl time.Duration, interval time.Duration, log *log.Entry) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Error(err)
	}
	cleaner := cleaner.New(cli, ctx, log, dir, ttl)

	for {
		_, err = cli.Ping(ctx)
		if err != nil {
			log.Fatal(err)
		}

		usageData, err := cleaner.GetUsageInfo()
		if err != nil {
			log.Error(err)
		}

		mtx.Lock()
		usage = usageData
		mtx.Unlock()

		if usage.System.Percents > threshold {
			mtx.Lock()
			timeStamp = time.Now().Unix()
			mtx.Unlock()
			log.Infof("Disk usage more then %.1f percents (%.1f usage), starting cleanup", threshold, usage.System.Percents)
			if usage.Docker.BuildCache.Reclaimable > 0 {
				log.Infof("BuildCache estimate reclaimable space: %s", units.HumanSize(float64(usage.Docker.BuildCache.Reclaimable)))
				reclaimed, err := cleaner.BuildCachePrune()
				if err != nil {
					log.Error(err)
				}
				log.Infof("BuildCache prune reclaimed: %s", units.HumanSize(float64(reclaimed)))
			}
			if usage.Docker.Containers.Reclaimable > 0 {
				log.Infof("Containers estimate reclaimable space: %s", units.HumanSize(float64(usage.Docker.Containers.Reclaimable)))
				reclaimed, err := cleaner.ContainersPrune()
				if err != nil {
					log.Error(err)
				}
				log.Infof("Containers prune reclaimed: %s", units.HumanSize(float64(reclaimed)))
			}
			if usage.Docker.Volumes.Reclaimable > 0 {
				log.Infof("Volumes estimate reclaimable space: %s", units.HumanSize(float64(usage.Docker.Volumes.Reclaimable)))
				reclaimed, err := cleaner.VolumesPrune()
				if err != nil {
					log.Error(err)
				}
				log.Infof("Volumes prune reclaimed: %s", units.HumanSize(float64(reclaimed)))
			}
			if usage.Docker.Images.Reclaimable > 0 {
				var reclaimed uint64

				log.Infof("Images estimate reclaimable space: %s", units.HumanSize(float64(usage.Docker.Images.Reclaimable)))
				log.Infoln("Prune dangling images")
				dangling := filters.NewArgs(
					filters.KeyValuePair{
						Key:   "until",
						Value: cleaner.TTL.String(),
					},
					filters.KeyValuePair{
						Key:   "dangling",
						Value: "true",
					},
				)
				dangling_reclaimed, err := cleaner.ImagesPrune(dangling)
				if err != nil {
					log.Error(err)
				}

				log.Infof("Prune erected images older then %s", cleaner.TTL.String())

				erected := filters.NewArgs(
					filters.KeyValuePair{
						Key:   "until",
						Value: cleaner.TTL.String(),
					},
					filters.KeyValuePair{
						Key:   "dangling",
						Value: "false",
					},
				)
				erected_reclaimed, err := cleaner.ImagesPrune(erected)
				if err != nil {
					log.Error(err)
				}
				reclaimed = dangling_reclaimed + erected_reclaimed

				log.Infof("Images prune reclaimed: %s", units.HumanSize(float64(reclaimed)))
			}
		}

		log.Infoln("Prune unused networks")
		cleanedNetworksCount, err := cleaner.NetworksPrune()
		if err != nil {
			log.Error(err)
		}
		log.Infof("Removed %d unused networks", cleanedNetworksCount)

		log.Infof("Waiting for interval: %s", interval)
		time.Sleep(interval)
	}
}

func Run() {
	log := makeLog()
	log.Infoln("Starting docker-cleaner", version.Info())
	log.Infoln("Build context", version.BuildContext())

	flags, err := params()
	if err != nil {
		log.Fatal(err)
	}

	// Cleanup
	go func(threshold float64, dir string, ttl time.Duration, interval time.Duration) {
		cleanup(threshold, dir, ttl, interval, log)
	}(flags.dockerSpaceThreshold, flags.dockerDir, flags.dockerTTL, flags.interval)

	prometheus.MustRegister(NewExporter())

	http.Handle(flags.exporterMetricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`<html>
      <head><title>Docker cleaner Exporter</title></head>
      <body>
      <h1>Docker cleaner Exporter</h1>
      <p><a href='` + flags.exporterMetricsPath + `'>Metrics</a></p>
      </body>
      </html>`))
		if err != nil {
			log.Error(err)
		}
	})

	err = http.ListenAndServe(fmt.Sprintf("%s:%d", flags.exporterHost, flags.exporterPort), nil)
	if err != nil {
		log.Fatal(err)
	}
}
