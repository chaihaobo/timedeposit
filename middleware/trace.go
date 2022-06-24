// Package middleware
// @author： Boice
// @createTime：2022/6/23 20:02
package middleware

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/bns-engineering/common/telemetry"
	ins "gitlab.com/bns-engineering/common/telemetry/instrumentation"
	otelcontrib "go.opentelemetry.io/contrib"
	"go.opentelemetry.io/otel/label"

	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/propagators"
	"go.opentelemetry.io/otel/semconv"
)

var propagator = otel.NewCompositeTextMapPropagator(propagators.TraceContext{},
	b3.B3{InjectEncoding: b3.B3MultipleHeader | b3.B3SingleHeader})
var optionshttp = []otelhttptrace.Option{otelhttptrace.WithPropagators(propagator)}

func TraceMiddleware(api *telemetry.API) gin.HandlerFunc {
	return func(c *gin.Context) {
		tracer := api.Tracer(trace.WithInstrumentationVersion(otelcontrib.SemVersion()))
		_, _, spanCtx := otelhttptrace.Extract(c, c.Request, optionshttp...)
		ctx, span := tracer.Start(trace.ContextWithRemoteSpanContext(c, spanCtx), c.Request.URL.Path)
		c.Request = c.Request.WithContext(ctx)
		defer span.End()
		c.Next()

		status := c.Writer.Status()
		attrs := semconv.HTTPAttributesFromHTTPStatusCode(status)
		spanStatus, spanMessage := semconv.SpanStatusFromHTTPStatusCode(status)
		span.SetAttributes(attrs...)
		span.SetStatus(spanStatus, spanMessage)
		if len(c.Errors) > 0 {
			span.SetAttributes(label.String("gin.errors", c.Errors.String()))
		}

		span.SetAttributes(ins.Method.String(c.Request.Method+" "+c.Request.URL.Path),
			ins.CorrelationID.String(span.SpanContext().TraceID.String()+"-"+span.SpanContext().SpanID.String()),
			semconv.HTTPClientIPKey.String(c.Request.Header.Get(ins.XForwardedFor)),
			semconv.HTTPUserAgentKey.String(c.Request.Header.Get(ins.UserAgent)))
	}
}
