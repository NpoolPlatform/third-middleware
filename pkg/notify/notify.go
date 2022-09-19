package notify

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	emailpb "github.com/NpoolPlatform/message/npool/third/mgr/v1/template/email"
	"github.com/NpoolPlatform/message/npool/third/mgr/v1/usedfor"
	"github.com/NpoolPlatform/third-manager/pkg/crud/v1/template/email"

	"strings"

	verifyemail "github.com/NpoolPlatform/third-middleware/pkg/verify/email"
)

func NotifyEmail(
	ctx context.Context,
	appID,
	fromAccount string,
	usedFor usedfor.UsedFor,
	receiverAccount,
	senderName,
	receiverName,
	langID string,
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

	body := strings.ReplaceAll(template.Body, "{{ FROM }}", fromAccount)
	body = strings.ReplaceAll(body, "{{ TO }}", receiverAccount)
	body = strings.ReplaceAll(body, "{{ RECEIVER }}", receiverName)
	body = strings.ReplaceAll(body, "{{ SENDER }}", senderName)

	return verifyemail.SendEmailByAWS(template.Subject, body, template.Sender, receiverAccount)
}
