// Package auth maintains the web based api for auth access.
package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	errs "github.com/ardanlabs/service/internal/errors"
	"github.com/ardanlabs/service/internal/handler/http/mid"
	"github.com/ardanlabs/service/internal/service/authclient"
	"github.com/ardanlabs/service/internal/sys/auth"
	"github.com/ardanlabs/service/pkg/foundation/web"
)

type app struct {
	auth *auth.Auth
}

func newApp(ath *auth.Auth) *app {
	return &app{
		auth: ath,
	}
}

func (a *app) token(ctx context.Context, r *http.Request) web.Encoder {
	kid := web.Param(r, "kid")
	if kid == "" {
		return errs.NewFieldErrors("kid", errors.New("missing kid"))
	}

	// The BearerBasic middleware function generates the claims.
	claims := mid.GetClaims(ctx)

	tkn, err := a.auth.GenerateToken(kid, claims)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	// Also set the token as an httpOnly cookie so browser clients never
	// need to touch the raw JWT in JS-accessible storage (localStorage),
	// which closes off token theft via XSS. Non-browser clients (curl,
	// admin CLI, service-to-service) keep using the JSON body above.
	if w := web.GetWriter(ctx); w != nil {
		cookie := http.Cookie{
			Name:     "auth_token",
			Value:    tkn,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		}
		if claims.ExpiresAt != nil {
			cookie.Expires = claims.ExpiresAt.Time
		}
		http.SetCookie(w, &cookie)
	}

	return token{Token: tkn}
}

// logout clears the httpOnly auth cookie set by token. JS cannot delete an
// httpOnly cookie itself, so the browser has to ask the server to do it -
// setting MaxAge to a negative value tells the browser to expire it now.
func (a *app) logout(ctx context.Context, r *http.Request) web.Encoder {
	if w := web.GetWriter(ctx); w != nil {
		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   -1,
		})
	}

	return nil
}

func (a *app) authenticate(ctx context.Context, r *http.Request) web.Encoder {
	// The middleware is actually handling the authentication. So if the code
	// gets to this handler, authentication passed.

	userID, err := mid.GetUserID(ctx)
	if err != nil {
		return errs.New(errs.Unauthenticated, err)
	}

	resp := authclient.AuthenticateResp{
		UserID: userID,
		Claims: mid.GetClaims(ctx),
	}

	return resp
}

func (a *app) authorize(ctx context.Context, r *http.Request) web.Encoder {
	var auth authclient.Authorize
	if err := web.Decode(r, &auth); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	if err := a.auth.Authorize(ctx, auth.Claims, auth.UserID, auth.Rule); err != nil {
		// NOTE(AI): 403 with a safe message; claims and rule stay in the logged cause (step (h)).
		return errs.Wrap(errs.PermissionDenied, fmt.Errorf("authorize: claims[%v] rule[%v]: %w", auth.Claims.Roles, auth.Rule, err), "you are not authorized for that action")
	}

	return nil
}
