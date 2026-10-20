// Package auth maintains the gRPC service for authentication.
package auth

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/ardanlabs/service/internal/generated/authpb"
	"github.com/ardanlabs/service/internal/handler/http/mid"
	"github.com/ardanlabs/service/internal/sys/auth"
	userbus "github.com/ardanlabs/service/internal/user/usecase"
	"github.com/ardanlabs/service/pkg/foundation/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Config holds the dependencies for the gRPC auth service.
type Config struct {
	Log     *logger.Logger
	Auth    *auth.Auth
	UserBus userbus.ExtBusiness
}

// App represents the gRPC service for authentication. Starting and stopping
// the gRPC server lives in internal/server/grpc.
type App struct {
	authpb.UnimplementedAuthServer
	log     *logger.Logger
	auth    *auth.Auth
	userBus userbus.ExtBusiness
}

// NewApp constructs the gRPC auth service.
//
// NOTE(AI): exported (was newApp) by Claude by direct request, so the server
// lifecycle can move to internal/server/grpc (DOMAIN_LAYOUT_PLAN.md step 2).
func NewApp(cfg Config) *App {
	return &App{
		log:     cfg.Log,
		auth:    cfg.Auth,
		userBus: cfg.UserBus,
	}
}

// Token generates a token for a given key ID.
func (a *App) Token(ctx context.Context, req *authpb.TokenRequest) (*authpb.TokenResponse, error) {
	kid := req.GetKid()
	if kid == "" {
		return nil, status.Error(codes.InvalidArgument, "kid is required")
	}

	// NOTE(AI): the AuthInterceptor validated the Basic auth credentials and
	// stored the user's real claims (subject, roles, expiry), same as HTTP.
	claims := mid.GetClaims(ctx)
	if claims.Subject == "" {
		return nil, status.Error(codes.Unauthenticated, "no claims")
	}

	token, err := a.auth.GenerateToken(kid, claims)
	if err != nil {
		a.log.Error(ctx, "token", "err", err)
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	trb := authpb.TokenResponse_builder{
		Token: new(token),
	}

	return trb.Build(), nil
}

// Authenticate validates a bearer token.
func (a *App) Authenticate(ctx context.Context, req *authpb.AuthenticateRequest) (*authpb.AuthenticateResponse, error) {
	reqToken := req.GetToken()
	if reqToken == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	a.log.Info(ctx, "authenticate", "method", "grpc")

	claims, err := a.auth.Authenticate(ctx, reqToken)
	if err != nil {
		a.log.Error(ctx, "authenticate", "err", err)
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	// NOTE(AI): the client reads the full Claims message (roles included) and
	// forwards it to Authorize, so it must be populated here.
	cb := authpb.Claims_builder{
		Id:        new(claims.ID),
		Issuer:    new(claims.Issuer),
		Subject:   new(claims.Subject),
		Audience:  claims.Audience,
		ExpiresAt: new(ndToInt64(claims.ExpiresAt)),
		NotBefore: new(ndToInt64(claims.NotBefore)),
		IssuedAt:  new(ndToInt64(claims.IssuedAt)),
		Roles:     claims.Roles,
	}

	arb := authpb.AuthenticateResponse_builder{
		UserId:  new(claims.Subject),
		Subject: new(claims.Subject),
		Roles:   claims.Roles,
		Claims:  cb.Build(),
	}

	return arb.Build(), nil
}

// Authorize checks if a user is authorized for a specific rule.
func (a *App) Authorize(ctx context.Context, req *authpb.AuthorizeRequest) (*authpb.AuthorizeResponse, error) {
	reqUserID := req.GetUserId()
	if reqUserID == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	reqRule := req.GetRule()
	if reqRule == "" {
		return nil, status.Error(codes.InvalidArgument, "rule is required")
	}

	reqClaims := req.GetClaims()

	// Convert gRPC request to auth client format
	claims := auth.Claims{
		ID:        reqClaims.GetId(),
		Issuer:    reqClaims.GetIssuer(),
		Subject:   reqClaims.GetSubject(),
		Audience:  reqClaims.GetAudience(),
		ExpiresAt: int64ToND(reqClaims.GetExpiresAt()),
		NotBefore: int64ToND(reqClaims.GetNotBefore()),
		IssuedAt:  int64ToND(reqClaims.GetIssuedAt()),
		Roles:     reqClaims.GetRoles(),
	}

	userID, err := uuid.Parse(reqUserID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
	}

	err = a.auth.Authorize(ctx, claims, userID, reqRule)
	if err != nil {
		// NOTE(AI): a policy denial is wrapped in ErrInvalidAuthzOPA
		// (ErrForbidden is never returned), so it was reported as Internal.
		if errors.Is(err, auth.ErrInvalidAuthzOPA) {
			a.log.Info(ctx, "authorize", "status", "denied", "err", err)
			return nil, status.Error(codes.PermissionDenied, "not authorized")
		}

		a.log.Error(ctx, "authorize", "err", err)
		return nil, status.Error(codes.Internal, "authorization failed")
	}

	return &authpb.AuthorizeResponse{}, nil
}

// =============================================================================

func (a *App) AuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// NOTE(AI): match the generated full method names; the proto package is
	// grpcauthapp, so the old "/auth.Auth/..." literals never matched and the
	// interceptor was a no-op. Authenticate is intentionally not intercepted:
	// the token arrives in the request body and the handler validates it.
	switch info.FullMethod {
	case authpb.Auth_Token_FullMethodName:
		return a.authorize(ctx, req, info, handler)

	default:
		return handler(ctx, req)
	}
}

// NOTE(AI): errors are returned as gRPC status errors so clients see
// codes.Unauthenticated instead of codes.Unknown.
func (a *App) authorize(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no metadata")
	}

	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return nil, status.Error(codes.Unauthenticated, "no authorization header")
	}

	ctx, err := mid.HandleAuthorization(ctx, authHeaders[0], a.userBus, a.auth)
	if err != nil {
		a.log.Info(ctx, "token", "status", "unauthenticated", "err", err)
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	return handler(ctx, req)
}

// NOTE(AI): optional JWT dates (e.g. nbf, which no token minter sets) are nil;
// map nil <-> 0 so encoding doesn't panic and decoding stays nil.
func ndToInt64(nd *jwt.NumericDate) int64 {
	if nd == nil {
		return 0
	}
	return nd.Unix()
}

func int64ToND(in int64) *jwt.NumericDate {
	if in == 0 {
		return nil
	}
	round, frac := math.Modf(float64(in))
	return jwt.NewNumericDate(time.Unix(int64(round), int64(frac*1e9)))
}
