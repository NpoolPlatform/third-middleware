package send

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	constant "github.com/NpoolPlatform/third-middleware/pkg/message/const"
	"github.com/NpoolPlatform/third-middleware/pkg/send"

	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	scodes "go.opentelemetry.io/otel/codes"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/third/mw/v1/send"

	commontracer "github.com/NpoolPlatform/third-middleware/pkg/tracer"
)

func (s *Server) SendMessage(
	ctx context.Context,
	in *npool.SendMessageRequest,
) (
	resp *npool.SendMessageResponse,
	err error,
) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "SendMessage")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if in.GetTo() == "" {
		logger.Sugar().Errorw("SendMessage", "To", in.GetTo())
		return &npool.SendMessageResponse{}, status.Error(codes.InvalidArgument, "To is invalid")
	}
	if in.GetContent() == "" {
		logger.Sugar().Errorw("SendMessage", "Content", in.GetContent())
		return &npool.SendMessageResponse{}, status.Error(codes.InvalidArgument, "Content is invalid")
	}

	switch in.GetAccountType() {
	case basetypes.SignMethod_Email:
		if in.GetFrom() == "" {
			logger.Sugar().Errorw("SendMessage", "From", in.GetFrom())
			return &npool.SendMessageResponse{}, status.Error(codes.InvalidArgument, "From is invalid")
		}
		if in.GetSubject() == "" {
			logger.Sugar().Errorw("SendMessage", "Subject", in.GetSubject())
			return &npool.SendMessageResponse{}, status.Error(codes.InvalidArgument, "Subject is invalid")
		}
	case basetypes.SignMethod_Mobile:
	default:
		logger.Sugar().Errorw("SendMessage", "AccountType", in.GetAccountType())
		return &npool.SendMessageResponse{}, status.Error(codes.InvalidArgument, "AccountType is invalid")
	}

	span = commontracer.TraceInvoker(span, "send", "send", "SendMessage")

	err = send.SendMessage(
		ctx,
		in.GetFrom(),
		in.GetTo(),
		in.GetSubject(),
		in.GetContent(),
		in.GetReplyTos(),
		in.GetToCCs(),
		in.GetAccountType(),
	)
	if err != nil {
		logger.Sugar().Errorw("SendMessage", "error", err)
		return &npool.SendMessageResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.SendMessageResponse{}, nil
}
