//nolint:nolintlint,dupl
package verify

import (
	"context"
	"time"

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

func VerifyGoogleRecaptchaV3(
	ctx context.Context,
	recaptchaToken string,
) error {
	err := do(ctx, func(_ctx context.Context, cli npool.MiddlewareClient) error {
		_, err := cli.VerifyGoogleRecaptchaV3(ctx, &npool.VerifyGoogleRecaptchaV3Request{
			RecaptchaToken: recaptchaToken,
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
