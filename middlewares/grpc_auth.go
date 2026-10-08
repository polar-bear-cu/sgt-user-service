package middlewares

import (
	"context"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	userv1 "github.com/polar-bear-cu/sgt-proto/gen/go/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type accessClaims struct {
	jwt.RegisteredClaims
}

type ctxKey string

const userIDKey ctxKey = "user_id"

const edgeHeader = "x-sgt-edge"

// GRPCAuth reads the bearer token if the caller sent one and puts the user id
// in context. Calls with no token pass through only when they did not come
// through envoy (service-to-service traffic); edge calls must carry a token.
func GRPCAuth(jwtSecret string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		fromEdge := len(md.Get(edgeHeader)) > 0
		if fromEdge && info.FullMethod == userv1.UserService_FindOrCreateUser_FullMethodName {
			return nil, status.Error(codes.PermissionDenied, "internal method")
		}
		if len(md.Get("authorization")) == 0 {
			if fromEdge {
				return nil, status.Error(codes.Unauthenticated, "missing token")
			}
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
