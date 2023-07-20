package api

import (
	"context"

	v1 "github.com/NpoolPlatform/message/npool/third/mw/v1"

	"github.com/NpoolPlatform/third-middleware/api/v1/oauth"
	"github.com/NpoolPlatform/third-middleware/api/v1/send"
	"github.com/NpoolPlatform/third-middleware/api/v1/verify"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	v1.UnimplementedMiddlewareServer
}

func Register(server grpc.ServiceRegistrar) {
	v1.RegisterMiddlewareServer(server, &Server{})
	verify.Register(server)
	send.Register(server)
	oauth.Register(server)
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := v1.RegisterMiddlewareHandlerFromEndpoint(context.Background(), mux, endpoint, opts); err != nil {
		return err
	}
	return nil
}
