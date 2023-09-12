package vrest

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type RestTraceMaker struct {
	Tracer trace.Tracer
}

type RestTrace struct {
	span trace.Span
}

func NewTraceMaker(tracer trace.Tracer) *RestTraceMaker {
	return &RestTraceMaker{
		Tracer: tracer,
	}
}

func (rtm *RestTraceMaker) NewTrace(req *Request) Trace {
	httpReq := req.Raw
	spanName := fmt.Sprintf("http.request %s %s", req.Method, req.Raw.URL.String())

	attributes := make([]attribute.KeyValue, 0, 4)
	attributes = append(attributes,
		attribute.String("http.method", httpReq.Method),
		attribute.String("http.url", httpReq.URL.String()),
		attribute.String("http.header", fmt.Sprintf("%v", headerWithoutAuth(httpReq.Header))),
	)
	if req.TraceBody {
		attributes = append(attributes, attribute.String("http.body", string(req.BodyBytes)))
	}

	ctx, span := rtm.Tracer.Start(httpReq.Context(), spanName, trace.WithAttributes(attributes...))

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(httpReq.Header))
	return &RestTrace{span: span}
}

func (rt *RestTrace) OnAfterRequest(req *Request) {
	attributes := make([]attribute.KeyValue, 0, 3)
	attributes = append(attributes,
		attribute.Int("http.status_code", req.Response.StatusCode()),
		attribute.String("http.response_header", fmt.Sprintf("%v", headerWithoutAuth(req.Response.Header()))),
	)
	if req.Response.TraceBody {
		attributes = append(attributes, attribute.String("http.response_body", string(req.Response.BodyBytes)))
	}
	rt.span.SetAttributes(attributes...)
}

func (rt *RestTrace) End() {
	rt.span.End()
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
