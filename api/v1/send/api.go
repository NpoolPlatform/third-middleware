package send

import (
	"github.com/NpoolPlatform/message/npool/third/mw/v1/send"
	"google.golang.org/grpc"
)

type Server struct {
	send.UnimplementedMiddlewareServer
}

func Register(server grpc.ServiceRegistrar) {
	send.RegisterMiddlewareServer(server, &Server{})
}
