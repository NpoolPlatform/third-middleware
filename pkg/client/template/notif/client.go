//nolint:nolintlint,dupl
package notif

import (
	"context"
	"fmt"
	"time"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	mgrpb "github.com/NpoolPlatform/message/npool/third/mgr/v1/template/notif"

	grpc2 "github.com/NpoolPlatform/go-service-framework/pkg/grpc"

	npool "github.com/NpoolPlatform/message/npool/third/mw/v1/template/notif"

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

func GetNotifTemplate(ctx context.Context, id string) (*mgrpb.NotifTemplate, error) {
	info, err := do(ctx, func(_ctx context.Context, cli npool.MiddlewareClient) (cruder.Any, error) {
		resp, err := cli.GetNotifTemplate(ctx, &npool.GetNotifTemplateRequest{
			ID: id,
		})
		if err != nil {
			return nil, fmt.Errorf("fail get notif: %v", err)
		}
		return resp.Info, nil
	})
	if err != nil {
		return nil, fmt.Errorf("fail get notif: %v", err)
	}
	return info.(*mgrpb.NotifTemplate), nil
}

func GetNotifTemplateOnly(ctx context.Context, conds *mgrpb.Conds) (*mgrpb.NotifTemplate, error) {
	info, err := do(ctx, func(_ctx context.Context, cli npool.MiddlewareClient) (cruder.Any, error) {
		resp, err := cli.GetNotifTemplateOnly(ctx, &npool.GetNotifTemplateOnlyRequest{
			Conds: conds,
		})
		if err != nil {
			return nil, fmt.Errorf("fail get notif: %v", err)
		}
		return resp.Info, nil
	})
	if err != nil {
		return nil, fmt.Errorf("fail get notif: %v", err)
	}
	return info.(*mgrpb.NotifTemplate), nil
}

func GetNotifTemplates(ctx context.Context, conds *mgrpb.Conds, offset, limit uint32) ([]*mgrpb.NotifTemplate, uint32, error) {
	var total uint32
	infos, err := do(ctx, func(_ctx context.Context, cli npool.MiddlewareClient) (cruder.Any, error) {
		resp, err := cli.GetNotifTemplates(ctx, &npool.GetNotifTemplatesRequest{
			Conds:  conds,
			Offset: offset,
			Limit:  limit,
		})
		if err != nil {
			return nil, fmt.Errorf("fail get notif: %v", err)
		}
		total = resp.GetTotal()
		return resp.Infos, nil
	})
	if err != nil {
		return nil, 0, fmt.Errorf("fail get notif: %v", err)
	}
	return infos.([]*mgrpb.NotifTemplate), total, nil
}
