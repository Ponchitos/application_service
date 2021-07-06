package prometheus

import (
	"github.com/Ponchitos/application_service/server/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type Histogram struct {
	histogramVec *prometheus.HistogramVec
	labelValues  []string
}

func (histogram *Histogram) With(labelValues ...string) metrics.Histogram {
	return &Histogram{
		histogramVec: histogram.histogramVec,
		labelValues:  append(histogram.labelValues, labelValues...),
	}
}

func (histogram *Histogram) Observe(value float64) {
	histogram.histogramVec.With(makeLabels(histogram.labelValues...)).Observe(value)
}

func NewHistogram(histogramVec *prometheus.HistogramVec) *Histogram {
	return &Histogram{histogramVec: histogramVec}
}

func NewHistogramFrom(opts prometheus.HistogramOpts, labelNames []string) *Histogram {
	histogramVec := prometheus.NewHistogramVec(opts, labelNames)
	prometheus.MustRegister(histogramVec)

	return NewHistogram(histogramVec)
}
