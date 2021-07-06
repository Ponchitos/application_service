package prometheus

import (
	"github.com/Ponchitos/application_service/server/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type Summary struct {
	summaryVec  *prometheus.SummaryVec
	labelValues []string
}

func (summary *Summary) With(labelValues ...string) metrics.Histogram {
	return &Summary{
		summaryVec:  summary.summaryVec,
		labelValues: append(summary.labelValues, labelValues...),
	}
}

func (summary *Summary) Observe(value float64) {
	summary.summaryVec.With(makeLabels(summary.labelValues...)).Observe(value)
}

func NewSummary(summaryVec *prometheus.SummaryVec) *Summary {
	return &Summary{summaryVec: summaryVec}
}

func NewSummaryFrom(opts prometheus.SummaryOpts, labelNames []string) *Summary {
	summaryVec := prometheus.NewSummaryVec(opts, labelNames)
	prometheus.MustRegister(summaryVec)

	return NewSummary(summaryVec)
}
