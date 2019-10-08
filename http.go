package goerrors

import (
	"io"
	"io/ioutil"
	"net/http"
	"runtime/debug"

	"github.com/gorilla/handlers"
)

type httpLog struct {
	Method     string
	Query      string
	StatusCode int
}

func httpLogInterceptor(_ io.Writer, params handlers.LogFormatterParams) {
	Log().WithField(HTTPRequestKey, httpLog{
		Method:     params.Request.Method,
		Query:      params.URL.String(),
		StatusCode: params.StatusCode,
	}).Info("HTTP request")
}

func HTTPLoggingHandler(h http.Handler) http.Handler {
	return handlers.CustomLoggingHandler(ioutil.Discard, h, httpLogInterceptor)
}

func HTTPRecoverer(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer OnRecover(func(err error) {
			Log().
				WithContext(CreateContext(r.Context()).AddHTTPRequest(r).AddStack(debug.Stack())).
				WithError(err).
				Error("HTTP recovered")

			w.WriteHeader(http.StatusInternalServerError)
		})

		handler.ServeHTTP(w, r)
	})
}
