package opentracing

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/opentracing/basictracer-go"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.wdf.sap.corp/velocity/trc"
)

const (
	TagUserAgent     = ext.StringTagName("http.user_agent")
	TagResponseSize  = ext.Uint32TagName("http.response_size")
	TagResponseType  = ext.StringTagName("http.response_type")
	TagHijacked      = ext.BoolTagName("http.hijacked")
	TagVsystemTenant = ext.StringTagName("vsystem.tenant")
	TagVsystemUser   = ext.StringTagName("vsystem.user")
)

type OpenTracingAPI1 struct {
	Tracer trc.Tracer
}

var OpenTracer = OpenTracingAPI1{}

func Const(tracer trc.Tracer) {
	OpenTracer.Tracer = tracer
	log.Printf("Const is called ")
}

// initial the open tracing
func Initopentracing(name string) {
	if OpenTracer.Tracer.IsInfo() {
		OpenTracer.Tracer.Infof("initializing opentracing with VelocityTracer")

		opentracing.SetGlobalTracer(NewVelocityTracer(name))
	} else {
		OpenTracer.Tracer.Infof("initializing opentracing with NoopTracer")

		opentracing.SetGlobalTracer(opentracing.NoopTracer{})
	}
}

// create recorder

const (
	trcInfoOperation    = "operation"
	trcInfoDuration     = "duration"
	trcInfoStartedUTC   = "started-at"
	trcInfoSpanID       = "span-id"
	trcInfoParentSpanID = "parent-span-id"
	trcInfoTagPrefix    = "tag:"
)

func NewVelocityTracer(name string) opentracing.Tracer {
	opts := basictracer.DefaultOptions()
	opts.Recorder = NewVelocityRecorder(name)

	return basictracer.NewWithOptions(opts)
}

type VelocityRecorder struct {
	name   string
	tracer trc.Tracer
}

func NewVelocityRecorder(name string) *VelocityRecorder {
	return &VelocityRecorder{
		name:   name,
		tracer: OpenTracer.Tracer,
	}
}

func (r *VelocityRecorder) RecordSpan(span basictracer.RawSpan) {
	infos := getTracerInfos(span)
	tracer := r.tracer.Sub(infos...)

	warn := false

	if code, ok := span.Tags[string(ext.HTTPStatusCode)]; ok {
		if v, ok := code.(uint16); ok && v >= 300 {
			warn = true
		}
	}

	if warn {
		tracer.Warningf("Telemetry")
	} else {
		tracer.Infof("Telemetry")
	}
}

func getTracerInfos(span basictracer.RawSpan) []trc.Info {
	infos := make([]trc.Info, 0)
	infos = append(infos,
		trc.NewInfo(trcInfoOperation, span.Operation),
		trc.NewInfo(trcInfoSpanID, fmt.Sprint(span.Context.SpanID)),
		trc.NewInfo(trcInfoTagPrefix+trcInfoDuration, span.Duration.String()),
		trc.NewInfo(trcInfoTagPrefix+trcInfoStartedUTC, span.Start.UTC().String()),
	)

	if span.ParentSpanID != 0 {
		infos = append(infos, trc.NewInfo(trcInfoParentSpanID, fmt.Sprint(span.ParentSpanID)))
	}

	for k, v := range span.Tags {
		infos = append(infos, trc.NewInfo(trcInfoTagPrefix+k, fmt.Sprint(v)))
	}

	for k, v := range span.Context.Baggage {
		infos = append(infos, trc.NewInfo(k, v))
	}

	return infos
}

// create span in the middleware
func Middleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		wireContext, _ := opentracing.GlobalTracer().Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(req.Header),
		)

		operation := fmt.Sprintf("%s %s", req.Method, req.URL)

		// if wireContext == nil, a root span will be created.
		options := []opentracing.StartSpanOption{
			ext.RPCServerOption(wireContext),
			opentracing.Tag{Key: string(ext.Component), Value: "component1"},
			opentracing.Tag{Key: string(ext.HTTPUrl), Value: req.URL.String()},
			opentracing.Tag{Key: string(ext.HTTPMethod), Value: req.Method},
			opentracing.Tag{Key: string(ext.PeerHostIPv4), Value: req.RemoteAddr},
			opentracing.Tag{Key: string(TagVsystemTenant), Value: "default"},
			opentracing.Tag{Key: string(TagVsystemUser), Value: "default"},
		}
		OpenTracer.Tracer.Infof("span middleware is called ")
		span, ctx := opentracing.StartSpanFromContext(req.Context(), operation, options...)
		defer span.Finish()

		next.ServeHTTP(w, req.WithContext(ctx))

	})

}

//attach span to outbround request
func AttachToRequest(ctx context.Context, req *http.Request) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		_ = opentracing.GlobalTracer().Inject(
			span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(req.Header),
		)
	}
}
