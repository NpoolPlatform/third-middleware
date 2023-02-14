//nolint:nolintlint,dupl
package frontend

import (
	"context"
	"fmt"
	"time"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	mgrpb "github.com/NpoolPlatform/message/npool/third/mgr/v1/template/frontend"

	grpc2 "github.com/NpoolPlatform/go-service-framework/pkg/grpc"

	npool "github.com/NpoolPlatform/message/npool/third/mw/v1/template/frontend"

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

func GetFrontendTemplate(ctx context.Context, id string) (*mgrpb.FrontendTemplate, error) {
	info, err := do(ctx, func(_ctx context.Context, cli npool.MiddlewareClient) (cruder.Any, error) {
		resp, err := cli.GetFrontendTemplate(ctx, &npool.GetFrontendTemplateRequest{
			ID: id,
		})
		if err != nil {
			return nil, fmt.Errorf("fail get frontend: %v", err)
		}
		return resp.Info, nil
	})
	if err != nil {
		return nil, fmt.Errorf("fail get frontend: %v", err)
	}
	return info.(*mgrpb.FrontendTemplate), nil
}

func GetFrontendTemplateOnly(ctx context.Context, conds *mgrpb.Conds) (*mgrpb.FrontendTemplate, error) {
	info, err := do(ctx, func(_ctx context.Context, cli npool.MiddlewareClient) (cruder.Any, error) {
		resp, err := cli.GetFrontendTemplateOnly(ctx, &npool.GetFrontendTemplateOnlyRequest{
			Conds: conds,
		})
		if err != nil {
			return nil, fmt.Errorf("fail get frontend: %v", err)
		}
		return resp.Info, nil
	})
	if err != nil {
		return nil, fmt.Errorf("fail get frontend: %v", err)
	}
	return info.(*mgrpb.FrontendTemplate), nil
}

func GetFrontendTemplates(ctx context.Context, conds *mgrpb.Conds, offset, limit uint32) ([]*mgrpb.FrontendTemplate, uint32, error) {
	var total uint32
	infos, err := do(ctx, func(_ctx context.Context, cli npool.MiddlewareClient) (cruder.Any, error) {
		resp, err := cli.GetFrontendTemplates(ctx, &npool.GetFrontendTemplatesRequest{
			Conds:  conds,
			Offset: offset,
			Limit:  limit,
		})
		if err != nil {
			return nil, fmt.Errorf("fail get frontend: %v", err)
		}
		total = resp.GetTotal()
		return resp.Infos, nil
	})
	if err != nil {
		return nil, 0, fmt.Errorf("fail get frontend: %v", err)
	}
	return infos.([]*mgrpb.FrontendTemplate), total, nil
}
