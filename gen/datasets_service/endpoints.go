// Code generated by goa v3.8.5, DO NOT EDIT.
//
// datasets_service endpoints
//
// Command:
// $ goa gen github.com/inst-api/parser/design

package datasetsservice

import (
	"context"

	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// Endpoints wraps the "datasets_service" service endpoints.
type Endpoints struct {
	CreateDatasetDraft goa.Endpoint
	UpdateDataset      goa.Endpoint
	FindSimilar        goa.Endpoint
	ParseDataset       goa.Endpoint
	GetDataset         goa.Endpoint
	GetProgress        goa.Endpoint
	ListDatasets       goa.Endpoint
}

// NewEndpoints wraps the methods of the "datasets_service" service with
// endpoints.
func NewEndpoints(s Service) *Endpoints {
	// Casting service to Auther interface
	a := s.(Auther)
	return &Endpoints{
		CreateDatasetDraft: NewCreateDatasetDraftEndpoint(s, a.JWTAuth),
		UpdateDataset:      NewUpdateDatasetEndpoint(s, a.JWTAuth),
		FindSimilar:        NewFindSimilarEndpoint(s, a.JWTAuth),
		ParseDataset:       NewParseDatasetEndpoint(s, a.JWTAuth),
		GetDataset:         NewGetDatasetEndpoint(s, a.JWTAuth),
		GetProgress:        NewGetProgressEndpoint(s, a.JWTAuth),
		ListDatasets:       NewListDatasetsEndpoint(s, a.JWTAuth),
	}
}

// Use applies the given middleware to all the "datasets_service" service
// endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.CreateDatasetDraft = m(e.CreateDatasetDraft)
	e.UpdateDataset = m(e.UpdateDataset)
	e.FindSimilar = m(e.FindSimilar)
	e.ParseDataset = m(e.ParseDataset)
	e.GetDataset = m(e.GetDataset)
	e.GetProgress = m(e.GetProgress)
	e.ListDatasets = m(e.ListDatasets)
}

// NewCreateDatasetDraftEndpoint returns an endpoint function that calls the
// method "create dataset draft" of service "datasets_service".
func NewCreateDatasetDraftEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*CreateDatasetDraftPayload)
		var err error
		sc := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"driver", "admin"},
			RequiredScopes: []string{},
		}
		ctx, err = authJWTFn(ctx, p.Token, &sc)
		if err != nil {
			return nil, err
		}
		return s.CreateDatasetDraft(ctx, p)
	}
}

// NewUpdateDatasetEndpoint returns an endpoint function that calls the method
// "update dataset" of service "datasets_service".
func NewUpdateDatasetEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*UpdateDatasetPayload)
		var err error
		sc := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"driver", "admin"},
			RequiredScopes: []string{},
		}
		ctx, err = authJWTFn(ctx, p.Token, &sc)
		if err != nil {
			return nil, err
		}
		return s.UpdateDataset(ctx, p)
	}
}

// NewFindSimilarEndpoint returns an endpoint function that calls the method
// "find similar" of service "datasets_service".
func NewFindSimilarEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*FindSimilarPayload)
		var err error
		sc := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"driver", "admin"},
			RequiredScopes: []string{},
		}
		ctx, err = authJWTFn(ctx, p.Token, &sc)
		if err != nil {
			return nil, err
		}
		return s.FindSimilar(ctx, p)
	}
}

// NewParseDatasetEndpoint returns an endpoint function that calls the method
// "parse dataset" of service "datasets_service".
func NewParseDatasetEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*ParseDatasetPayload)
		var err error
		sc := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"driver", "admin"},
			RequiredScopes: []string{},
		}
		ctx, err = authJWTFn(ctx, p.Token, &sc)
		if err != nil {
			return nil, err
		}
		return s.ParseDataset(ctx, p)
	}
}

// NewGetDatasetEndpoint returns an endpoint function that calls the method
// "get dataset" of service "datasets_service".
func NewGetDatasetEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*GetDatasetPayload)
		var err error
		sc := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"driver", "admin"},
			RequiredScopes: []string{},
		}
		ctx, err = authJWTFn(ctx, p.Token, &sc)
		if err != nil {
			return nil, err
		}
		return s.GetDataset(ctx, p)
	}
}

// NewGetProgressEndpoint returns an endpoint function that calls the method
// "get progress" of service "datasets_service".
func NewGetProgressEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*GetProgressPayload)
		var err error
		sc := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"driver", "admin"},
			RequiredScopes: []string{},
		}
		ctx, err = authJWTFn(ctx, p.Token, &sc)
		if err != nil {
			return nil, err
		}
		return s.GetProgress(ctx, p)
	}
}

// NewListDatasetsEndpoint returns an endpoint function that calls the method
// "list datasets" of service "datasets_service".
func NewListDatasetsEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*ListDatasetsPayload)
		var err error
		sc := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"driver", "admin"},
			RequiredScopes: []string{},
		}
		ctx, err = authJWTFn(ctx, p.Token, &sc)
		if err != nil {
			return nil, err
		}
		return s.ListDatasets(ctx, p)
	}
}
