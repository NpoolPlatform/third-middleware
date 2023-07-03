//nolint:nolintlint,dupl
package send

import (
	"context"
	"fmt"
	"time"

	grpc2 "github.com/NpoolPlatform/go-service-framework/pkg/grpc"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"

	npool "github.com/NpoolPlatform/message/npool/third/mw/v1/send"

	servicename "github.com/NpoolPlatform/third-middleware/pkg/servicename"
)

func do(ctx context.Context, fn func(_ctx context.Context, cli npool.MiddlewareClient) (cruder.Any, error)) (cruder.Any, error) {
	_ctx, cancel := context.WithTimeout(ctx, 10*time.Second) //nolint
	defer cancel()

	conn, err := grpc2.GetGRPCConn(servicename.ServiceDomain, grpc2.GRPCTAG)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	cli := npool.NewMiddlewareClient(conn)

	return fn(_ctx, cli)
}

func SendMessage(ctx context.Context, in *npool.SendMessageRequest) error {
	_, err := do(ctx, func(_ctx context.Context, cli npool.MiddlewareClient) (cruder.Any, error) {
		_, err := cli.SendMessage(ctx, in)
		if err != nil {
			return nil, fmt.Errorf("fail send message: %v", err)
		}
		return nil, nil
	})
	if err != nil {
		return fmt.Errorf("fail send message: %v", err)
	}
	return nil
}
