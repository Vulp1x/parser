// Code generated by goa v3.10.2, DO NOT EDIT.
//
// datasets_service HTTP server
//
// Command:
// $ goa gen github.com/inst-api/parser/design

package server

import (
	"context"
	"net/http"

	datasetsservice "github.com/inst-api/parser/gen/datasets_service"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// Server lists the datasets_service service endpoint HTTP handlers.
type Server struct {
	Mounts             []*MountPoint
	CreateDatasetDraft http.Handler
	UpdateDataset      http.Handler
	FindSimilar        http.Handler
	GetProgress        http.Handler
	ParseDataset       http.Handler
	GetDataset         http.Handler
	GetParsingProgress http.Handler
	ListDatasets       http.Handler
}

// MountPoint holds information about the mounted endpoints.
type MountPoint struct {
	// Method is the name of the service method served by the mounted HTTP handler.
	Method string
	// Verb is the HTTP method used to match requests to the mounted handler.
	Verb string
	// Pattern is the HTTP request path pattern used to match requests to the
	// mounted handler.
	Pattern string
}

// New instantiates HTTP handlers for all the datasets_service service
// endpoints using the provided encoder and decoder. The handlers are mounted
// on the given mux using the HTTP verb and path defined in the design.
// errhandler is called whenever a response fails to be encoded. formatter is
// used to format errors returned by the service methods prior to encoding.
// Both errhandler and formatter are optional and can be nil.
func New(
	e *datasetsservice.Endpoints,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(ctx context.Context, err error) goahttp.Statuser,
) *Server {
	return &Server{
		Mounts: []*MountPoint{
			{"CreateDatasetDraft", "POST", "/api/datasets/draft/"},
			{"UpdateDataset", "PUT", "/api/datasets/{dataset_id}/"},
			{"FindSimilar", "POST", "/api/datasets/{dataset_id}/start/"},
			{"GetProgress", "GET", "/api/datasets/{dataset_id}/progress/"},
			{"ParseDataset", "POST", "/api/datasets/{dataset_id}/parse/"},
			{"GetDataset", "GET", "/api/datasets/{dataset_id}/"},
			{"GetParsingProgress", "GET", "/api/datasets/{dataset_id}/parsing_progress/"},
			{"ListDatasets", "GET", "/api/datasets/"},
		},
		CreateDatasetDraft: NewCreateDatasetDraftHandler(e.CreateDatasetDraft, mux, decoder, encoder, errhandler, formatter),
		UpdateDataset:      NewUpdateDatasetHandler(e.UpdateDataset, mux, decoder, encoder, errhandler, formatter),
		FindSimilar:        NewFindSimilarHandler(e.FindSimilar, mux, decoder, encoder, errhandler, formatter),
		GetProgress:        NewGetProgressHandler(e.GetProgress, mux, decoder, encoder, errhandler, formatter),
		ParseDataset:       NewParseDatasetHandler(e.ParseDataset, mux, decoder, encoder, errhandler, formatter),
		GetDataset:         NewGetDatasetHandler(e.GetDataset, mux, decoder, encoder, errhandler, formatter),
		GetParsingProgress: NewGetParsingProgressHandler(e.GetParsingProgress, mux, decoder, encoder, errhandler, formatter),
		ListDatasets:       NewListDatasetsHandler(e.ListDatasets, mux, decoder, encoder, errhandler, formatter),
	}
}

// Service returns the name of the service served.
func (s *Server) Service() string { return "datasets_service" }

// Use wraps the server handlers with the given middleware.
func (s *Server) Use(m func(http.Handler) http.Handler) {
	s.CreateDatasetDraft = m(s.CreateDatasetDraft)
	s.UpdateDataset = m(s.UpdateDataset)
	s.FindSimilar = m(s.FindSimilar)
	s.GetProgress = m(s.GetProgress)
	s.ParseDataset = m(s.ParseDataset)
	s.GetDataset = m(s.GetDataset)
	s.GetParsingProgress = m(s.GetParsingProgress)
	s.ListDatasets = m(s.ListDatasets)
}

// MethodNames returns the methods served.
func (s *Server) MethodNames() []string { return datasetsservice.MethodNames[:] }

// Mount configures the mux to serve the datasets_service endpoints.
func Mount(mux goahttp.Muxer, h *Server) {
	MountCreateDatasetDraftHandler(mux, h.CreateDatasetDraft)
	MountUpdateDatasetHandler(mux, h.UpdateDataset)
	MountFindSimilarHandler(mux, h.FindSimilar)
	MountGetProgressHandler(mux, h.GetProgress)
	MountParseDatasetHandler(mux, h.ParseDataset)
	MountGetDatasetHandler(mux, h.GetDataset)
	MountGetParsingProgressHandler(mux, h.GetParsingProgress)
	MountListDatasetsHandler(mux, h.ListDatasets)
}

// Mount configures the mux to serve the datasets_service endpoints.
func (s *Server) Mount(mux goahttp.Muxer) {
	Mount(mux, s)
}

// MountCreateDatasetDraftHandler configures the mux to serve the
// "datasets_service" service "create dataset draft" endpoint.
func MountCreateDatasetDraftHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("POST", "/api/datasets/draft/", f)
}

// NewCreateDatasetDraftHandler creates a HTTP handler which loads the HTTP
// request and calls the "datasets_service" service "create dataset draft"
// endpoint.
func NewCreateDatasetDraftHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(ctx context.Context, err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeCreateDatasetDraftRequest(mux, decoder)
		encodeResponse = EncodeCreateDatasetDraftResponse(encoder)
		encodeError    = EncodeCreateDatasetDraftError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "create dataset draft")
		ctx = context.WithValue(ctx, goa.ServiceKey, "datasets_service")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountUpdateDatasetHandler configures the mux to serve the "datasets_service"
// service "update dataset" endpoint.
func MountUpdateDatasetHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("PUT", "/api/datasets/{dataset_id}/", f)
}

// NewUpdateDatasetHandler creates a HTTP handler which loads the HTTP request
// and calls the "datasets_service" service "update dataset" endpoint.
func NewUpdateDatasetHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(ctx context.Context, err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeUpdateDatasetRequest(mux, decoder)
		encodeResponse = EncodeUpdateDatasetResponse(encoder)
		encodeError    = EncodeUpdateDatasetError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "update dataset")
		ctx = context.WithValue(ctx, goa.ServiceKey, "datasets_service")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountFindSimilarHandler configures the mux to serve the "datasets_service"
// service "find similar" endpoint.
func MountFindSimilarHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("POST", "/api/datasets/{dataset_id}/start/", f)
}

// NewFindSimilarHandler creates a HTTP handler which loads the HTTP request
// and calls the "datasets_service" service "find similar" endpoint.
func NewFindSimilarHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(ctx context.Context, err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeFindSimilarRequest(mux, decoder)
		encodeResponse = EncodeFindSimilarResponse(encoder)
		encodeError    = EncodeFindSimilarError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "find similar")
		ctx = context.WithValue(ctx, goa.ServiceKey, "datasets_service")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountGetProgressHandler configures the mux to serve the "datasets_service"
// service "get progress" endpoint.
func MountGetProgressHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/api/datasets/{dataset_id}/progress/", f)
}

// NewGetProgressHandler creates a HTTP handler which loads the HTTP request
// and calls the "datasets_service" service "get progress" endpoint.
func NewGetProgressHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(ctx context.Context, err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeGetProgressRequest(mux, decoder)
		encodeResponse = EncodeGetProgressResponse(encoder)
		encodeError    = EncodeGetProgressError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "get progress")
		ctx = context.WithValue(ctx, goa.ServiceKey, "datasets_service")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountParseDatasetHandler configures the mux to serve the "datasets_service"
// service "parse dataset" endpoint.
func MountParseDatasetHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("POST", "/api/datasets/{dataset_id}/parse/", f)
}

// NewParseDatasetHandler creates a HTTP handler which loads the HTTP request
// and calls the "datasets_service" service "parse dataset" endpoint.
func NewParseDatasetHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(ctx context.Context, err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeParseDatasetRequest(mux, decoder)
		encodeResponse = EncodeParseDatasetResponse(encoder)
		encodeError    = EncodeParseDatasetError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "parse dataset")
		ctx = context.WithValue(ctx, goa.ServiceKey, "datasets_service")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountGetDatasetHandler configures the mux to serve the "datasets_service"
// service "get dataset" endpoint.
func MountGetDatasetHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/api/datasets/{dataset_id}/", f)
}

// NewGetDatasetHandler creates a HTTP handler which loads the HTTP request and
// calls the "datasets_service" service "get dataset" endpoint.
func NewGetDatasetHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(ctx context.Context, err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeGetDatasetRequest(mux, decoder)
		encodeResponse = EncodeGetDatasetResponse(encoder)
		encodeError    = EncodeGetDatasetError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "get dataset")
		ctx = context.WithValue(ctx, goa.ServiceKey, "datasets_service")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountGetParsingProgressHandler configures the mux to serve the
// "datasets_service" service "get parsing progress" endpoint.
func MountGetParsingProgressHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/api/datasets/{dataset_id}/parsing_progress/", f)
}

// NewGetParsingProgressHandler creates a HTTP handler which loads the HTTP
// request and calls the "datasets_service" service "get parsing progress"
// endpoint.
func NewGetParsingProgressHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(ctx context.Context, err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeGetParsingProgressRequest(mux, decoder)
		encodeResponse = EncodeGetParsingProgressResponse(encoder)
		encodeError    = EncodeGetParsingProgressError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "get parsing progress")
		ctx = context.WithValue(ctx, goa.ServiceKey, "datasets_service")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountListDatasetsHandler configures the mux to serve the "datasets_service"
// service "list datasets" endpoint.
func MountListDatasetsHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/api/datasets/", f)
}

// NewListDatasetsHandler creates a HTTP handler which loads the HTTP request
// and calls the "datasets_service" service "list datasets" endpoint.
func NewListDatasetsHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(ctx context.Context, err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeListDatasetsRequest(mux, decoder)
		encodeResponse = EncodeListDatasetsResponse(encoder)
		encodeError    = EncodeListDatasetsError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "list datasets")
		ctx = context.WithValue(ctx, goa.ServiceKey, "datasets_service")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}
