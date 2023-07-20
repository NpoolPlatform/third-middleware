package oauth

import (
	"github.com/NpoolPlatform/message/npool/third/mw/v1/oauth"
	"google.golang.org/grpc"
)

type Server struct {
	oauth.UnimplementedMiddlewareServer
}

func Register(server grpc.ServiceRegistrar) {
	oauth.RegisterMiddlewareServer(server, &Server{})
}
