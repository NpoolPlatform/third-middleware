package send

import (
	"context"
	"fmt"

	aws "github.com/NpoolPlatform/third-middleware/pkg/aws"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func SendMessage(
	ctx context.Context,
	from, to, subject, content string,
	replyTos, toCCs []string,
	accountType basetypes.SignMethod,
) error {
	switch accountType {
	case basetypes.SignMethod_Email:
		replyTos = append(replyTos, from)
		return aws.SendEmail(subject, content, from, to, replyTos...)
	case basetypes.SignMethod_Mobile:
		return aws.SendSMS(content, to)
	}
	return fmt.Errorf("signmethod is invalid")
}
