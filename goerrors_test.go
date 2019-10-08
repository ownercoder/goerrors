package goerrors

import (
	"net/http"
	"net/http/httptest"
	"runtime/debug"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func init() {
	opts := &Options{
		// kernel_stage DSN is here for testing purposes only
		// use SENTRY_DSN env var in prod
		SentryDSN: "https://9292cf26950e46d285cc9f0e0e371a24:0538baffa31b425b9e5a7e21335963dd@sentry.butik.ru/11",

		// must be false in prod
		SentrySyncMode: true,

		// play with it to remove useless frames from top of the call stack
		SentrySkipFrames: 1,
	}

	// use one of these formatters
	formatter := &TextFormatter{
		TextFormatter: logrus.TextFormatter{ForceColors: true, FullTimestamp: true},
	}
	// formatter := &JSONFormatter{}

	// and don't forget to set correct logging level
	level := DebugLevel

	if err := InitLog(level, formatter, opts); err != nil {
		panic(err)
	}
}

func testHTTPServer(beforeServe func(router *mux.Router), handler func(writer http.ResponseWriter, request *http.Request)) {
	req, _ := http.NewRequest(http.MethodGet, "/test?this=query#plz", nil)
	w := httptest.NewRecorder()

	r := mux.NewRouter()
	r.Use(HTTPLoggingHandler)
	r.Use(HTTPRecoverer)
	r.Path("/test").Methods("GET").HandlerFunc(handler)

	beforeServe(r)
	r.ServeHTTP(w, req)
}

func ExampleLogging() {
	testHTTPServer(
		func(router *mux.Router) {
			//
		},
		func(writer http.ResponseWriter, request *http.Request) {
			// you should add context into logging entries when possible
			Log().WithContext(request.Context()).Trace("test trace")
			Log().WithContext(request.Context()).Debug("test debug")
			Log().WithContext(request.Context()).Info("test info")
			Log().WithContext(request.Context()).Warn("test warn")
			Log().WithContext(request.Context()).Error("test error")

			// but also it works without context
			Log().Info("test info without context")
		},
	)
	// Output:
}

func ExampleLoggingWithExtra() {
	testHTTPServer(
		func(router *mux.Router) {
			//
		},
		func(writer http.ResponseWriter, request *http.Request) {
			// this is an only wrapper for easy calling AddStack, AddHTTPRequest, ... funcs
			// nothing more
			ctx := CreateContext(request.Context())

			ctx.AddStack(debug.Stack())
			ctx.AddHTTPRequest(request)
			ctx.AddFingerprint([]string{"finger", "print"})

			Log().WithContext(ctx).Error("test error")
		},
	)
	// Output:
}

func ExamplePanicRecovery() {
	testHTTPServer(
		func(router *mux.Router) {
			//
		},
		func(writer http.ResponseWriter, request *http.Request) {
			panic("Aaa!")

			// recoverer automatically adds error, stack and request info into logging entry,
			// handles Sentry sending
			// and returns 500 as HTTP response
		},
	)
	// Output:
}
