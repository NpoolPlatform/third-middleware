package verify

import (
	"github.com/NpoolPlatform/message/npool/third/mw/v1/verify"
	"google.golang.org/grpc"
)

type Server struct {
	verify.UnimplementedMiddlewareServer
}

func Register(server grpc.ServiceRegistrar) {
	verify.RegisterMiddlewareServer(server, &Server{})
}
