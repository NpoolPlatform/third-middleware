package sms

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appuser/mgr/v2/signmethod"
	smspb "github.com/NpoolPlatform/message/npool/third/mgr/v1/template/sms"
	"github.com/NpoolPlatform/message/npool/third/mgr/v1/usedfor"
	"github.com/NpoolPlatform/third-manager/pkg/crud/v1/template/sms"
	"github.com/NpoolPlatform/third-middleware/pkg/verify/code"
)

func SendCode(
	ctx context.Context,
	appID,
	langID,
	account string,
	accountType signmethod.SignMethodType,
	usedFor usedfor.UsedFor,
) error {
	template, err := sms.RowOnly(ctx, &smspb.Conds{
		AppID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: appID,
		},
		LangID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: langID,
		},
		UsedFor: &npool.Int32Val{
			Op:    cruder.EQ,
			Value: int32(usedFor.Number()),
		},
	})
	if err != nil {
		return err
	}
	if template == nil {
		return fmt.Errorf("failed get email template")
	}

	body, err := code.BuildBody(ctx, appID, account, nil, accountType, usedFor, template.Message)
	if err != nil {
		return fmt.Errorf("fail build email body: %v", err)
	}

	err = sendSMSByAWS(body, account)
	if err != nil {
		return fmt.Errorf("fail to send email: %v", err)
	}

	return nil
}
