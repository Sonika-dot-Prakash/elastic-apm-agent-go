package apmgoji

import "net/http"

type responseWriter struct {
	http.ResponseWriter
	status int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, 0}
}

// Status returns the status code that was monitored.
func (w *responseWriter) Status() int {
	return w.status
}

// wrapResponseWriter wraps an underlying http.ResponseWriter so that it can
// trace the http response codes. It also checks for various http interfaces
// (Flusher, Pusher, CloseNotifier, Hijacker) and if the underlying
// http.ResponseWriter implements them it generates an unnamed struct with the
// appropriate fields.
//
// This code is generated because we have to account for all the permutations
// of the interfaces.
func wrapResponseWriter(w http.ResponseWriter) (http.ResponseWriter, *responseWriter) {
	hFlusher, okFlusher := w.(http.Flusher)
	hPusher, okPusher := w.(http.Pusher)
	hCloseNotifier, okCloseNotifier := w.(http.CloseNotifier)
	hHijacker, okHijacker := w.(http.Hijacker)
	hRoundTripper, okRoundTripper := w.(http.RoundTripper)

	mw := newResponseWriter(w)
	type monitoredResponseWriter interface {
		http.ResponseWriter
		Status() int
	}
	switch {
	case okFlusher && okPusher && okCloseNotifier && okHijacker:
		w = struct {
			monitoredResponseWriter
			http.Flusher
			http.Pusher
			http.CloseNotifier
			http.Hijacker
		}{mw, hFlusher, hPusher, hCloseNotifier, hHijacker}
	case okFlusher && okPusher && okCloseNotifier:
		w = struct {
			monitoredResponseWriter
			http.Flusher
			http.Pusher
			http.CloseNotifier
		}{mw, hFlusher, hPusher, hCloseNotifier}
	case okFlusher && okPusher && okHijacker:
		w = struct {
			monitoredResponseWriter
			http.Flusher
			http.Pusher
			http.Hijacker
		}{mw, hFlusher, hPusher, hHijacker}
	case okFlusher && okCloseNotifier && okHijacker:
		w = struct {
			monitoredResponseWriter
			http.Flusher
			http.CloseNotifier
			http.Hijacker
		}{mw, hFlusher, hCloseNotifier, hHijacker}
	case okPusher && okCloseNotifier && okHijacker:
		w = struct {
			monitoredResponseWriter
			http.Pusher
			http.CloseNotifier
			http.Hijacker
		}{mw, hPusher, hCloseNotifier, hHijacker}
	case okFlusher && okPusher:
		w = struct {
			monitoredResponseWriter
			http.Flusher
			http.Pusher
		}{mw, hFlusher, hPusher}
	case okFlusher && okCloseNotifier:
		w = struct {
			monitoredResponseWriter
			http.Flusher
			http.CloseNotifier
		}{mw, hFlusher, hCloseNotifier}
	case okFlusher && okHijacker:
		w = struct {
			monitoredResponseWriter
			http.Flusher
			http.Hijacker
		}{mw, hFlusher, hHijacker}
	case okPusher && okCloseNotifier:
		w = struct {
			monitoredResponseWriter
			http.Pusher
			http.CloseNotifier
		}{mw, hPusher, hCloseNotifier}
	case okPusher && okHijacker:
		w = struct {
			monitoredResponseWriter
			http.Pusher
			http.Hijacker
		}{mw, hPusher, hHijacker}
	case okCloseNotifier && okHijacker:
		w = struct {
			monitoredResponseWriter
			http.CloseNotifier
			http.Hijacker
		}{mw, hCloseNotifier, hHijacker}
	case okFlusher:
		w = struct {
			monitoredResponseWriter
			http.Flusher
		}{mw, hFlusher}
	case okPusher:
		w = struct {
			monitoredResponseWriter
			http.Pusher
		}{mw, hPusher}
	case okCloseNotifier:
		w = struct {
			monitoredResponseWriter
			http.CloseNotifier
		}{mw, hCloseNotifier}
	case okHijacker:
		w = struct {
			monitoredResponseWriter
			http.Hijacker
		}{mw, hHijacker}
	case okRoundTripper:
		w = struct {
			monitoredResponseWriter
			http.RoundTripper
		}{mw, hRoundTripper}
	default:
		w = mw
	}

	return w, mw
}

func getStatus(status int) int {
	var finalStatus int
	if status == 0 {
		finalStatus = 200
	} else {
		finalStatus = status
	}
	return finalStatus
}
