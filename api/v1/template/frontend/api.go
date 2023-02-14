package frontend

import (
	"github.com/NpoolPlatform/message/npool/third/mw/v1/template/frontend"

	"google.golang.org/grpc"
)

type Server struct {
	frontend.UnimplementedMiddlewareServer
}

func Register(server grpc.ServiceRegistrar) {
	frontend.RegisterMiddlewareServer(server, &Server{})
}
