package email

import (
	"github.com/NpoolPlatform/message/npool/third/mw/v1/template/email"

	"google.golang.org/grpc"
)

type Server struct {
	email.UnimplementedMiddlewareServer
}

func Register(server grpc.ServiceRegistrar) {
	email.RegisterMiddlewareServer(server, &Server{})
}
