package prome_test

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/alfred-zhong/goutil/prome"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func ExampleClient_basic1() {
	client := prome.NewClientWithOption("test", "/foo")
	if err := client.ListenAndServe(":9000"); err != nil {
		panic(err)
	}
}

func ExampleClient_basic2() {
	client := prome.NewClientWithOption(
		"test", "/foo",
		prome.WithRuntimeEnable(true, 5),
	)
	if err := client.ListenAndServe(":9000"); err != nil {
		panic(err)
	}
}

func ExampleClient_withoutruntime() {
	client := prome.NewClientWithOption(
		"test", "/foo",
		prome.WithRuntimeEnable(false, 0),
	)
	if err := client.ListenAndServe(":9000"); err != nil {
		panic(err)
	}
}

func ExampleClient_labels() {
	client := prome.NewClientWithOption(
		"test", "/foo",
		prome.WithConstLables(
			"env", "test",
			"foo", "bar",
		),
	)
	if err := client.ListenAndServe(":9000"); err != nil {
		panic(err)
	}
}

func ExampleClient_gin() {
	client := prome.NewClientWithOption(
		"test", "",
		prome.WithRuntimeEnable(false, 0),
	)

	e := gin.New()
	e.Use(client.MiddlewareRequestCount(""))
	e.Use(client.MiddlewareRequestDuration("", nil))

	e.GET("/foo", gin.WrapH(client.Handler()))
	e.GET("/hello/:name", func(c *gin.Context) {
		c.String(200, "indeed, %s", c.Param("name"))
	})
	e.GET("/sleep", func(c *gin.Context) {
		time.Sleep(time.Duration(rand.Int63n(1000)) * time.Millisecond)
	})
	if err := e.Run(":9527"); err != nil {
		panic(err)
	}
}

func ExampleClient_custommetrics() {
	client := prome.NewClientWithOption(
		"test", "/foo",
		prome.WithRuntimeEnable(false, 0),
	)
	client.ConstLabels = prometheus.Labels{
		"env": "test",
	}

	counter := client.AddCounter(prometheus.CounterOpts{
		Name: "test_counter",
	})
	counterVec := client.AddCounterVec(prometheus.CounterOpts{
		Name: "test_counterVec",
	}, []string{"name"})

	gauge := client.AddGauge(prometheus.GaugeOpts{
		Name: "test_gauge",
	})
	gaugeVec := client.AddGaugeVec(prometheus.GaugeOpts{
		Name: "test_gaugeVec",
	}, []string{"name"})

	histogram := client.AddHistogram(prometheus.HistogramOpts{
		Name:    "test_histogram",
		Buckets: []float64{10, 20, 30, 40, 50},
	})
	histogramVec := client.AddHistogramVec(prometheus.HistogramOpts{
		Name:    "test_histogramVec",
		Buckets: []float64{10, 20, 30, 40, 50},
	}, []string{"name"})

	summary := client.AddSummary(prometheus.SummaryOpts{
		Name:       "test_summary",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	})
	summaryVec := client.AddSummaryVec(prometheus.SummaryOpts{
		Name:       "test_summaryVec",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, []string{"name"})

	go func() {
		for {
			counter.Inc()
			counterVec.WithLabelValues("hysteria").Inc()

			gauge.Inc()
			gaugeVec.WithLabelValues("roxanne").Inc()

			histogram.Observe(rand.Float64() * 100)
			histogramVec.WithLabelValues("cassandra").Observe(rand.Float64() * 100)

			summary.Observe(rand.Float64() * 100)
			summaryVec.WithLabelValues("hysteria").Observe(rand.Float64() * 30)
			summaryVec.WithLabelValues("roxanne").Observe(rand.Float64() * 50)
			summaryVec.WithLabelValues("riful").Observe(rand.Float64() * 100)

			time.Sleep(time.Second)
		}
	}()

	if err := client.ListenAndServe(":9000"); err != nil {
		panic(err)
	}
}

func ExampleClient_shutdown() {
	client := prome.NewClientWithOption("test", "/foo")
	go func() {
		if err := client.ListenAndServe(":9000"); err != nil {
			panic(err)
		}
		fmt.Println("server shutdown")
	}()

	time.Sleep(5 * time.Second)
	if err := client.Close(); err != nil {
		panic(err)
	}
	time.Sleep(10 * time.Second)
}
