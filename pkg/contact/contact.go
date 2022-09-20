package contact

import (
	"context"
	"fmt"
	"strings"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	"github.com/NpoolPlatform/message/npool/appuser/mgr/v2/signmethod"
	contactpb "github.com/NpoolPlatform/message/npool/third/mgr/v1/contact"
	"github.com/NpoolPlatform/message/npool/third/mgr/v1/usedfor"
	contact "github.com/NpoolPlatform/third-manager/pkg/crud/v1/contact"
	"github.com/NpoolPlatform/third-middleware/pkg/verify/email"
)

const AccountType = signmethod.SignMethodType_Email

func ContactViaEmail(
	ctx context.Context,
	appID string,
	usedFor usedfor.UsedFor,
	sender,
	subject,
	body,
	senderName string,
) error {
	info, err := contact.RowOnly(ctx, &contactpb.Conds{
		AppID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: appID,
		},
		UsedFor: &npool.Int32Val{
			Op:    cruder.EQ,
			Value: int32(usedFor.Number()),
		},
		AccountType: &npool.Int32Val{
			Op:    cruder.EQ,
			Value: int32(AccountType.Number()),
		},
	})
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("failed get contact")
	}
	sendBody := fmt.Sprintf("From: %v<br>Name: %v<br>%v", sender, senderName, body)
	sendBody = strings.ReplaceAll(sendBody, "\n", "<br>")

	err = email.SendEmailByAWS(subject, sendBody, info.Sender, info.Account, sender)
	if err != nil {
		return fmt.Errorf("fail send email: %v", err)
	}
	return nil
}
