//nolint:nolintlint,dupl
package oauth

import (
	"context"
	"time"

	grpc2 "github.com/NpoolPlatform/go-service-framework/pkg/grpc"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	npool "github.com/NpoolPlatform/message/npool/third/mw/v1/oauth"

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

func GetOAuthAccessToken(ctx context.Context, clientName basetypes.SignMethod, clientID, clientSecret, code string) (*string, error) {
	info, err := do(ctx, func(_ctx context.Context, cli npool.MiddlewareClient) (cruder.Any, error) {
		resp, err := cli.GetOAuthAccessToken(ctx, &npool.GetOAuthAccessTokenRequest{
			ClientName:   clientName,
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Code:         code,
		})
		if err != nil {
			return nil, err
		}
		return resp.AccessToken, nil
	})
	if err != nil {
		return nil, err
	}
	return info.(*string), nil
}

func GetOAuthUserInfo(ctx context.Context, clientName basetypes.SignMethod, token string) (*npool.ThirdUserInfo, error) {
	info, err := do(ctx, func(_ctx context.Context, cli npool.MiddlewareClient) (cruder.Any, error) {
		resp, err := cli.GetOAuthUserInfo(ctx, &npool.GetOAuthUserInfoRequest{
			ClientName:  clientName,
			AccessToken: token,
		})
		if err != nil {
			return nil, err
		}
		return resp.Info, nil
	})
	if err != nil {
		return nil, err
	}
	return info.(*npool.ThirdUserInfo), nil
}
