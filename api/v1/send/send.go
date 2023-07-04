package send

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	npool "github.com/NpoolPlatform/message/npool/third/mw/v1/send"
	send1 "github.com/NpoolPlatform/third-middleware/pkg/mw/send"
)

func (s *Server) SendMessage(
	ctx context.Context,
	in *npool.SendMessageRequest,
) (
	resp *npool.SendMessageResponse,
	err error,
) {
	handler, err := send1.NewHandler(
		ctx,
		send1.WithSubject(in.GetSubject()),
		send1.WithContent(in.GetContent()),
		send1.WithFrom(in.GetFrom()),
		send1.WithTo(in.GetTo()),
		send1.WithToCCs(in.GetToCCs()),
		send1.WithReplyTos(in.GetReplyTos()),
		send1.WithAccountType(in.GetAccountType()),
	)
	if err != nil {
		return &npool.SendMessageResponse{}, status.Error(codes.Internal, err.Error())
	}

	err = handler.SendMessage(ctx)
	if err != nil {
		logger.Sugar().Errorw("SendMessage", "error", err)
		return &npool.SendMessageResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.SendMessageResponse{}, nil
}
