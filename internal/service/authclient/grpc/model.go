package grpc

import (
	"math"
	"time"

	"github.com/ardanlabs/service/internal/generated/authpb"
	"github.com/ardanlabs/service/internal/service/authclient"
	"github.com/ardanlabs/service/internal/sys/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func authorizeRequestToGRPC(auth authclient.Authorize) *authpb.AuthorizeRequest {
	cb := authpb.Claims_builder{
		Id:        new(auth.Claims.ID),
		Issuer:    new(auth.Claims.Issuer),
		Subject:   new(auth.Claims.Subject),
		Audience:  auth.Claims.Audience,
		ExpiresAt: new(ndToInt64(auth.Claims.ExpiresAt)),
		NotBefore: new(ndToInt64(auth.Claims.NotBefore)),
		IssuedAt:  new(ndToInt64(auth.Claims.IssuedAt)),
		Roles:     auth.Claims.Roles,
	}
	claims := cb.Build()

	arb := authpb.AuthorizeRequest_builder{
		UserId: new(auth.UserID.String()),
		Rule:   new(auth.Rule),
		Claims: claims,
	}

	req := arb.Build()
	return req
}

func authenticateRespFromGRPC(a *authpb.AuthenticateResponse) (authclient.AuthenticateResp, error) {
	userID := a.GetUserId()
	uID, err := uuid.Parse(userID)
	if err != nil {
		return authclient.AuthenticateResp{}, err
	}

	c := a.GetClaims()
	claims := auth.Claims{
		ID:        c.GetId(),
		Issuer:    c.GetIssuer(),
		Subject:   c.GetSubject(),
		Audience:  c.GetAudience(),
		ExpiresAt: int64ToND(c.GetExpiresAt()),
		NotBefore: int64ToND(c.GetNotBefore()),
		IssuedAt:  int64ToND(c.GetIssuedAt()),
		Roles:     c.GetRoles(),
	}

	resp := authclient.AuthenticateResp{
		UserID: uID,
		Claims: claims,
	}

	return resp, nil
}

// NOTE(AI): optional JWT dates (e.g. nbf) are nil; map nil <-> 0 so encoding
// doesn't panic and decoding stays nil, matching the HTTP client's JSON shape.
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
