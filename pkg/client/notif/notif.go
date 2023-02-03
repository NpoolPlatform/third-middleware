//nolint:nolintlint,dupl
package notif

import (
	"context"
	"time"

	usedfor "github.com/NpoolPlatform/message/npool/notif/mgr/v1/notif"

	grpc2 "github.com/NpoolPlatform/go-service-framework/pkg/grpc"

	npool "github.com/NpoolPlatform/message/npool/third/mw/v1/notif"

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

func NotifEmail(
	ctx context.Context,
	appID string,
	fromAccount *string,
	usedFor usedfor.EventType,
	receiverAccount,
	langID string,
	senderName,
	receiverName,
	message *string,
	useTemplate bool,
	title,
	content string,
) error {
	err := do(ctx, func(_ctx context.Context, cli npool.MiddlewareClient) error {
		_, err := cli.NotifEmail(ctx, &npool.NotifEmailRequest{
			AppID:           appID,
			FromAccount:     fromAccount,
			UsedFor:         usedFor,
			ReceiverAccount: receiverAccount,
			LangID:          langID,
			SenderName:      senderName,
			ReceiverName:    receiverName,
			Message:         message,
			UseTemplate:     useTemplate,
			Title:           title,
			Content:         content,
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
