package verify

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool/appuser/mgr/v2/signmethod"
	constant "github.com/NpoolPlatform/third-middleware/pkg/message/const"
	"github.com/NpoolPlatform/third-middleware/pkg/verify"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	scodes "go.opentelemetry.io/otel/codes"

	npool "github.com/NpoolPlatform/message/npool/third/mw/v1/verify"

	commontracer "github.com/NpoolPlatform/third-middleware/pkg/tracer"
	tracer "github.com/NpoolPlatform/third-middleware/pkg/tracer/verify"
)

func (s *Server) SendCode(ctx context.Context, in *npool.SendCodeRequest) (resp *npool.SendCodeResponse, err error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "SendCode")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = tracer.Trace(span, in)

	if _, err = uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("SendCode", "error", err)
		return &npool.SendCodeResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	if _, err = uuid.Parse(in.GetLangID()); err != nil {
		logger.Sugar().Errorw("SendCode", "error", err)
		return &npool.SendCodeResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	if in.GetAccount() == "" {
		logger.Sugar().Errorw("SendCode", "Account", in.GetAccount())
		return &npool.SendCodeResponse{}, status.Error(codes.InvalidArgument, "invalid account")
	}

	switch in.GetAccountType() {
	case signmethod.SignMethodType_Email:
	case signmethod.SignMethodType_Mobile:
	default:
		logger.Sugar().Errorw("SendCode", "error", err)
		return &npool.SendCodeResponse{}, status.Error(codes.InvalidArgument, "invalid account type")
	}

	span = commontracer.TraceInvoker(span, "verify", "verify", "SendCode")

	err = verify.SendCode(ctx, in.GetAppID(), in.GetLangID(), in.GetAccount(), in.ToUsername, in.GetAccountType(), in.GetUsedFor())
	if err != nil {
		logger.Sugar().Errorw("GetAuths", "error", err)
		return &npool.SendCodeResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.SendCodeResponse{}, nil
}

func (s *Server) VerifyCode(ctx context.Context, in *npool.VerifyCodeRequest) (resp *npool.VerifyCodeResponse, err error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "SendCode")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = tracer.TraceVerify(span, in)

	if _, err = uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("SendCode", "error", err)
		return &npool.VerifyCodeResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	if in.GetAccount() == "" {
		logger.Sugar().Errorw("SendCode", "error", err)
		return &npool.VerifyCodeResponse{}, status.Error(codes.InvalidArgument, "Account is empty")
	}

	switch in.GetAccountType() {
	case signmethod.SignMethodType_Email:
	case signmethod.SignMethodType_Mobile:
	case signmethod.SignMethodType_Google:
	default:
		logger.Sugar().Errorw("SendCode", "error", err)
		return &npool.VerifyCodeResponse{}, status.Error(codes.InvalidArgument, "invalid account type")
	}

	span = commontracer.TraceInvoker(span, "verify", "verify", "SendCode")

	err = verify.VerifyCode(ctx, in.GetAppID(), in.GetAccount(), in.GetCode(), in.GetAccountType(), in.GetUsedFor())
	if err != nil {
		logger.Sugar().Errorw("GetAuths", "error", err)
		return &npool.VerifyCodeResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.VerifyCodeResponse{}, nil
}

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
