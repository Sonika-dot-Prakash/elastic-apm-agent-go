package apmgoji

import (
	"net/http"

	"go.elastic.co/apm/module/apmhttp/v2"
	"go.elastic.co/apm/v2"
)

// ServeConfig specifies the tracing configuration when using TraceAndServe.
type ServeConfig struct {
	Resource string
}

type middleware struct {
	// engine *web.Mux
	tracer *apm.Tracer
}

// TraceAndServe serves the handler h using the given ResponseWriter and Request, applying tracing
// according to the specified config.
func TraceAndServe(h http.Handler, w http.ResponseWriter, r *http.Request, cfg *ServeConfig) {
	Info.Println("Inside TraceandServe...")
	if cfg == nil {
		cfg = new(ServeConfig)
	}
	Debug.Println("cfg: ", cfg)
	m := &middleware{
		tracer: apm.DefaultTracer(),
	}
	tx, body, req := apmhttp.StartTransactionWithBody(m.tracer, cfg.Resource, r)
	defer tx.End()
	rw, ddrw := wrapResponseWriter(w)
	Debug.Println("ddrw.status: ", ddrw.status)
	Debug.Println("req.Response: ", req.Response)
	Debug.Println("req.Header: ", req.Header)
	Debug.Println("body: ", body)
	defer func() {
		httpStatus := getStatus(ddrw.status)
		Debug.Println("httpStatus: ", httpStatus)
		if v := recover(); v != nil {
			Info.Println("v is not nil.")
			Debug.Println("v: ", v)
			w.WriteHeader(http.StatusInternalServerError)
			e := m.tracer.Recovered(v)
			e.SetTransaction(tx)
			setContext(&e.Context, req, httpStatus, body)
			e.Send()
		}
		w.WriteHeader(httpStatus)
		tx.Result = apmhttp.StatusCodeResult(httpStatus)
		if tx.Sampled() {
			setContext(&tx.Context, req, httpStatus, body)
		}
		body.Discard()
	}()
	h.ServeHTTP(rw, r.WithContext(r.Context())) // this may cause panic if context is nil
}

func setContext(ctx *apm.Context, req *http.Request, status int, body *apm.BodyCapturer) {
	ctx.SetFramework("goji", "")
	ctx.SetHTTPRequest(req)
	ctx.SetHTTPRequestBody(body)
	ctx.SetHTTPStatusCode(status)
	ctx.SetHTTPResponseHeaders(req.Header)
}
