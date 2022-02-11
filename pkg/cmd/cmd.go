package cmd

import (
	"context"
	"flag"
	"fmt"
	"github.com/docker/docker/client"
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

var usage = &cleaner.DiskUsage{}
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

	flag.StringVar(&flags.dockerDir, "docker.dir", "/var/lib/docker", "Docker block device.")
	flag.Float64Var(&flags.dockerSpaceThreshold, "docker.threshold", 50.0, "Docker used space threshold in percents")
	flag.DurationVar(&flags.dockerTTL, "docker.ttl", 48*time.Hour, "Docker TTL")
	flag.DurationVar(&flags.interval, "cleaner.interval", 15*time.Second, "Cleaner interval")
	flag.StringVar(&flags.exporterHost, "exporter.host", "0.0.0.0", "Exporter host")
	flag.IntVar(&flags.exporterPort, "exporter.port", 9203, "Exporter port")
	flag.DurationVar(&flags.exporterTimeout, "exporter.timeout", 15*time.Second, "Exporter timeout")
	flag.StringVar(&flags.exporterMetricsPath, "exporter.telemetry-path", "/metrics", "Exporter path under which to expose metrics.")

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

func cleanup(threshold float64, dir string, ttl time.Duration, interval time.Duration, log *log.Entry) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	cleaner := cleaner.New(cli, ctx, log, dir, ttl)

	for {
		_, err = cli.Ping(ctx)
		if err != nil {
			return err
		}

		usageData, err := cleaner.GetUsageInfo()
		if err != nil {
			return err
		}

		mtx.Lock()
		usage = usageData
		mtx.Unlock()

		if usage.System.Percents > threshold {
			mtx.Lock()
			timeStamp = time.Now().Unix()
			mtx.Unlock()
			log.Infof("Disk usage more then %.1f percents (%.1f usage), starting cleanup", threshold, usage.System.Percents)
			if usage.Docker.BuildCache.Size > 0 {
				reclamed, err := cleaner.BuildCachePrune()
				if err != nil {
					return err
				}
				log.Infof("Build cache prune reclamed Gb: %.1f", float64(reclamed)/(1<<30))
			}
			if usage.Docker.Containers.Size > 0 {
				reclamed, err := cleaner.ContainersPrune()
				if err != nil {
					return err
				}
				log.Infof("Containers prune reclamed Gb: %.1f", float64(reclamed)/(1<<30))
			}
			if usage.Docker.Volumes.Size > 0 {
				reclaimed, err := cleaner.VolumesPrune()
				if err != nil {
					return err
				}
				log.Infof("Volumes prune reclaimed Gb: %.1f", float64(reclaimed)/(1<<30))
			}
			if usage.Docker.Images.Size > 0 {
				reclaimed, err := cleaner.ImagesPrune()
				if err != nil {
					return err
				}
				log.Infof("Image prune reclaimed Gb: %.1f", float64(reclaimed)/(1<<30))
			}
		}
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
		err = cleanup(threshold, dir, ttl, interval, log)
		if err != nil {
			log.Error(err)
		}
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
