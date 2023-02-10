package notif

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/third/mw/v1/notif"
	constant "github.com/NpoolPlatform/third-middleware/pkg/message/const"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/NpoolPlatform/third-middleware/pkg/verify/email"
)

func (s *Server) SendNotifEmail(
	ctx context.Context,
	in *npool.SendNotifEmailRequest,
) (
	resp *npool.SendNotifEmailResponse,
	err error,
) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "SendNotifEmail")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if in.GetSubject() == "" {
		logger.Sugar().Errorw("SendNotifEmail", "error", err)
		return &npool.SendNotifEmailResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	if in.GetContent() == "" {
		logger.Sugar().Errorw("SendNotifEmail", "error", err)
		return &npool.SendNotifEmailResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	if in.GetFrom() == "" {
		logger.Sugar().Errorw("SendNotifEmail", "error", err)
		return &npool.SendNotifEmailResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	if in.GetTo() == "" {
		logger.Sugar().Errorw("SendNotifEmail", "error", err)
		return &npool.SendNotifEmailResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	err = email.SendEmailByAWS(in.GetSubject(), in.GetContent(), in.GetFrom(), in.GetTo())
	if err != nil {
		logger.Sugar().Errorw("SendNotifEmail", "error", err)
		return &npool.SendNotifEmailResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	return &npool.SendNotifEmailResponse{}, nil
}