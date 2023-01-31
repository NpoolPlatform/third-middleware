package notif

import (
	"github.com/NpoolPlatform/message/npool/third/mw/v1/template/notif"

	"google.golang.org/grpc"
)

type Server struct {
	notif.UnimplementedMiddlewareServer
}

func Register(server grpc.ServiceRegistrar) {
	notif.RegisterMiddlewareServer(server, &Server{})
}
