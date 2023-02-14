package frontend

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/third/mw/v1/template/frontend"
	constant "github.com/NpoolPlatform/third-middleware/pkg/message/const"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mgrpb "github.com/NpoolPlatform/message/npool/third/mgr/v1/template/frontend"
	mgrcli "github.com/NpoolPlatform/third-manager/pkg/client/template/frontend"

	usedfor "github.com/NpoolPlatform/message/npool/third/mgr/v1/usedfor"
)

func (s *Server) GetFrontendTemplate(
	ctx context.Context,
	in *npool.GetFrontendTemplateRequest,
) (
	resp *npool.GetFrontendTemplateResponse,
	err error,
) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetFrontendTemplate")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if _, err = uuid.Parse(in.GetID()); err != nil {
		logger.Sugar().Errorw("GetFrontendTemplate", "error", err)
		return &npool.GetFrontendTemplateResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := mgrcli.GetFrontendTemplate(ctx, in.GetID())
	if err != nil {
		logger.Sugar().Errorw("GetFrontendTemplate", "error", err)
		return &npool.GetFrontendTemplateResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetFrontendTemplateResponse{
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
		case uint32(usedfor.UsedFor_WithdrawalRequest):
		case uint32(usedfor.UsedFor_WithdrawalCompleted):
		case uint32(usedfor.UsedFor_DepositReceived):
		case uint32(usedfor.UsedFor_KYCApproved):
		case uint32(usedfor.UsedFor_KYCRejected):
		case uint32(usedfor.UsedFor_Announcement):
		default:
			return fmt.Errorf("UsedFor is invalid")
		}
	}
	return nil
}

func (s *Server) GetFrontendTemplates(
	ctx context.Context,
	in *npool.GetFrontendTemplatesRequest,
) (
	resp *npool.GetFrontendTemplatesResponse,
	err error,
) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetFrontendTemplate")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	err = validateConds(in.GetConds())
	if err != nil {
		logger.Sugar().Errorw("GetFrontendTemplate", "error", err)
		return &npool.GetFrontendTemplatesResponse{}, status.Error(codes.Internal, err.Error())
	}
	infos, total, err := mgrcli.GetFrontendTemplates(ctx, in.GetConds(), int32(in.GetOffset()), int32(in.GetLimit()))
	if err != nil {
		logger.Sugar().Errorw("GetFrontendTemplate", "error", err)
		return &npool.GetFrontendTemplatesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetFrontendTemplatesResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetFrontendTemplateOnly(
	ctx context.Context,
	in *npool.GetFrontendTemplateOnlyRequest,
) (
	resp *npool.GetFrontendTemplateOnlyResponse,
	err error,
) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetFrontendTemplateOnly")
	defer span.End()
	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	err = validateConds(in.GetConds())
	if err != nil {
		logger.Sugar().Errorw("GetFrontendTemplateOnly", "error", err)
		return &npool.GetFrontendTemplateOnlyResponse{}, status.Error(codes.Internal, err.Error())
	}
	info, err := mgrcli.GetFrontendTemplateOnly(ctx, in.GetConds())
	if err != nil {
		logger.Sugar().Errorw("GetFrontendTemplateOnly", "error", err)
		return &npool.GetFrontendTemplateOnlyResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetFrontendTemplateOnlyResponse{
		Info: info,
	}, nil
}
