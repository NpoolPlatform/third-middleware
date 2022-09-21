//nolint:nolintlint,dupl
package verify

import (
	"context"
	"time"

	"github.com/NpoolPlatform/message/npool/third/mgr/v1/usedfor"

	grpc2 "github.com/NpoolPlatform/go-service-framework/pkg/grpc"

	npool "github.com/NpoolPlatform/message/npool/third/mw/v1/contact"

	constant "github.com/NpoolPlatform/third-middleware/pkg/message/const"
)

var timeout = 10 * time.Second

type handler func(context.Context, npool.MiddlewareClient) error

func do(ctx context.Context, handler handler) error {
	_ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	conn, err := grpc2.GetGRPCConn(constant.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return err
	}

	defer conn.Close()

	cli := npool.NewMiddlewareClient(conn)

	return handler(_ctx, cli)
}

func ContactViaEmail(ctx context.Context,
	appID string,
	usedFor usedfor.UsedFor,
	sender,
	subject,
	body,
	senderName string,
) error {
	err := do(ctx, func(_ctx context.Context, cli npool.MiddlewareClient) error {
		_, err := cli.ContactViaEmail(ctx, &npool.ContactViaEmailRequest{
			AppID:      appID,
			UsedFor:    usedFor,
			Sender:     sender,
			Subject:    subject,
			Body:       body,
			SenderName: senderName,
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
