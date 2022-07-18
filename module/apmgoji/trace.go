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
	Debug.Println("w headers:", w.Header())
	if cfg == nil {
		cfg = new(ServeConfig)
	}
	Debug.Println("cfg: ", cfg)
	m := &middleware{
		tracer: apm.DefaultTracer(),
	}
	tx, body, req := apmhttp.StartTransactionWithBody(m.tracer, cfg.Resource, r)
	defer tx.End()
	// Debug.Println(tx.TraceContext().State)
	rw, resp := apmhttp.WrapResponseWriter(w)
	Debug.Println("resp.StatusCode: ", resp.StatusCode)
	Debug.Println("req.Response: ", req.Response)
	Debug.Println("req.Header: ", req.Header)
	Debug.Println("body: ", body)
	defer func() {
		Debug.Println("r.response: ", r.Response)
		panicked := false
		if v := recover(); v != nil {
			Info.Println("v is not nil.")
			Debug.Println("v: ", v)
			w.WriteHeader(http.StatusInternalServerError)
			e := m.tracer.Recovered(v)
			e.SetTransaction(tx)
			setContext(&e.Context, req, http.StatusInternalServerError, body)
			e.Send()
			panicked = true
		}
		if panicked {
			resp.StatusCode = http.StatusInternalServerError
			apmhttp.SetTransactionContext(tx, req, resp, body)
			// 	w.WriteHeader(http.StatusInternalServerError)
			// 	tx.Result = apmhttp.StatusCodeResult(http.StatusInternalServerError)
			// 	if tx.Sampled() {
			// 		setContext(&tx.Context, req, http.StatusInternalServerError, body)
			// 	}
		} else {
			apmhttp.SetTransactionContext(tx, req, resp, body)
			// 	w.WriteHeader(httpStatus)
			// 	tx.Result = apmhttp.StatusCodeResult(httpStatus)
			// 	if tx.Sampled() {
			// 		setContext(&tx.Context, req, httpStatus, body)
			// 	}
		}
		body.Discard()
	}()
	h.ServeHTTP(rw, req)
	if resp.StatusCode == 0 {
		resp.StatusCode = http.StatusOK
	}
	Debug.Println("resp.StatusCode now: ", resp.StatusCode)
}

func setContext(ctx *apm.Context, req *http.Request, status int, body *apm.BodyCapturer) {
	ctx.SetFramework("goji", "")
	ctx.SetHTTPRequest(req)
	ctx.SetHTTPRequestBody(body)
	ctx.SetHTTPStatusCode(status)
	ctx.SetHTTPResponseHeaders(req.Header)
}
