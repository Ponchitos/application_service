package prometheus

import (
	"github.com/Ponchitos/application_service/server/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type Gauge struct {
	gaugeVec    *prometheus.GaugeVec
	labelValues []string
}

func (gauge *Gauge) With(labelValues ...string) metrics.Gauge {
	return &Gauge{
		gaugeVec:    gauge.gaugeVec,
		labelValues: append(gauge.labelValues, labelValues...),
	}
}

func (gauge *Gauge) Set(value float64) {
	gauge.gaugeVec.With(makeLabels(gauge.labelValues...)).Set(value)
}

func (gauge *Gauge) Add(delta float64) {
	gauge.gaugeVec.With(makeLabels(gauge.labelValues...)).Add(delta)
}

func NewGauge(gaugeVec *prometheus.GaugeVec) *Gauge {
	return &Gauge{gaugeVec: gaugeVec}
}

func NewGaugeFrom(opts prometheus.GaugeOpts, labelNames []string) *Gauge {
	gaugeVec := prometheus.NewGaugeVec(opts, labelNames)
	prometheus.MustRegister(gaugeVec)

	return NewGauge(gaugeVec)
}
