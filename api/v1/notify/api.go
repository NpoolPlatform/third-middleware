package notify

import (
	"github.com/NpoolPlatform/message/npool/third/mw/v1/notify"
	"google.golang.org/grpc"
)

type Server struct {
	notify.UnimplementedMiddlewareServer
}

func Register(server grpc.ServiceRegistrar) {
	notify.RegisterMiddlewareServer(server, &Server{})
}
