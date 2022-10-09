package contact

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/message/npool/third/mgr/v1/usedfor"
	npool "github.com/NpoolPlatform/message/npool/third/mw/v1/contact"
	"github.com/NpoolPlatform/third-middleware/pkg/contact"
	constant "github.com/NpoolPlatform/third-middleware/pkg/message/const"
	commontracer "github.com/NpoolPlatform/third-middleware/pkg/tracer"
	tracer "github.com/NpoolPlatform/third-middleware/pkg/tracer/contact"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) ContactViaEmail(ctx context.Context, in *npool.ContactViaEmailRequest) (resp *npool.ContactViaEmailResponse, err error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "ContactViaEmail")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = tracer.Trace(span, in)

	if _, err = uuid.Parse(in.AppID); err != nil {
		logger.Sugar().Errorw("ContactViaEmail", "error", err)
		return &npool.ContactViaEmailResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	switch in.UsedFor {
	case usedfor.UsedFor_Contact:
	default:
		logger.Sugar().Errorw("NotifyEmail", "UsedFor is invalid")
		return &npool.ContactViaEmailResponse{}, status.Error(codes.InvalidArgument, "UsedFor is invalid")
	}

	if in.GetSender() == "" {
		logger.Sugar().Errorw("NotifyEmail", "error", err)
		return &npool.ContactViaEmailResponse{}, status.Error(codes.InvalidArgument, "Sender is empty")
	}

	if in.GetSubject() == "" {
		logger.Sugar().Errorw("NotifyEmail", "error", err)
		return &npool.ContactViaEmailResponse{}, status.Error(codes.InvalidArgument, "Subject is empty")
	}

	if in.GetBody() == "" {
		logger.Sugar().Errorw("NotifyEmail", "error", err)
		return &npool.ContactViaEmailResponse{}, status.Error(codes.InvalidArgument, "Body is empty")
	}

	if in.GetSenderName() == "" {
		logger.Sugar().Errorw("NotifyEmail", "error", err)
		return &npool.ContactViaEmailResponse{}, status.Error(codes.InvalidArgument, "SenderName is empty")
	}

	span = commontracer.TraceInvoker(span, "verify", "verify", "SendCode")

	err = contact.ContactViaEmail(
		ctx, in.GetAppID(),
		in.GetUsedFor(),
		in.GetSender(),
		in.GetSubject(),
		in.GetBody(),
		in.GetSender(),
	)
	if err != nil {
		logger.Sugar().Errorw("NotifyEmail", "error", err)
		return &npool.ContactViaEmailResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.ContactViaEmailResponse{}, nil
}
