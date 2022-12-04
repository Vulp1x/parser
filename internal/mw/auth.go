package mw

import (
	"context"
	"net/http"
	"reflect"
	"runtime/debug"

	"github.com/dgrijalva/jwt-go"
	"github.com/inst-api/parser/internal/sessions"
	"github.com/inst-api/parser/pkg/logger"
)

type contextKey string

func (c contextKey) String() string {
	return "middlewares context value " + string(c)
}

func (c contextKey) Write(r *http.Request, val interface{}) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), c, val))
}

const (
	// claimsRequestKey key for using in context.
	claimsRequestKey contextKey = "Claims"
	// tokenRequestKey key for using in context.
	tokenRequestKey contextKey = "Token"
)

// CheckSession check sesssion middleware.
func CheckSession(securityConfig sessions.Configuration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Bearer")
			if tokenString == "" {
				Error(w, r, http.StatusUnauthorized, "No token provided")

				return
			}

			token, err := jwt.ParseWithClaims(tokenString, &sessions.SessionClaims{}, securityConfig.KeyFunc)
			if err != nil {
				InternalError(w, r, "Failed to parser token: %v", err)

				return
			}

			if claims, ok := token.Claims.(*sessions.SessionClaims); ok && token.Valid {
				LogEntrySetField(r, "user_id", claims.UserID)
				logger.Debugf(r.Context(), "Successfully checked token")

				next.ServeHTTP(w, claimsRequestKey.Write(r, claims))
			} else {
				InternalError(w, r, "Token is valid: %t or Claims are SessionClaims: %v", token.Valid, reflect.TypeOf(token.Claims))

				return
			}
		}

		return http.HandlerFunc(fn)
	}
}

// GetToken is used to get User JSON Web Token from request context.
func GetToken(r *http.Request) string {
	signedToken, ok := r.Context().Value(tokenRequestKey).(string)
	if !ok {
		panic("can`t get signed token from context when expected" +
			string(debug.Stack()))
	}

	return signedToken
}

// GetClaims is used to get claims(userID etc) from request context.
func GetClaims(r *http.Request) *sessions.SessionClaims {
	claims, ok := r.Context().Value(claimsRequestKey).(*sessions.SessionClaims)
	if !ok {
		panic("can`t get claims from context when expected" +
			string(debug.Stack()))
	}

	return claims
}
