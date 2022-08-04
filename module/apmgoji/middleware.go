package apmgoji

import (
	"fmt"
	// "log"
	"net/http"
	// "os"
	// "path/filepath"

	"github.com/zenazn/goji/web"

	"go.elastic.co/apm/module/apmhttp/v2"
	"go.elastic.co/apm/v2"
	"go.elastic.co/apm/v2/stacktrace"
)

// var Info *log.Logger
// var Debug *log.Logger
// var Error *log.Logger
// var Warn *log.Logger

// var t *apm.Tracer

func init() {
	stacktrace.RegisterLibraryPackage(
		"github.com/zenazn",
	)
	//t, _ = apm.NewTracer("", "")
	// filePath, _ := filepath.Abs("C:\\Users\\Sonika.Prakash\\GitHub\\goji web app\\apmgojiLogs.log")
	// openLogFile, _ := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// Info = log.New(openLogFile, "\tINFO\t", log.Ldate|log.Ltime|log.Lmsgprefix|log.Lshortfile)
	// Debug = log.New(openLogFile, "\tDEBUG\t", log.Ldate|log.Ltime|log.Lmsgprefix|log.Lshortfile)
	// Error = log.New(openLogFile, "\tERROR\t", log.Ldate|log.Ltime|log.Lmsgprefix|log.Lshortfile)
	// Warn = log.New(openLogFile, "\tWARN\t", log.Ldate|log.Ltime|log.Lmsgprefix|log.Lshortfile)
}

// ServeConfig specifies the tracing configuration when using TraceAndServe.
type ServeConfig struct {
	Resource string
}

// type middleware struct {
// 	// engine *web.Mux
// 	tracer *apm.Tracer
// }

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
// func Middleware() func(*web.C, http.Handler) http.Handler {
// 	Info.Println("Inside Middleware function...")
// 	return func(c *web.C, h http.Handler) http.Handler {
// 		Info.Println("Inside 1st return function...")
// 		Debug.Println("c: ", c)
// 		Debug.Println("h: ", h)
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			Info.Println("Inside 2nd return function...")
// 			reqMethod := r.Method
// 			urlPattern := web.GetMatch(*c).RawPattern()
// 			if urlPattern != nil {
// 				Debug.Println(w.Header())
// 				Debug.Println("reqMethod: ", reqMethod)
// 				Debug.Println("urlPattern: ", urlPattern)
// 				Debug.Println("r.Header: ", r.Header)
// 				Debug.Println("r.Response: ", r.Response)
// 				Debug.Println("r.URL: ", r.URL.Path)
// 				requestName := reqMethod + fmt.Sprintf(" %s", urlPattern)
// 				Debug.Println("requestName: ", requestName)
// 				m := &middleware{
// 					tracer: apm.DefaultTracer(),
// 				}
// 				Debug.Println("m: ", m)

// 				rw, ddrw := wrapResponseWriter(w)
// 				h.ServeHTTP(rw, r.WithContext(r.Context()))
// 				tx, body, req := apmhttp.StartTransactionWithBody(m.tracer, requestName, r)
// 				Debug.Println("tx: ", tx)
// 				Debug.Println("body: ", body)
// 				Debug.Println("req: ", req)
// 				Debug.Println("req.Response: ", req.Response)
// 				Debug.Println("req.Body: ", req.Body) // is nil for GET requests
// 				// Debug.Println("req.Response.Status: ", req.Response.Status)
// 				// Debug.Println("req.Response.StatusCode: ", req.Response.StatusCode)
// 				defer tx.End()

// 				defer func() {
// 					httpStatus := getStatus(ddrw.status)
// 					Debug.Println("httpStatus: ", httpStatus)
// 					if v := recover(); v != nil {
// 						Debug.Println("v: ", v)
// 						w.WriteHeader(http.StatusInternalServerError)
// 						e := m.tracer.Recovered(v)
// 						e.SetTransaction(tx)
// 						setContext(&e.Context, req, httpStatus, body)
// 						e.Send()
// 					}
// 					w.WriteHeader(httpStatus)
// 					tx.Result = apmhttp.StatusCodeResult(httpStatus)
// 					if tx.Sampled() {
// 						setContext(&tx.Context, req, httpStatus, body)
// 					}
// 					body.Discard()
// 				}()
// 			} // else {
// 			// 	http.NotFound(w, r)
// 			// }
// 		})
// 	}
// }

// func setContext(ctx *apm.Context, req *http.Request, status int, body *apm.BodyCapturer) {
// 	ctx.SetFramework("goji", "")
// 	ctx.SetHTTPRequest(req)
// 	ctx.SetHTTPRequestBody(body)
// 	ctx.SetHTTPStatusCode(status)
// 	ctx.SetHTTPResponseHeaders(req.Header)
// }

// Middleware returns a goji middleware function that will trace incoming requests.
func Middleware() func(*web.C, http.Handler) http.Handler {
	// Info.Println("\nNew request...")
	m := &middleware{
		tracer: apm.DefaultTracer(),
	}
	return func(c *web.C, h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Debug.Println("context web.C: ", *c)
			resource := r.Method
			p := web.GetMatch(*c).RawPattern()
			// Debug.Println("method: ", resource)
			// Debug.Println("p: ", p)
			if p != nil {
				resource += fmt.Sprintf(" %s", p)
				// Debug.Println("resource: ", resource)
			} else {
				p = r.URL.Path
				resource += fmt.Sprintf(" %s", p)
				// Debug.Println("resource: ", resource)
			}
			m.TraceAndServe(h, w, r, &ServeConfig{
				Resource: resource,
			})
		})
	}
}

type middleware struct {
	// engine *web.Mux
	tracer *apm.Tracer
}

// TraceAndServe serves the handler h using the given ResponseWriter and Request, applying tracing
// according to the specified config.
func (m *middleware) TraceAndServe(h http.Handler, w http.ResponseWriter, r *http.Request, cfg *ServeConfig) {
	// Info.Println("Inside TraceandServe...")
	// Debug.Println("w headers:", w.Header())
	if cfg == nil {
		cfg = new(ServeConfig)
	}
	// Debug.Println("cfg: ", cfg)

	tx, body, req := apmhttp.StartTransactionWithBody(m.tracer, cfg.Resource, r)
	defer tx.End()
	// Debug.Println(tx.TraceContext().State)
	rw, resp := apmhttp.WrapResponseWriter(w)
	// Debug.Println("resp.StatusCode: ", resp.StatusCode)
	// Debug.Println("req.Response: ", req.Response)
	// Debug.Println("req.Header: ", req.Header)
	// Debug.Println("body: ", body)
	defer func() {
		// Debug.Println("r.response: ", r.Response)
		panicked := false
		if v := recover(); v != nil {
			// Info.Println("v is not nil.")
			// Debug.Println("v: ", v)
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
	// Debug.Println("resp.StatusCode now: ", resp.StatusCode)
}

func setContext(ctx *apm.Context, req *http.Request, status int, body *apm.BodyCapturer) {
	ctx.SetFramework("goji", "")
	ctx.SetHTTPRequest(req)
	ctx.SetHTTPRequestBody(body)
	ctx.SetHTTPStatusCode(status)
	ctx.SetHTTPResponseHeaders(req.Header)
}
