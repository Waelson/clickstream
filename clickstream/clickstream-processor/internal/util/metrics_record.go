package util

import (
	"github.com/prometheus/client_golang/prometheus"
)

type MetricsRecord interface {
	ResetGaugeCampaign()
	WithLabelValues(campaignName string, id string, value float64)
}

type metricsRecord struct {
	gaugeCampaign *prometheus.GaugeVec
}

func (m *metricsRecord) WithLabelValues(campaignName string, id string, value float64) {
	m.gaugeCampaign.WithLabelValues(campaignName).Set(value)
}

func (m *metricsRecord) ResetGaugeCampaign() {
	m.gaugeCampaign.Reset()
}

func NewMetricsRecord() MetricsRecord {
	campaignGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "campaign",
			Help: "Number of link clicks per campaign",
		},
		[]string{"campaign_name"},
	)

	prometheus.MustRegister(campaignGauge)

	return &metricsRecord{
		gaugeCampaign: campaignGauge,
	}
}
