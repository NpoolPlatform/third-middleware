package notif

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/message/npool"
	usedfor "github.com/NpoolPlatform/message/npool/notif/mgr/v1/notif"
	notifpb "github.com/NpoolPlatform/message/npool/third/mgr/v1/template/notif"

	notifcli "github.com/NpoolPlatform/third-manager/pkg/client/template/notif"

	"strings"

	verifyemail "github.com/NpoolPlatform/third-middleware/pkg/verify/email"
)

func NotifEmail(
	ctx context.Context,
	appID string,
	fromAccount *string,
	usedFor usedfor.EventType,
	receiverAccount,
	langID string,
	senderName,
	receiverName,
	message *string,
	useTemplate bool,
	title,
	content string,
) error {
	subject := title
	body := content

	template, err := notifcli.GetNotifTemplateOnly(ctx, &notifpb.Conds{
		AppID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: appID,
		},
		LangID: &npool.StringVal{
			Op:    cruder.EQ,
			Value: langID,
		},
		UsedFor: &npool.Uint32Val{
			Op:    cruder.EQ,
			Value: uint32(usedFor.Number()),
		},
	})
	if err != nil {
		return err
	}
	if template == nil {
		return fmt.Errorf("failed get email template")
	}

	if useTemplate {
		body = template.Content

		if fromAccount != nil {
			body = strings.ReplaceAll(body, "{{ FROM }}", *fromAccount)
		}

		body = strings.ReplaceAll(body, "{{ TO }}", receiverAccount)

		if senderName != nil {
			body = strings.ReplaceAll(body, "{{ SENDER }}", *senderName)
		}

		if receiverName != nil {
			body = strings.ReplaceAll(body, "{{ RECEIVER }}", *receiverName)
		}

		if message != nil {
			body = strings.ReplaceAll(body, "{{ MESSAGE }}", *message)
		}
	}

	return verifyemail.SendEmailByAWS(subject, body, template.Sender, receiverAccount)
}
