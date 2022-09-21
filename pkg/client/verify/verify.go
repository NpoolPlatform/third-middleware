//nolint:nolintlint,dupl
package verify

import (
	"context"
	"time"

	"github.com/NpoolPlatform/message/npool/appuser/mgr/v2/signmethod"
	"github.com/NpoolPlatform/message/npool/third/mgr/v1/usedfor"

	grpc2 "github.com/NpoolPlatform/go-service-framework/pkg/grpc"

	npool "github.com/NpoolPlatform/message/npool/third/mw/v1/verify"

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

func SendCode(ctx context.Context,
	appID,
	langID,
	account string,
	accountType signmethod.SignMethodType,
	usedFor usedfor.UsedFor,
	toUsername *string,
) error {
	err := do(ctx, func(_ctx context.Context, cli npool.MiddlewareClient) error {
		_, err := cli.SendCode(ctx, &npool.SendCodeRequest{
			AppID:       appID,
			LangID:      langID,
			Account:     account,
			AccountType: accountType,
			UsedFor:     usedFor,
			ToUsername:  toUsername,
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

func VerifyCode(
	ctx context.Context,
	appID string,
	account *string,
	code string,
	accountType signmethod.SignMethodType,
	usedFor usedfor.UsedFor,
) error {
	err := do(ctx, func(_ctx context.Context, cli npool.MiddlewareClient) error {
		_, err := cli.VerifyCode(ctx, &npool.VerifyCodeRequest{
			AppID:       appID,
			Account:     account,
			AccountType: accountType,
			UsedFor:     usedFor,
			Code:        code,
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
