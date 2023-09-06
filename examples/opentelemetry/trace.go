package opentelemetry

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/fond-of-vertigo/vrest"
)

type Tracer struct {
	Tracer trace.Tracer
}

func (o *Tracer) NewTrace(req *vrest.Request) vrest.Trace {
	httpReq := req.Raw
	spanName := fmt.Sprintf("http.request %s %s", req.Method, req.Raw.URL.String())
	ctx, span := o.Tracer.Start(
		httpReq.Context(),
		spanName,
		trace.WithAttributes(
			attribute.String("http.method", httpReq.Method),
			attribute.String("http.url", httpReq.URL.String()),
			attribute.String("http.header", fmt.Sprintf("%v", headerWithoutAuth(httpReq.Header))),
			attribute.String("http.body", string(req.BodyBytes)),
		),
	)

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(httpReq.Header))
	return &Trace{span: span}
}

type Trace struct {
	span trace.Span
}

func (o *Trace) End() {
	o.span.End()
}

func (o *Trace) OnAfterRequest(req *vrest.Request) {
	o.span.SetAttributes(
		attribute.Int("http.status_code", req.Response.StatusCode()),
		attribute.String("http.response_header", fmt.Sprintf("%v", headerWithoutAuth(req.Response.Header()))),
		attribute.String("http.response_body", string(req.Response.BodyBytes)),
	)
}

func headerWithoutAuth(header http.Header) http.Header {
	h := http.Header{}
	for k, v := range header {
		if k != "Authorization" {
			h[k] = v
		}
	}
	return h
}

func main() {
	tracer := otel.GetTracerProvider().Tracer("")
	_ = vrest.New().
		SetTraceMaker(&Tracer{Tracer: tracer})
}
