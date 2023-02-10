//nolint:nolintlint,dupl
package notif

import (
	"context"
	"fmt"
	"time"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"

	grpc2 "github.com/NpoolPlatform/go-service-framework/pkg/grpc"

	npool "github.com/NpoolPlatform/message/npool/third/mw/v1/notif"

	constant "github.com/NpoolPlatform/third-middleware/pkg/message/const"
)

var timeout = 10 * time.Second

type handler func(context.Context, npool.MiddlewareClient) (cruder.Any, error)

func do(ctx context.Context, handler handler) (cruder.Any, error) {
	_ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	conn, err := grpc2.GetGRPCConn(constant.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	cli := npool.NewMiddlewareClient(conn)

	return handler(_ctx, cli)
}

func SendNotifEmail(ctx context.Context, subject, content, from, to string) error {
	_, err := do(ctx, func(_ctx context.Context, cli npool.MiddlewareClient) (cruder.Any, error) {
		resp, err := cli.SendNotifEmail(ctx, &npool.SendNotifEmailRequest{
			Subject: subject,
			Content: content,
			From:    from,
			To:      to,
		})
		if err != nil {
			return nil, fmt.Errorf("fail get notif: %v", err)
		}
		return resp, nil
	})
	if err != nil {
		return fmt.Errorf("fail send email: %v", err)
	}
	return nil
}
