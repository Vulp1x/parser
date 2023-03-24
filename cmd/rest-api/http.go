package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	swagger "github.com/go-openapi/runtime/middleware"
	datasetsservice "github.com/inst-api/parser/gen/datasets_service"
	datasetsservicesvr "github.com/inst-api/parser/gen/http/datasets_service/server"
	"github.com/inst-api/parser/internal/mw"
	"github.com/inst-api/parser/internal/service/multipart"
	"github.com/inst-api/parser/pkg/logger"
	"github.com/riandyrn/otelchi"
	goahttp "goa.design/goa/v3/http"
	httpmdlwr "goa.design/goa/v3/http/middleware"
	"goa.design/goa/v3/middleware"
)

// handleHTTPServer starts configures and starts a HTTP server on the given
// URL. It shuts down the server if any error is received in the error channel.
func handleHTTPServer(
	ctx context.Context,
	host, port string,
	datasetsServiceEndpoints *datasetsservice.Endpoints,
	wg *sync.WaitGroup,
	errc chan error,
	debug bool,
) {
	// Provide the transport specific request decoder and response encoder.
	// The goa http package has built-in support for JSON, XML and gob.
	// Other encodings can be used by providing the corresponding functions,
	// see goa.design/implement/encoding.
	var (
		dec = goahttp.RequestDecoder
		enc = goahttp.ResponseEncoder
	)

	// Build the service HTTP request multiplexer and configure it to serve
	// HTTP requests to the service endpoints.
	var mux goahttp.Muxer
	{
		mux = goahttp.NewMuxer()
	}

	opts := swagger.SwaggerUIOpts{SpecURL: "openapi3.yaml"}
	mux.Handle("GET", "/openapi3.yaml", http.FileServer(http.Dir("./gen/http")).ServeHTTP)
	mux.Handle("GET", "/openapi.yaml", http.FileServer(http.Dir("./gen/http")).ServeHTTP)

	mux.Handle("GET", "/docs", swagger.SwaggerUI(opts, nil).ServeHTTP)

	// Wrap the endpoints with the transport specific layers. The generated
	// server packages contains code generated from the design which maps
	// the service input and output data structures to HTTP requests and
	// responses.
	var (
		datasetsServiceServer *datasetsservicesvr.Server
	)
	{
		eh := errorHandler()
		datasetsServiceServer = datasetsservicesvr.New(datasetsServiceEndpoints, mux, dec, enc, eh, nil, multipart.DatasetsServiceUploadFilesDecoderFunc)

		datasetsServiceServer.Use(mw.RequestLoggerWithDebug(mux, debug))
		datasetsServiceServer.Use(httpmdlwr.RequestID())

		if debug {
			// datasetsServiceServer.Use(httpmdlwr.RequestLoggerWithDebug(mux, os.Stdout))
			// tasksServiceServer.Use(httpmdlwr.RequestLoggerWithDebug(mux, os.Stdout))
			// locationsServiceServer.Use(httpmdlwr.RequestLoggerWithDebug(mux, os.Stdout))
			// adminServiceServer.Use(httpmdlwr.RequestLoggerWithDebug(mux, os.Stdout))

			// datasetsServiceServer.Use(mw.RequestLoggerWithDebug(mux, true))
			// routesServiceServer.Use(mw.RequestLoggerWithDebug(mux))
			// locationsServiceServer.Use(mw.RequestLoggerWithDebug(mux))
			// adminServiceServer.Use(mw.RequestLoggerWithDebug(mux))

		}
	}
	// Configure the mux.
	datasetsservicesvr.Mount(mux, datasetsServiceServer)

	router := chi.NewRouter()
	// router.Use(cors.Handler(cors.Options{
	// 	// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
	// 	AllowedOrigins: []string{"https://*", "http://*"},
	// 	AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
	// 	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowedHeaders:   []string{"Accept", "datasetsorization", "Content-Type", "X-CSRF-Token"},
	// 	ExposedHeaders:   []string{"Link"},
	// 	AllowCredentials: false,
	// 	MaxAge:           300, // Maximum value not ignored by any of major browsers
	// 	Debug:            true,
	// }))

	router.Use(cors.AllowAll().Handler)
	router.Use(otelchi.Middleware(
		"chi-parser",
		// otelchi.WithRequestMethodInSpanName(true),
		otelchi.WithChiRoutes(router),
	))

	router.Mount("/", mux)

	srv := &http.Server{Addr: fmt.Sprintf("%s:%s", host, port), Handler: router}
	for _, m := range datasetsServiceServer.Mounts {
		logger.Infof(ctx, "HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}

	defer wg.Done()

	// Start HTTP server in a separate goroutine.
	go func() {
		logger.Infof(ctx, "HTTP server listening on %s:%s", host, port)
		// errc <- srv.ListenAndServeTLS("deploy/cert.pem", "deploy/key.pem")
		errc <- srv.ListenAndServe()
	}()

	<-ctx.Done()
	logger.Infof(ctx, "shutting down HTTP server at %s", host)

	// Shutdown gracefully with a 10s timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = srv.Shutdown(ctx)
}

// errorHandler returns a function that writes and logs the given error.
// The function also writes and logs the error unique ID so that it's possible
// to correlate.
func errorHandler() func(context.Context, http.ResponseWriter, error) {
	return func(ctx context.Context, w http.ResponseWriter, err error) {
		id := ctx.Value(middleware.RequestIDKey).(string)
		_, _ = w.Write([]byte(fmt.Sprintf("[%s] encoding: %v", id, err)))
		logger.Infof(ctx, "[%s] ERROR: %s", id, err.Error())
	}
}
