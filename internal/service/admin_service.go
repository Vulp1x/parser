package service

//
// import (
// 	"context"
// 	"fmt"
//
// 	"github.com/inst-api/parser/internal/domain"
// 	"github.com/inst-api/parser/internal/store/bots"
// 	"github.com/inst-api/parser/pkg/api"
// 	"goa.design/goa/v3/security"
// )
//
// type botsStore interface {
// 	SaveBots(ctx context.Context, bots domain.Bots) (int64, error)
// }
//
// // admin_service service example implementation.
// // The example methods log the requests and return zero values.
// type adminServicesrvc struct {
// 	api.UnimplementedParserServer
// 	botsStore botsStore
// }
//
// // NewAdminService returns the admin_service service implementation.
// func NewAdminService(botsStore *bots.Store) api.ParserServer {
// 	return &adminServicesrvc{botsStore: botsStore}
// }
//
// // JWTAuth implements the authorization logic for service "admin_service" for
// // the "jwt" security scheme.
// func (s *adminServicesrvc) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
// 	//
// 	// TBD: add authorization logic.
// 	//
// 	// In case of authorization failure this function should return
// 	// one of the generated error structs, e.g.:
// 	//
// 	//    return ctx, myservice.MakeUnauthorizedError("invalid token")
// 	//
// 	// Alternatively this function may return an instance of
// 	// goa.ServiceError with a Name field value that matches one of
// 	// the design error names, e.g:
// 	//
// 	//    return ctx, goa.PermanentError("unauthorized", "invalid token")
// 	//
// 	return ctx, fmt.Errorf("not implemented")
// }
//
// // создать драфт задачи
// func (s *adminServicesrvc) SaveBots(ctx context.Context, request *api.SaveBotsRequest) (*api.SaveBotsResponse, error) {
// 	count, err := s.botsStore.SaveBots(ctx, request.Bots)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &api.SaveBotsResponse{
// 		BotsSaved: int32(count),
// 	}, nil
// }
