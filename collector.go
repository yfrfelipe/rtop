package main

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/crypto/ssh"
)

type localCollector struct {
	totalMem *prometheus.Desc
	freeMem  *prometheus.Desc
	totalCpu *prometheus.Desc
	ipv4Tx   *prometheus.Desc
	ipv4Rx   *prometheus.Desc

	client *ssh.Client
}

func newCollector(sshClient *ssh.Client) *localCollector {
	return &localCollector{
		totalMem: prometheus.NewDesc("tz_total_mem", "Total memory.", nil, nil),
		freeMem:  prometheus.NewDesc("tz_free_mem", "Total free memory.", nil, nil),
		totalCpu: prometheus.NewDesc("tz_total_cpu", "Total CPU.", nil, nil),
		ipv4Rx:   prometheus.NewDesc("tz_ipv4_rx", "Total IPv4 Rx.", nil, nil),
		ipv4Tx:   prometheus.NewDesc("tz_ipv4_tx", "Total IPv4 Tx.", nil, nil),
		client:   sshClient,
	}
}

func (col *localCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- col.totalMem
}

func (col *localCollector) Collect(ch chan<- prometheus.Metric) {

	fmt.Println("Collecting")

	stats := make(chan Stats)
	go collectStats(stats, col.client)

	result := <-stats

	ch <- prometheus.MustNewConstMetric(col.totalMem, prometheus.GaugeValue, float64(result.MemTotal))
	ch <- prometheus.MustNewConstMetric(col.freeMem, prometheus.GaugeValue, float64(result.MemFree))
	ch <- prometheus.MustNewConstMetric(col.totalCpu, prometheus.GaugeValue, float64(result.CPU.System))
}
