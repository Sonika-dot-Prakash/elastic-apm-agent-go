package apmgoji

import (
	"fmt"
	"net/http"

	"github.com/zenazn/goji/web"

	"go.elastic.co/apm/module/apmhttp/v2"
	"go.elastic.co/apm/v2"
	"go.elastic.co/apm/v2/stacktrace"
)

type middleware struct {
	// engine *web.Mux
	tracer *apm.Tracer
}

// Middleware returns a new Goji middleware handler for tracing
// requests and reporting errors.
// func Middleware(engine *web.Mux) web.HandlerFunc {
// 	m := &middleware{
// 		engine: engine,
// 	}
// 	return m.handle
// }

// func (m *middleware) handle(c *web.C, w http.ResponseWriter, r *http.Request) {}

// Middleware returns a new Goji middleware handler for tracing requests
func Middleware() func(*web.C, http.Handler) http.Handler {
	return func(c *web.C, h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqMethod := r.Method
			urlPattern := web.GetMatch(*c).RawPattern()
			if urlPattern != nil {
				requestName := reqMethod + fmt.Sprintf(" %s", urlPattern)
				m := &middleware{
					tracer: apm.DefaultTracer(),
				}
				tx, body, req := apmhttp.StartTransactionWithBody(m.tracer, requestName, r)
				defer tx.End()

				defer func() {
					if v := recover(); v != nil {
						w.WriteHeader(http.StatusInternalServerError)
						e := m.tracer.Recovered(v)
						e.SetTransaction(tx)
						setContext(&e.Context, req, body)
						e.Send()
					}
					w.WriteHeader(req.Response.StatusCode)
					tx.Result = apmhttp.StatusCodeResult(req.Response.StatusCode)
					if tx.Sampled() {
						setContext(&tx.Context, req, body)
					}
					body.Discard()
				}()
			}
		})
	}
}

func setContext(ctx *apm.Context, req *http.Request, body *apm.BodyCapturer) {
	ctx.SetFramework("goji", "")
	ctx.SetHTTPRequest(req)
	ctx.SetHTTPRequestBody(body)
	ctx.SetHTTPStatusCode(req.Response.StatusCode)
	ctx.SetHTTPResponseHeaders(req.Header)
}

func init() {
	stacktrace.RegisterLibraryPackage(
		"github.com/zenazn",
	)
}
