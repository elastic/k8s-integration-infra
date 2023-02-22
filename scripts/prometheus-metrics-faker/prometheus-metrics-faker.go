package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	interval = 9
)

func main() {

	gauges := flag.Int("gauge", 0, "number of metrics of type gauge")
	counters := flag.Int("counter", 0, "number of metrics of type counter")
	histograms := flag.Int("histogram", 0, "number of metrics of type histogram")
	// summaries := flag.Int("summary", 0, "number of metrics of type summary")
	labels := flag.Int("labels", 1, "number of labels per metrics")
	cardinality := flag.Int("cardinality", 1, "number of different label value per metrics")
	flag.Parse()

	constPromLabels := prometheus.Labels{
		"job": "prometheus-metrics-faker",
	}

	var labelsList []string

	for j := 1; j <= *labels; j++ {
		// promLabels[fmt.Sprintf("key%d", i)] = fmt.Sprintf("value%d", rand.Intn(100))
		labelsList = append(labelsList, fmt.Sprintf("key%d", j))
		// valuesList = append(valuesList, fmt.Sprintf("value%d", rand.Intn(100)))
	}
	switch {
	case *gauges > 0:
		gaugeInit("gauge", *gauges, labelsList, *cardinality, constPromLabels)
	case *counters > 0:
		counterInit("counter", *counters, labelsList, *cardinality, constPromLabels)
	case *histograms > 0:
		histogramInit("histogram", *histograms, labelsList, *cardinality, constPromLabels)
		// case *summaries > 0:
		// 	metricsInit("summary", *summaries, *labels)
	}

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9000", nil)
}

func counterInit(metricType string, metricsNumber int, labelList []string, cardinality int, constPromLabels prometheus.Labels) {
	for i := 1; i <= metricsNumber; i++ {

		metrObj := prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:        fmt.Sprintf("%s_metric_%d", metricType, i),
			Help:        fmt.Sprintf("Number of %s_metric_%d", metricType, i),
			ConstLabels: constPromLabels,
		},
			labelList,
		)

		prometheus.MustRegister(metrObj)

		for j := 1; j <= cardinality; j++ {
			var valuesList []string

			for k := 1; k <= len(labelList); k++ {
				valuesList = append(valuesList, fmt.Sprintf("value%d%d%d", i, k, j))
			}

			metrObjWithLab := metrObj.WithLabelValues(valuesList...)
			updateCounter(metrObjWithLab)
		}

	}
}

func updateCounter(metricObj prometheus.Counter) {
	go func() {
		for {
			metricObj.Inc()
			time.Sleep(interval * time.Second)
		}
	}()
}

func gaugeInit(metricType string, metricsNumber int, labelList []string, cardinality int, constPromLabels prometheus.Labels) {
	for i := 1; i <= metricsNumber; i++ {

		metrObj := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:        fmt.Sprintf("%s_metric_%d", metricType, i),
			Help:        fmt.Sprintf("Number of %s_metric_%d", metricType, i),
			ConstLabels: constPromLabels,
		},
			labelList,
		)

		prometheus.MustRegister(metrObj)

		for j := 1; j <= cardinality; j++ {
			var valuesList []string

			for k := 1; k <= len(labelList); k++ {
				valuesList = append(valuesList, fmt.Sprintf("value%d%d%d", i, k, j))
			}

			metrObjWithLab := metrObj.WithLabelValues(valuesList...)
			updateGauge(metrObjWithLab)
		}

	}
}

func updateGauge(metricObj prometheus.Gauge) {
	go func() {
		for {
			metricObj.Set(float64(rand.Intn(100)))
			time.Sleep(interval * time.Second)
		}
	}()
}

func histogramInit(metricType string, metricsNumber int, labelList []string, cardinality int, constPromLabels prometheus.Labels) {
	for i := 1; i <= metricsNumber; i++ {

		metrObj := prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:        fmt.Sprintf("%s_metric_%d", metricType, i),
			Help:        fmt.Sprintf("Number of %s_metric_%d", metricType, i),
			ConstLabels: constPromLabels,
		},
			labelList,
		)

		prometheus.MustRegister(metrObj)

		for j := 1; j <= cardinality; j++ {
			var valuesList []string

			for k := 1; k <= len(labelList); k++ {
				valuesList = append(valuesList, fmt.Sprintf("value%d%d%d", i, k, j))
			}

			metrObjWithLab := metrObj.WithLabelValues(valuesList...)
			updateHistogram(metrObjWithLab)
		}
	}
}

func updateHistogram(metricObj prometheus.Observer) {
	go func() {
		for {
			metricObj.Observe(float64(rand.Intn(100)))
			time.Sleep(interval * time.Second)
		}
	}()
}
