package verify

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	constant "github.com/NpoolPlatform/third-middleware/pkg/message/const"
	"github.com/NpoolPlatform/third-middleware/pkg/verify"

	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	scodes "go.opentelemetry.io/otel/codes"

	npool "github.com/NpoolPlatform/message/npool/third/mw/v1/verify"

	commontracer "github.com/NpoolPlatform/third-middleware/pkg/tracer"
)

func (s *Server) VerifyGoogleRecaptchaV3(
	ctx context.Context,
	in *npool.VerifyGoogleRecaptchaV3Request,
) (
	resp *npool.VerifyGoogleRecaptchaV3Response,
	err error,
) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "SendCode")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if in.GetRecaptchaToken() == "" {
		logger.Sugar().Errorw("SendCode", "error", err)
		return &npool.VerifyGoogleRecaptchaV3Response{}, status.Error(codes.InvalidArgument, "RecaptchaToken is empty")
	}

	span = commontracer.TraceInvoker(span, "verify", "verify", "SendCode")

	err = verify.VerifyGoogleRecaptchaV3(in.GetRecaptchaToken())
	if err != nil {
		logger.Sugar().Errorw("GetAuths", "error", err)
		return &npool.VerifyGoogleRecaptchaV3Response{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.VerifyGoogleRecaptchaV3Response{}, nil
}
