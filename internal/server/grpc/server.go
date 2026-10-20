// Package grpc starts and stops the gRPC server for the auth service.
package grpc

import (
	"context"
	"net"

	"github.com/ardanlabs/service/internal/generated/authpb"
	grpcauthapp "github.com/ardanlabs/service/internal/handler/grpc/auth"
	"github.com/ardanlabs/service/pkg/foundation/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server owns the lifecycle of the gRPC server.
//
// NOTE(AI): moved out of internal/handler/grpc/auth by Claude by direct
// request (DOMAIN_LAYOUT_PLAN.md step 2); same startup/shutdown behavior.
type Server struct {
	log *logger.Logger
	gs  *grpc.Server
}

// Start registers the auth service on a new gRPC server and starts serving
// on the listener in the background.
func Start(ctx context.Context, log *logger.Logger, lis net.Listener, api *grpcauthapp.App) *Server {
	log.Info(context.Background(), "auth service", "status", "start auth server")

	gs := grpc.NewServer(
		grpc.UnaryInterceptor(api.AuthInterceptor),
	)

	authpb.RegisterAuthServer(gs, api)
	reflection.Register(gs)

	go func() {
		if err := gs.Serve(lis); err != nil {
			log.Error(ctx, "startup", "status", "auth server error", "err", err)
		}
	}()

	return &Server{log: log, gs: gs}
}

// Shutdown stops the gRPC server gracefully.
func (s *Server) Shutdown() {
	s.log.Info(context.Background(), "shutdown", "status", "gRPC server stopped")
	s.gs.GracefulStop()
}
