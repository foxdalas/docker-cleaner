package cmd

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Exporter collects metrics from a docker-cleaner daemon.
type Exporter struct {
	up   *prometheus.Desc
	last *prometheus.Desc

	diskUsage         *prometheus.Desc
	diskUsagePercents *prometheus.Desc
	diskReclaimable   *prometheus.Desc
}

const (
	namespace = "docker_cleaner"
)

func NewExporter() *Exporter {
	return &Exporter{
		up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "up"),
			"Could the docker-cleaner server be reached.",
			nil,
			nil,
		),
		last: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "last"),
			"Cleaner last run.",
			nil,
			nil,
		),
		diskUsage: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "disk_usage"),
			"Docker daemon disk usage",
			[]string{"type"},
			nil,
		),
		diskReclaimable: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "disk_reclaimable"),
			"Docker daemon disk reclaimable",
			[]string{"type"},
			nil,
		),
		diskUsagePercents: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "disk_usage_percents"),
			"Docker daemon disk usage.",
			nil,
			nil,
		),
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.up
	ch <- e.last
	ch <- e.diskUsage
	ch <- e.diskUsagePercents
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	status := 1
	ch <- prometheus.MustNewConstMetric(e.up, prometheus.GaugeValue, float64(status))

	mtx.Lock()
	if timeStamp > 0 {
		ch <- prometheus.MustNewConstMetric(e.last, prometheus.CounterValue, float64(usage.Docker.Volumes.Size))
	}

	ch <- prometheus.MustNewConstMetric(e.diskUsage, prometheus.GaugeValue, float64(usage.System.Bytes), "total")
	ch <- prometheus.MustNewConstMetric(e.diskUsage, prometheus.GaugeValue, float64(usage.Docker.BuildCache.Size), "build_cache")
	ch <- prometheus.MustNewConstMetric(e.diskUsage, prometheus.GaugeValue, float64(usage.Docker.Containers.Size), "containers")
	ch <- prometheus.MustNewConstMetric(e.diskUsage, prometheus.GaugeValue, float64(usage.Docker.Volumes.Size), "volumes")
	ch <- prometheus.MustNewConstMetric(e.diskUsage, prometheus.GaugeValue, float64(usage.Docker.Images.Size), "images")

	ch <- prometheus.MustNewConstMetric(e.diskReclaimable, prometheus.GaugeValue, float64(usage.Docker.BuildCache.Reclaimable), "build_cache")
	ch <- prometheus.MustNewConstMetric(e.diskReclaimable, prometheus.GaugeValue, float64(usage.Docker.Containers.Reclaimable), "containers")
	ch <- prometheus.MustNewConstMetric(e.diskReclaimable, prometheus.GaugeValue, float64(usage.Docker.Volumes.Reclaimable), "volumes")
	ch <- prometheus.MustNewConstMetric(e.diskReclaimable, prometheus.GaugeValue, float64(usage.Docker.Images.Reclaimable), "images")

	ch <- prometheus.MustNewConstMetric(e.diskUsagePercents, prometheus.GaugeValue, usage.System.Percents)

	mtx.Unlock()
}
