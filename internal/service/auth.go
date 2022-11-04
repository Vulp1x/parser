package service

import (
	"context"
	"fmt"
	"reflect"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	datasetsservice "github.com/inst-api/parser/gen/datasets_service"
	"github.com/inst-api/parser/internal/sessions"
	"github.com/inst-api/parser/pkg/logger"
	"goa.design/goa/v3/security"
)

type contextKey string

const userIDContextKey = contextKey("UserID")

// auth_service service example implementation.
// The example methods log the requests and return zero values.
type authService struct {
	securityCfg sessions.Configuration
}

// JWTAuth implements the authorization logic for service "auth_service" for
// the "jwt" security scheme.
func (s *authService) JWTAuth(ctx context.Context, tokenString string, scheme *security.JWTScheme) (context.Context, error) {
	if tokenString == "" {
		return ctx, datasetsservice.BadRequest("No token provided")
	}

	token, err := jwt.ParseWithClaims(tokenString, &sessions.SessionClaims{}, s.securityCfg.KeyFunc)
	if err != nil {
		logger.Infof(ctx, "Failed to parse token: %v", err)

		return ctx, datasetsservice.Unauthorized("failed to parse token")
	}

	// logger.Debugf(ctx, "got sec scheme %+v", scheme)

	if claims, ok := token.Claims.(*sessions.SessionClaims); ok && token.Valid {
		ctx = logger.WithFields(ctx, logger.Fields{"user_id": claims.UserID})
		ctx = AddUserIDToContext(ctx, claims.UserID)
		logger.Debugf(ctx, "Successfully checked token")
		return ctx, nil
	}

	logger.Infof(ctx, "Token is valid: %t or Claims are SessionClaims: %#v", token.Valid, reflect.TypeOf(token.Claims))

	return ctx, fmt.Errorf("internal err")
}

// AddUserIDToContext добавляет user_id в контекст
func AddUserIDToContext(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(userIDContextKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("no user id in ctx: %#v", ctx)
	}

	return userID, nil
}
