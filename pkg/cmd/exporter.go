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
		ch <- prometheus.MustNewConstMetric(e.last, prometheus.CounterValue, float64(usage.VolumesUsage))
	}

	ch <- prometheus.MustNewConstMetric(e.diskUsage, prometheus.GaugeValue, float64(usage.DiskUsage), "total")
	ch <- prometheus.MustNewConstMetric(e.diskUsage, prometheus.GaugeValue, float64(usage.BuildCacheUsage), "build_cache")
	ch <- prometheus.MustNewConstMetric(e.diskUsage, prometheus.GaugeValue, float64(usage.ContainerUsage), "containers")
	ch <- prometheus.MustNewConstMetric(e.diskUsage, prometheus.GaugeValue, float64(usage.VolumesUsage), "volumes")

	ch <- prometheus.MustNewConstMetric(e.diskUsagePercents, prometheus.GaugeValue, usage.DiskUsagePercents)

	mtx.Unlock()
}
