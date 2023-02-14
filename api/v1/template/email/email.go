package email

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/third/mw/v1/template/email"
	constant "github.com/NpoolPlatform/third-middleware/pkg/message/const"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mgrpb "github.com/NpoolPlatform/message/npool/third/mgr/v1/template/email"
	mgrcli "github.com/NpoolPlatform/third-manager/pkg/client/template/email"

	usedfor "github.com/NpoolPlatform/message/npool/third/mgr/v1/usedfor"
)

func (s *Server) GetEmailTemplate(
	ctx context.Context,
	in *npool.GetEmailTemplateRequest,
) (
	resp *npool.GetEmailTemplateResponse,
	err error,
) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetEmailTemplate")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err = uuid.Parse(in.GetID()); err != nil {
		logger.Sugar().Errorw("GetEmailTemplate", "error", err)
		return &npool.GetEmailTemplateResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := mgrcli.GetEmailTemplate(ctx, in.GetID())
	if err != nil {
		logger.Sugar().Errorw("GetEmailTemplate", "error", err)
		return &npool.GetEmailTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetEmailTemplateResponse{
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
		case int32(usedfor.UsedFor_WithdrawalRequest):
		case int32(usedfor.UsedFor_WithdrawalCompleted):
		case int32(usedfor.UsedFor_DepositReceived):
		case int32(usedfor.UsedFor_KYCApproved):
		case int32(usedfor.UsedFor_KYCRejected):
		case int32(usedfor.UsedFor_Announcement):
		default:
			return fmt.Errorf("UsedFor is invalid")
		}
	}
	return nil
}

func (s *Server) GetEmailTemplates(
	ctx context.Context,
	in *npool.GetEmailTemplatesRequest,
) (
	resp *npool.GetEmailTemplatesResponse,
	err error,
) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetEmailTemplate")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	err = validateConds(in.GetConds())
	if err != nil {
		logger.Sugar().Errorw("GetEmailTemplate", "error", err)
		return &npool.GetEmailTemplatesResponse{}, status.Error(codes.Internal, err.Error())
	}
	infos, total, err := mgrcli.GetEmailTemplates(ctx, in.GetConds(), int32(in.GetOffset()), int32(in.GetLimit()))
	if err != nil {
		logger.Sugar().Errorw("GetEmailTemplate", "error", err)
		return &npool.GetEmailTemplatesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetEmailTemplatesResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetEmailTemplateOnly(
	ctx context.Context,
	in *npool.GetEmailTemplateOnlyRequest,
) (
	resp *npool.GetEmailTemplateOnlyResponse,
	err error,
) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetEmailTemplateOnly")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	err = validateConds(in.GetConds())
	if err != nil {
		logger.Sugar().Errorw("GetEmailTemplateOnly", "error", err)
		return &npool.GetEmailTemplateOnlyResponse{}, status.Error(codes.Internal, err.Error())
	}
	info, err := mgrcli.GetEmailTemplateOnly(ctx, in.GetConds())
	if err != nil {
		logger.Sugar().Errorw("GetEmailTemplateOnly", "error", err)
		return &npool.GetEmailTemplateOnlyResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetEmailTemplateOnlyResponse{
		Info: info,
	}, nil
}
