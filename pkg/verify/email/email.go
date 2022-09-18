package email

import (
	"context"
	"fmt"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appuser/mgr/v2/signmethod"
	"github.com/NpoolPlatform/message/npool/third/mgr/v1/usedfor"
	"github.com/NpoolPlatform/third-middleware/pkg/verify/code"
	"strings"

	emailpb "github.com/NpoolPlatform/message/npool/third/mgr/v1/template/email"
	"github.com/NpoolPlatform/third-manager/pkg/crud/v1/template/email"
)

func SendCode(
	ctx context.Context,
	appID,
	langID,
	account string,
	toUsername *string,
	accountType signmethod.SignMethodType,
	usedFor usedfor.UsedFor,
) error {
	template, err := email.RowOnly(ctx, &emailpb.Conds{
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

	toUsernameN := &template.DefaultToUsername
	if toUsername != nil && *toUsername != "" {
		toUsernameN = toUsername
	}

	body, err := code.BuildBody(ctx, appID, account, toUsernameN, accountType, usedFor, template.Body)
	if err != nil {
		return fmt.Errorf("fail build email body: %v", err)
	}

	if toUsername == nil || *toUsername == "" {
		body = strings.ReplaceAll(body, code.NameTemplate, template.DefaultToUsername)
	}

	err = sendEmailByAWS(template.Subject, body, template.Sender, account)
	if err != nil {
		return fmt.Errorf("fail to send email: %v", err)
	}

	return nil
}
