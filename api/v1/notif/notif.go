package notif

import (
	"context"

	usedfor "github.com/NpoolPlatform/message/npool/notif/mgr/v1/notif"
	"github.com/NpoolPlatform/third-middleware/pkg/notif"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/third/mw/v1/notif"
	constant "github.com/NpoolPlatform/third-middleware/pkg/message/const"
	commontracer "github.com/NpoolPlatform/third-middleware/pkg/tracer"
	tracer "github.com/NpoolPlatform/third-middleware/pkg/tracer/notif"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//nolint:gocyclo
func (s *Server) NotifEmail(ctx context.Context, in *npool.NotifEmailRequest) (resp *npool.NotifEmailResponse, err error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "NotifEmail")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = tracer.Trace(span, in)

	if _, err = uuid.Parse(in.GetAppID()); err != nil {
		logger.Sugar().Errorw("NotifEmail", "error", err)
		return &npool.NotifEmailResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	if in.FromAccount != nil && in.GetFromAccount() == "" {
		logger.Sugar().Errorw("NotifEmail", "error", err)
		return &npool.NotifEmailResponse{}, status.Error(codes.InvalidArgument, "from account account is empty")
	}

	switch in.GetUsedFor() {
	case usedfor.EventType_WithdrawalRequest:
	case usedfor.EventType_WithdrawalCompleted:
	case usedfor.EventType_DepositReceived:
	case usedfor.EventType_KYCApproved:
	case usedfor.EventType_KYCRejected:
	default:
		logger.Sugar().Errorw("NotifEmail", "error", err)
		return &npool.NotifEmailResponse{}, status.Error(codes.InvalidArgument, "UsedFor is invalid")
	}

	if in.GetReceiverAccount() == "" {
		logger.Sugar().Errorw("NotifEmail", "error", err)
		return &npool.NotifEmailResponse{}, status.Error(codes.InvalidArgument, "receiver account is empty")
	}

	if _, err = uuid.Parse(in.GetLangID()); err != nil {
		logger.Sugar().Errorw("NotifEmail", "error", err)
		return &npool.NotifEmailResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	if in.SenderName != nil && in.GetSenderName() == "" {
		logger.Sugar().Errorw("NotifEmail", "error", "SenderName is empty")
		return &npool.NotifEmailResponse{}, status.Error(codes.InvalidArgument, "SenderName is empty")
	}

	if in.ReceiverName != nil && in.GetReceiverName() == "" {
		logger.Sugar().Errorw("NotifEmail", "error", "ReceiverName is empty")
		return &npool.NotifEmailResponse{}, status.Error(codes.InvalidArgument, "ReceiverName is empty")
	}

	if in.Message != nil && in.GetMessage() == "" {
		logger.Sugar().Errorw("NotifEmail", "error", "Message is empty")
		return &npool.NotifEmailResponse{}, status.Error(codes.InvalidArgument, "Message is empty")
	}

	span = commontracer.TraceInvoker(span, "notif", "notif", "NotifEmail")

	err = notif.NotifEmail(
		ctx,
		in.GetAppID(),
		in.FromAccount,
		in.GetUsedFor(),
		in.GetReceiverAccount(),
		in.GetLangID(),
		in.SenderName,
		in.ReceiverName,
		in.Message,
		in.GetUseTemplate(),
		in.GetTitle(),
		in.GetContent(),
	)
	if err != nil {
		logger.Sugar().Errorw("NotifEmail", "error", err)
		return &npool.NotifEmailResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.NotifEmailResponse{}, nil
}
