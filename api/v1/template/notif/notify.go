package notif

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/third/mw/v1/template/notif"
	constant "github.com/NpoolPlatform/third-middleware/pkg/message/const"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mgrpb "github.com/NpoolPlatform/message/npool/third/mgr/v1/template/notif"
	mgrcli "github.com/NpoolPlatform/third-manager/pkg/client/template/notif"

	usedfor "github.com/NpoolPlatform/message/npool/notif/mgr/v1/notif"
)

func (s *Server) GetNotifTemplate(
	ctx context.Context,
	in *npool.GetNotifTemplateRequest,
) (
	resp *npool.GetNotifTemplateResponse,
	err error,
) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetNotifTemplate")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err = uuid.Parse(in.GetID()); err != nil {
		logger.Sugar().Errorw("GetNotifTemplate", "error", err)
		return &npool.GetNotifTemplateResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := mgrcli.GetNotifTemplate(ctx, in.GetID())
	if err != nil {
		logger.Sugar().Errorw("GetNotifTemplate", "error", err)
		return &npool.GetNotifTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetNotifTemplateResponse{
		Info: info,
	}, nil
}

func validateConds(in *mgrpb.Conds) error {
	if in.ID != nil {
		if _, err := uuid.Parse(in.GetID().GetValue()); err != nil {
			logger.Sugar().Errorw("validateConds", "ID", in.GetID().GetValue(), "error", err)
			return err
		}
	}
	if in.AppID != nil {
		if _, err := uuid.Parse(in.GetAppID().GetValue()); err != nil {
			logger.Sugar().Errorw("validateConds", "AppID", in.GetAppID().GetValue(), "error", err)
			return err
		}
	}
	if in.LangID != nil {
		if _, err := uuid.Parse(in.GetLangID().GetValue()); err != nil {
			logger.Sugar().Errorw("validateConds", "LangID", in.GetLangID().GetValue(), "error", err)
			return err
		}
	}
	if in.UsedFor != nil {
		switch in.GetUsedFor().GetValue() {
		case uint32(usedfor.EventType_WithdrawalRequest):
		case uint32(usedfor.EventType_WithdrawalCompleted):
		case uint32(usedfor.EventType_DepositReceived):
		case uint32(usedfor.EventType_KYCApproved):
		case uint32(usedfor.EventType_KYCRejected):
		case uint32(usedfor.EventType_Announcement):
		default:
			return fmt.Errorf("EventType is invalid")
		}
	}
	return nil
}

func (s *Server) GetNotifTemplates(
	ctx context.Context,
	in *npool.GetNotifTemplatesRequest,
) (
	resp *npool.GetNotifTemplatesResponse,
	err error,
) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetNotifTemplate")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	err = validateConds(in.GetConds())
	if err != nil {
		logger.Sugar().Errorw("GetNotifTemplate", "error", err)
		return &npool.GetNotifTemplatesResponse{}, status.Error(codes.Internal, err.Error())
	}
	infos, total, err := mgrcli.GetNotifTemplates(ctx, in.GetConds(), int32(in.GetOffset()), int32(in.GetLimit()))
	if err != nil {
		logger.Sugar().Errorw("GetNotifTemplate", "error", err)
		return &npool.GetNotifTemplatesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetNotifTemplatesResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetNotifTemplateOnly(
	ctx context.Context,
	in *npool.GetNotifTemplateOnlyRequest,
) (
	resp *npool.GetNotifTemplateOnlyResponse,
	err error,
) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetNotifTemplateOnly")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	err = validateConds(in.GetConds())
	if err != nil {
		logger.Sugar().Errorw("GetNotifTemplateOnly", "error", err)
		return &npool.GetNotifTemplateOnlyResponse{}, status.Error(codes.Internal, err.Error())
	}
	info, err := mgrcli.GetNotifTemplateOnly(ctx, in.GetConds())
	if err != nil {
		logger.Sugar().Errorw("GetNotifTemplateOnly", "error", err)
		return &npool.GetNotifTemplateOnlyResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetNotifTemplateOnlyResponse{
		Info: info,
	}, nil
}
