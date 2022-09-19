package contact

import (
	"context"
	"github.com/NpoolPlatform/message/npool/third/mgr/v1/usedfor"
)

func ContactViaEmail(
	ctx context.Context,
	appID string,
	usedFor usedfor.UsedFor,
	sender,
	subject,
	body,
	senderName string,
) error {

	return nil
}
