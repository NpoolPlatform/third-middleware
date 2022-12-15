package notify

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/third/mw/v1/notify"
	constant "github.com/NpoolPlatform/third-middleware/pkg/message/const"
	"github.com/NpoolPlatform/third-middleware/pkg/notify"
	commontracer "github.com/NpoolPlatform/third-middleware/pkg/tracer"
	tracer "github.com/NpoolPlatform/third-middleware/pkg/tracer/notify"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) NotifyEmail(ctx context.Context, in *npool.NotifyEmailRequest) (resp *npool.NotifyEmailResponse, err error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "NotifyEmail")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = tracer.Trace(span, in)

	if _, err = uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("NotifyEmail", "error", err)
		return &npool.NotifyEmailResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	if in.GetFromAccount() == "" {
		logger.Sugar().Errorw("NotifyEmail", "error", err)
		return &npool.NotifyEmailResponse{}, status.Error(codes.InvalidArgument, "from account is empty")
	}

	if in.GetReceiverAccount() == "" {
		logger.Sugar().Errorw("NotifyEmail", "error", err)
		return &npool.NotifyEmailResponse{}, status.Error(codes.InvalidArgument, "receiver account is empty")
	}

	if _, err = uuid.Parse(in.GetLangID()); err != nil {
		logger.Sugar().Errorw("NotifyEmail", "error", err)
		return &npool.NotifyEmailResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	span = commontracer.TraceInvoker(span, "verify", "verify", "SendCode")

	err = notify.NotifyEmail(
		ctx,
		in.GetAppID(),
		in.GetFromAccount(),
		in.GetUsedFor(),
		in.GetReceiverAccount(),
		in.GetSenderName(),
		in.GetReceiverName(),
		in.GetLangID(),
	)
	if err != nil {
		logger.Sugar().Errorw("NotifyEmail", "error", err)
		return &npool.NotifyEmailResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.NotifyEmailResponse{}, nil
}
