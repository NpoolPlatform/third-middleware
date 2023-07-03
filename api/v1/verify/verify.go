package verify

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	npool "github.com/NpoolPlatform/message/npool/third/mw/v1/verify"
	verify1 "github.com/NpoolPlatform/third-middleware/pkg/mw/verify"
)

func (s *Server) VerifyGoogleRecaptchaV3(
	ctx context.Context,
	in *npool.VerifyGoogleRecaptchaV3Request,
) (
	resp *npool.VerifyGoogleRecaptchaV3Response,
	err error,
) {
	handler, err := verify1.NewHandler(
		ctx,
		verify1.WithRecaptchaToken(in.GetRecaptchaToken()),
	)
	if err != nil {
		return &npool.VerifyGoogleRecaptchaV3Response{}, status.Error(codes.Internal, err.Error())
	}

	err = handler.VerifyGoogleRecaptchaV3(ctx)
	if err != nil {
		logger.Sugar().Errorw("GetAuths", "error", err)
		return &npool.VerifyGoogleRecaptchaV3Response{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.VerifyGoogleRecaptchaV3Response{}, nil
}
