//nolint:nolintlint,dupl
package email

import (
	"context"
	"fmt"
	"time"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	mgrpb "github.com/NpoolPlatform/message/npool/third/mgr/v1/template/email"

	grpc2 "github.com/NpoolPlatform/go-service-framework/pkg/grpc"

	npool "github.com/NpoolPlatform/message/npool/third/mw/v1/template/email"

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

func GetEmailTemplate(ctx context.Context, id string) (*mgrpb.EmailTemplate, error) {
	info, err := do(ctx, func(_ctx context.Context, cli npool.MiddlewareClient) (cruder.Any, error) {
		resp, err := cli.GetEmailTemplate(ctx, &npool.GetEmailTemplateRequest{
			ID: id,
		})
		if err != nil {
			return nil, fmt.Errorf("fail get email: %v", err)
		}
		return resp.Info, nil
	})
	if err != nil {
		return nil, fmt.Errorf("fail get email: %v", err)
	}
	return info.(*mgrpb.EmailTemplate), nil
}

func GetEmailTemplateOnly(ctx context.Context, conds *mgrpb.Conds) (*mgrpb.EmailTemplate, error) {
	info, err := do(ctx, func(_ctx context.Context, cli npool.MiddlewareClient) (cruder.Any, error) {
		resp, err := cli.GetEmailTemplateOnly(ctx, &npool.GetEmailTemplateOnlyRequest{
			Conds: conds,
		})
		if err != nil {
			return nil, fmt.Errorf("fail get email: %v", err)
		}
		return resp.Info, nil
	})
	if err != nil {
		return nil, fmt.Errorf("fail get email: %v", err)
	}
	return info.(*mgrpb.EmailTemplate), nil
}

func GetEmailTemplates(ctx context.Context, conds *mgrpb.Conds, offset, limit uint32) ([]*mgrpb.EmailTemplate, uint32, error) {
	var total uint32
	infos, err := do(ctx, func(_ctx context.Context, cli npool.MiddlewareClient) (cruder.Any, error) {
		resp, err := cli.GetEmailTemplates(ctx, &npool.GetEmailTemplatesRequest{
			Conds:  conds,
			Offset: offset,
			Limit:  limit,
		})
		if err != nil {
			return nil, fmt.Errorf("fail get email: %v", err)
		}
		total = resp.GetTotal()
		return resp.Infos, nil
	})
	if err != nil {
		return nil, 0, fmt.Errorf("fail get email: %v", err)
	}
	return infos.([]*mgrpb.EmailTemplate), total, nil
}
