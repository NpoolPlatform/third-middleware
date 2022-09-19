package contact

import (
	"github.com/NpoolPlatform/message/npool/third/mw/v1/contact"
	"google.golang.org/grpc"
)

type Server struct {
	contact.UnimplementedMiddlewareServer
}

func Register(server grpc.ServiceRegistrar) {
	contact.RegisterMiddlewareServer(server, &Server{})
}
