package log

import (
	"bytes"
	"github.com/gin-contrib/logger"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
	"io"
	"net/http"
)

var (
	C zerolog.Context
)

func Middleware() gin.HandlerFunc {
	return logger.SetLogger(
		logger.WithLogger(func(c *gin.Context, l zerolog.Logger) zerolog.Logger {
			C = l.Level(zerolog.DebugLevel).With().Caller().
				Str("id", requestid.Get(c)).
				Str("path", c.Request.URL.Path).
				Interface("hdr", c.Request.Header)
			lc := C
			if trace.SpanFromContext(c.Request.Context()).SpanContext().IsValid() {
				lc = lc.
					Str("trace_id", trace.SpanFromContext(c.Request.Context()).SpanContext().TraceID().String()).
					Str("span_id", trace.SpanFromContext(c.Request.Context()).SpanContext().SpanID().String())
			}
			if c.Request.Method != http.MethodGet {
				bodyByt, _ := io.ReadAll(c.Request.Body)
				c.Request.Body = io.NopCloser(bytes.NewReader(bodyByt))
				if binding.Default(c.Request.Method, c.ContentType()) == binding.JSON {
					lc = lc.RawJSON("body", bodyByt)
				} else {
					lc = lc.Str("body", string(bodyByt))
				}
			}
			return lc.Logger()
		}),
	)
}
