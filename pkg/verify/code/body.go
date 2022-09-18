package code

import (
	"context"
	"fmt"
	"github.com/NpoolPlatform/message/npool/appuser/mgr/v2/signmethod"
	"github.com/NpoolPlatform/message/npool/third/mgr/v1/usedfor"
	"github.com/NpoolPlatform/third-gateway/pkg/middleware/code"
	"github.com/google/uuid"
	"strings"
	"time"
)

const (
	CodeTemplate = "{{ CODE }}"
	NameTemplate = "{{ NAME }}"
)

func BuildBody(
	ctx context.Context,
	appID,
	account string,
	toUsername *string,
	accountType signmethod.SignMethodType,
	usedFor usedfor.UsedFor,
	template string,
) (
	string, error,
) {
	switch usedFor {
	case usedfor.UsedFor_Signup:
		fallthrough // nolint
	case usedfor.UsedFor_Signin:
		fallthrough // nolint
	case usedfor.UsedFor_SetWithdrawAddress:
		fallthrough // nolint
	case usedfor.UsedFor_Withdraw:
		fallthrough //nolint
	case usedfor.UsedFor_SetTransferTargetUser:
		fallthrough //nolint
	case usedfor.UsedFor_Transfer:
		fallthrough //nolint
	case usedfor.UsedFor_Update:
		vCode := code.Generate6NumberCode()
		userCode := code.UserCode{
			AppID:       uuid.MustParse(appID),
			Account:     account,
			AccountType: accountType.String(),
			UsedFor:     usedFor.String(),
			Code:        vCode,
			NextAt:      time.Now().Add(1 * time.Minute),
			ExpireAt:    time.Now().Add(10 * time.Minute),
		}
		if ok, err := code.Nextable(ctx, &userCode); err != nil || !ok {
			return "", fmt.Errorf("wait for next code generation")
		}

		err := code.CreateCodeCache(ctx, &userCode)
		if err != nil {
			return "", fmt.Errorf("fail create code cache: %v", err)
		}

		str := strings.ReplaceAll(template, CodeTemplate, vCode)

		if toUsername != nil && *toUsername != "" {
			str = strings.ReplaceAll(str, NameTemplate, *toUsername)
		}

		return str, nil

	case usedfor.UsedFor_Contact:
		return "", fmt.Errorf("NOT IMPLEMENTED")
	}

	return "", fmt.Errorf("invalid used for")
}
