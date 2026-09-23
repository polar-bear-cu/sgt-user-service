package middlewares

import (
	"context"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ctxKey string

const userIDKey ctxKey = "user_id"

// GRPCAuth reads the bearer token if the caller sent one and puts the user id
// in context. Calls with no token at all (service-to-service traffic) pass
// through untouched - callers that need a real identity check for the id.
func GRPCAuth(jwtSecret string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok || len(md.Get("authorization")) == 0 {
			return handler(ctx, req)
		}

		raw, ok := strings.CutPrefix(md.Get("authorization")[0], "Bearer ")
		if !ok || raw == "" {
			return nil, status.Error(codes.Unauthenticated, "invalid bearer token")
		}

		var claims accessClaims
		_, err := jwt.ParseWithClaims(raw, &claims, func(*jwt.Token) (any, error) {
			return []byte(jwtSecret), nil
		}, jwt.WithValidMethods([]string{"HS256"}))
		if err != nil || claims.Subject == "" {
			return nil, status.Error(codes.Unauthenticated, "invalid or expired token")
		}

		return handler(context.WithValue(ctx, userIDKey, claims.Subject), req)
	}
}

// UserIDFromContext returns the caller's user id and whether a token was
// present at all. ok=false means the call had no token (trusted internal call).
func UserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok
}
