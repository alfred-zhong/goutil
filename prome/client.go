package prome

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Client represents the client for prometheus server to pull data from.
type Client struct {
	ServiceName string
	Path        string

	// Enable metrics of runtime. Default enabled.
	EnableRuntime bool
	// Interval of update memstats. Unit: seconds. Default 60s.
	UpdateMemStatsInterval int

	// Labels which will always be attached to metrics.
	ConstLabels prometheus.Labels

	srv               *http.Server
	cancel            context.CancelFunc
	handler           http.Handler
	runtimeCollectors []prometheus.Collector
	collectors        []prometheus.Collector
}

// ListenAndServe listen on the addr and provide access for prometheus server to
// pull data.
func (c *Client) ListenAndServe(addr string) error {
	if c.handler == nil {
		reg := prometheus.NewRegistry()

		// Register collectors.
		if c.EnableRuntime {
			registerRuntime(c.ServiceName, &c.runtimeCollectors, c.ConstLabels)
			reg.MustRegister(c.runtimeCollectors...)
			constructs = append(constructs, updateRuntimeGuage)

			ctx, cancel := context.WithCancel(context.Background())
			c.cancel = cancel
			go updateRuntimeMemstats(ctx, c.UpdateMemStatsInterval)
		}
		reg.MustRegister(c.collectors...)

		c.handler = decorateHandler(
			promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
		)
	}

	http.Handle(c.Path, c.handler)
	c.srv = &http.Server{
		Addr: addr,
	}
	if err := c.srv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Close shutdown of listening.
func (c *Client) Close() error {
	if c.cancel != nil {
		c.cancel()
	}
	if c.srv != nil {
		return c.srv.Shutdown(context.Background())
	}
	return nil
}

// Handler returns the http handler which can be used for fetch metrics data.
func (c *Client) Handler() http.Handler {
	if c.handler == nil {
		reg := prometheus.NewRegistry()

		// Register collectors.
		if c.EnableRuntime {
			registerRuntime(c.ServiceName, &c.runtimeCollectors, c.ConstLabels)
			reg.MustRegister(c.runtimeCollectors...)
			constructs = append(constructs, updateRuntimeGuage)

			ctx, cancel := context.WithCancel(context.Background())
			c.cancel = cancel
			go updateRuntimeMemstats(ctx, c.UpdateMemStatsInterval)
		}
		reg.MustRegister(c.collectors...)

		c.handler = decorateHandler(
			promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
		)
	}
	return c.handler
}

var constructs []func()

type decoratedHandler struct {
	h http.Handler
}

func (d *decoratedHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	for _, f := range constructs {
		f()
	}
	d.h.ServeHTTP(rw, r)
}

func decorateHandler(h http.Handler) http.Handler {
	return &decoratedHandler{h}
}

// NewClient creates and returns a new Client instance.
func NewClient(serviceName string, path string) *Client {
	return &Client{
		ServiceName:   serviceName,
		Path:          path,
		EnableRuntime: true,
	}
}

// Option defines argument for create Server.
type Option func(*Client)

// WithRuntimeEnable creates Option which determines whether to enables metrics
// of runtime. updateMemstatsInterval make sense only if enable set to true.
func WithRuntimeEnable(enable bool, updateMemstatsInterval int) Option {
	return func(c *Client) {
		c.EnableRuntime = enable
		if enable && updateMemstatsInterval > 0 {
			c.UpdateMemStatsInterval = updateMemstatsInterval
		}
	}
}

// WithConstLables creates Option which append const labels to the client.
func WithConstLables(kv ...string) Option {
	return func(c *Client) {
		m := make(map[string]string, len(kv)/2)
		for i := 0; i < len(kv); i += 2 {
			if i+1 < len(kv) {
				m[kv[i]] = kv[i+1]
			}
		}
		c.ConstLabels = prometheus.Labels(m)
	}
}

// NewClientWithOption creates Client with optional Option.
func NewClientWithOption(serviceName, path string, opts ...Option) *Client {
	c := &Client{
		ServiceName:   serviceName,
		Path:          path,
		EnableRuntime: true,
	}

	for _, opt := range opts {
		opt(c)
	}
	return c
}
