//nolint:nolintlint,dupl
package notify

import (
	"context"
	"time"

	"github.com/NpoolPlatform/message/npool/third/mgr/v1/usedfor"

	grpc2 "github.com/NpoolPlatform/go-service-framework/pkg/grpc"

	npool "github.com/NpoolPlatform/message/npool/third/mw/v1/notify"

	constant "github.com/NpoolPlatform/third-manager/pkg/message/const"
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

func NotifyEmail(ctx context.Context,
	appID,
	fromAccount string,
	usedFor usedfor.UsedFor,
	receiverAccount,
	langID,
	senderName,
	receiverName string,
) error {
	err := do(ctx, func(_ctx context.Context, cli npool.MiddlewareClient) error {
		_, err := cli.NotifyEmail(ctx, &npool.NotifyEmailRequest{
			AppID:           appID,
			FromAccount:     fromAccount,
			UsedFor:         usedFor,
			ReceiverAccount: receiverAccount,
			LangID:          langID,
			SenderName:      senderName,
			ReceiverName:    receiverName,
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
