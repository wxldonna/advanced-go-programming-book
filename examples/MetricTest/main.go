package main

import (
	"context"
	"flag"
	"fmt"

	"time"

	"go.opentelemetry.io/contrib/instrumentation/runtime"

	"go.opentelemetry.io/otel/metric/instrument"

	"go.opentelemetry.io/otel/attribute"

	"go.opentelemetry.io/contrib/instrumentation/host"

	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.wdf.sap.corp/velocity/trc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

var tracer = trc.InitTraceTopic("Initlization", "executable package")

var trcFlag = flag.String("trc", "debug", "e.g. -trc=debug,main:warning")

// logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		tracer.Debugf("Request URL %s", r.RequestURI)
		// get the context from request
		ctx := r.Context()
		// add the context
		ctx = context.WithValue(ctx, "version", "1.0.0")

		r = r.WithContext(ctx)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func init() {
	trc.Application = "1Server"

}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) //
	title := vars["title"]
	tracer.Infof("get products from parameter title %s", string(title))

	ctx := r.Context()

	tracer.Debugf("version is %s", ctx.Value("version"))
	fmt.Fprintf(w, "version is %s", ctx.Value("version"))
	SynchMetricx(ctx)
}

func main() {
	flag.Parse()
	log.Printf("trcFlag is %v", *trcFlag)
	SynchMetricx(context.Background())
	AsynchMetricx()
	log.Printf("trcFlag is %v", *trcFlag)
	if err := trc.ReconfigFromString(*trcFlag); err != nil {
		log.Fatal(err)
	}
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.HandleFunc("/products/{title}", GetProducts)
	RegisterHandlers(r, NewResource("newService"))
	log.Fatal(http.ListenAndServe(":8181", r))
}

// NewPrometheusMetricsProvider sets up the metrics handler for Prometheus. Create an http.Handler
// to expose this on an endpoint (`mux.Handle("/metrics", NewPrometheusMetricsProvider())`).
func NewPrometheusMetricsProvider(service *resource.Resource) (*prometheus.Exporter, error) {
	//nolint:exhaustruct
	config := prometheus.Config{}

	ctrl := controller.New(
		processor.NewFactory(
			selector.NewWithHistogramDistribution(
				histogram.WithExplicitBoundaries(config.DefaultHistogramBoundaries),
			),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
		controller.WithResource(service),
	)

	if err := host.Start(); err != nil {
		return nil, fmt.Errorf("failed to start host instrumentation: %w", err)
	}

	exp, err := prometheus.New(config, ctrl)
	if err != nil {
		return nil, err
	}

	global.SetMeterProvider(ctrl)

	if err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second)); err != nil {
		return nil, err
	}

	return exp, nil
}

// RegisterHandlers registers the metrics handler under the `/metrics` path.
func RegisterHandlers(router *mux.Router, service *resource.Resource) error {
	metricsHandler, err := NewPrometheusMetricsProvider(service)
	if err != nil {
		return err
	}

	return router.Handle("/metrics", metricsHandler).GetError()
}

// NewResource sets up the resource for the service.
func NewResource(name string) *resource.Resource {
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(name),
			semconv.DeploymentEnvironmentKey.String("production"),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	return res
}

func AsynchMetricx() {
	provider := global.MeterProvider()
	meter := provider.Meter("connections")
	counter, err := meter.AsyncInt64().UpDownCounter("sap_axino_abap_connections")
	if err != nil {
		// ...
	}

	err = meter.RegisterCallback(
		[]instrument.Asynchronous{counter},
		func(ctx context.Context) {
			// collect current value

			counter.Observe(ctx, 100)
		},
	)
	if err != nil {
		// ...
	}

}

func SynchMetricx(ctx context.Context) {
	provider := global.MeterProvider()
	meter := provider.Meter("counter1")
	counter, err := meter.SyncInt64().Counter("test1")
	if err != nil {
		panic(err)
	}
	counter.Add(ctx, 1, attribute.String("key1", "value1"))
}
