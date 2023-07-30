package oauth

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	npool "github.com/NpoolPlatform/message/npool/third/mw/v1/oauth"
	oauth1 "github.com/NpoolPlatform/third-middleware/pkg/mw/oauth"
)

func (s *Server) GetOAuthAccessToken(
	ctx context.Context,
	in *npool.GetOAuthAccessTokenRequest,
) (
	resp *npool.GetOAuthAccessTokenResponse,
	err error,
) {
	handler, err := oauth1.NewHandler(
		ctx,
		oauth1.WithClientName(in.GetClientName()),
		oauth1.WithClientID(in.GetClientID()),
		oauth1.WithClientSecret(in.GetClientSecret()),
		oauth1.WithCode(in.GetCode()),
	)
	if err != nil {
		return &npool.GetOAuthAccessTokenResponse{}, status.Error(codes.Internal, err.Error())
	}

	info, err := handler.GetAccessToken(ctx)
	if err != nil {
		logger.Sugar().Errorw("GetOAuthAccessTokenResponse", "error", err)
		return &npool.GetOAuthAccessTokenResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetOAuthAccessTokenResponse{
		Info: info,
	}, nil
}

func (s *Server) GetOAuthUserInfo(
	ctx context.Context,
	in *npool.GetOAuthUserInfoRequest,
) (
	resp *npool.GetOAuthUserInfoResponse,
	err error,
) {
	handler, err := oauth1.NewHandler(
		ctx,
		oauth1.WithClientName(in.GetClientName()),
		oauth1.WithAccessToken(in.GetAccessToken()),
	)
	if err != nil {
		return &npool.GetOAuthUserInfoResponse{}, status.Error(codes.Internal, err.Error())
	}

	info, err := handler.GetThirdUserInfo(ctx)
	if err != nil {
		logger.Sugar().Errorw("GetOAuthUserInfoResponse", "error", err)
		return &npool.GetOAuthUserInfoResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetOAuthUserInfoResponse{
		Info: info,
	}, nil
}
