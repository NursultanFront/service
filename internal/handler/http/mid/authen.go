package mid

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/mail"
	"strings"
	"time"

	errs "github.com/ardanlabs/service/internal/errors"
	"github.com/ardanlabs/service/internal/service/authclient"
	"github.com/ardanlabs/service/internal/sys/auth"
	"github.com/ardanlabs/service/internal/types/role"
	userbus "github.com/ardanlabs/service/internal/user/usecase"
	"github.com/ardanlabs/service/pkg/foundation/web"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Safe client-facing messages for 401 responses. The real cause (token parsing,
// OPA evaluation, database lookups) is kept in errs.Error.Cause and only logged.
//
// NOTE(AI): added by Claude by direct request (migration plan step (h) remainder).
const (
	msgInvalidToken       = "invalid or missing token"
	msgInvalidCredentials = "invalid credentials"
)

// Authenticate is a middleware function that integrates with an authentication client
// to validate user credentials and attach user data to the request context.
func Authenticate(client authclient.Authenticator) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			resp, err := client.Authenticate(ctx, authorizationHeader(r))
			if err != nil {
				return errs.Wrap(errs.Unauthenticated, err, msgInvalidToken)
			}

			ctx = setUserID(ctx, resp.UserID)
			ctx = setClaims(ctx, resp.Claims)

			return next(ctx, r)
		}

		return h
	}

	return m
}

// authorizationHeader returns the bearer token for the incoming request,
// preferring the Authorization header (used by non-browser/service clients)
// and falling back to the httpOnly auth_token cookie (used by the browser
// admin UI, which never has direct JS access to the raw token).
func authorizationHeader(r *http.Request) string {
	if v := r.Header.Get("authorization"); v != "" {
		return v
	}

	c, err := r.Cookie("auth_token")
	if err != nil {
		return ""
	}

	return "Bearer " + c.Value
}

// Bearer processes JWT authentication logic.
func Bearer(ath *auth.Auth) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			ctx, err := HandleAuthentication(ctx, ath, authorizationHeader(r))
			if err != nil {
				return err
			}

			return next(ctx, r)
		}

		return h
	}

	return m
}

func HandleAuthentication(ctx context.Context, ath *auth.Auth, authorizationHeader string) (context.Context, *errs.Error) {
	claims, err := ath.Authenticate(ctx, authorizationHeader)
	if err != nil {
		return ctx, errs.Wrap(errs.Unauthenticated, err, msgInvalidToken)
	}

	if claims.Subject == "" {
		return ctx, errs.Errorf(errs.Unauthenticated, "authorize: you are not authorized for that action, no claims")
	}

	subjectID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return ctx, errs.Wrap(errs.Unauthenticated, fmt.Errorf("parsing subject: %w", err), "invalid token subject") // NOTE(AI): safe message (step (h))
	}

	ctx = setUserID(ctx, subjectID)
	ctx = setClaims(ctx, claims)

	return ctx, nil
}

// Basic processes basic authentication logic.
func Basic(ath *auth.Auth, userBus userbus.ExtBusiness) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			authorizationHeader := r.Header.Get("authorization")
			ctx, err := HandleAuthorization(ctx, authorizationHeader, userBus, ath)
			if err != nil {
				return err
			}

			return next(ctx, r)
		}

		return h
	}

	return m
}

func HandleAuthorization(ctx context.Context, authorizationHeader string, userBus userbus.ExtBusiness, ath *auth.Auth) (context.Context, *errs.Error) {
	email, pass, ok := parseBasicAuth(authorizationHeader)
	if !ok {
		return ctx, errs.Errorf(errs.Unauthenticated, "invalid Basic auth")
	}

	addr, err := mail.ParseAddress(email)
	if err != nil {
		return ctx, errs.Wrap(errs.Unauthenticated, err, msgInvalidCredentials)
	}

	usr, err := userBus.Authenticate(ctx, *addr, pass)
	if err != nil {
		return ctx, errs.Wrap(errs.Unauthenticated, err, msgInvalidCredentials)
	}

	claims := auth.Claims{
		Subject:   usr.ID.String(),
		Issuer:    ath.Issuer(),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(8760 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		Roles:     role.ParseToString(usr.Roles),
	}

	subjectID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return ctx, errs.Wrap(errs.Unauthenticated, fmt.Errorf("parsing subject: %w", err), "invalid token subject") // NOTE(AI): safe message (step (h))
	}

	ctx = setUserID(ctx, subjectID)
	ctx = setClaims(ctx, claims)

	return ctx, nil
}

func parseBasicAuth(auth string) (string, string, bool) {
	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Basic" {
		return "", "", false
	}

	c, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", false
	}

	username, password, ok := strings.Cut(string(c), ":")
	if !ok {
		return "", "", false
	}

	return username, password, true
}
