package prometheus

import (
	"github.com/Ponchitos/application_service/server/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type Counter struct {
	counterVec  *prometheus.CounterVec
	labelValues []string
}

func (counter *Counter) With(labelValues ...string) metrics.Counter {
	return &Counter{
		counterVec:  counter.counterVec,
		labelValues: append(counter.labelValues, labelValues...),
	}
}

func (counter *Counter) Add(delta float64) {
	counter.counterVec.With(makeLabels(counter.labelValues...)).Add(delta)
}

func NewCounter(counterVec *prometheus.CounterVec) *Counter {
	return &Counter{counterVec: counterVec}
}

func NewCounterFrom(opts prometheus.CounterOpts, labelNames []string) *Counter {
	counterVec := prometheus.NewCounterVec(opts, labelNames)
	prometheus.MustRegister(counterVec)

	return NewCounter(counterVec)
}
