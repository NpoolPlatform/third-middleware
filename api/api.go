package api

import (
	"context"

	"github.com/NpoolPlatform/third-middleware/api/v1/notif"
	templateemail "github.com/NpoolPlatform/third-middleware/api/v1/template/email"
	templatefrontend "github.com/NpoolPlatform/third-middleware/api/v1/template/frontend"
	templatesms "github.com/NpoolPlatform/third-middleware/api/v1/template/sms"

	"github.com/NpoolPlatform/third-middleware/api/v1/contact"

	v1 "github.com/NpoolPlatform/message/npool/third/mw/v1"
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
	templatefrontend.Register(server)
	templateemail.Register(server)
	templatesms.Register(server)
	notif.Register(server)
	contact.Register(server)
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := v1.RegisterMiddlewareHandlerFromEndpoint(context.Background(), mux, endpoint, opts); err != nil {
		return err
	}
	return nil
}
