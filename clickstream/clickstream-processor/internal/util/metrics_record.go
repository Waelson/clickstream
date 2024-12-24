package util

import (
	"github.com/prometheus/client_golang/prometheus"
)

type MetricsRecord interface {
	ResetGaugeFeatureFlag()
	WithLabelValues(campaignName string, id string, value float64)
}

type metricsRecord struct {
	gaugeFeatureFlag *prometheus.GaugeVec
}

func (m *metricsRecord) WithLabelValues(campaignName string, id string, value float64) {
	m.gaugeFeatureFlag.WithLabelValues(campaignName, id).Set(value)
}

func (m *metricsRecord) ResetGaugeFeatureFlag() {
	m.gaugeFeatureFlag.Reset()
}

func NewMetricsRecord() MetricsRecord {
	campaignGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "campaign_click_count",
			Help: "Number of link clicks per campaign",
		},
		[]string{"campaign_id", "click_count"},
	)

	prometheus.MustRegister(campaignGauge)

	return &metricsRecord{
		gaugeFeatureFlag: campaignGauge,
	}
}
