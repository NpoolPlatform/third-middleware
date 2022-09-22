package verify

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/message/npool/appuser/mgr/v2/signmethod"
	"github.com/NpoolPlatform/message/npool/third/mgr/v1/usedfor"
	verifycode "github.com/NpoolPlatform/third-middleware/pkg/verify/code"
	"github.com/NpoolPlatform/third-middleware/pkg/verify/email"
	"github.com/NpoolPlatform/third-middleware/pkg/verify/google"
	"github.com/NpoolPlatform/third-middleware/pkg/verify/sms"
	"github.com/google/uuid"
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
	switch accountType {
	case signmethod.SignMethodType_Email:
		err := email.SendCode(ctx, appID, langID, account, toUsername, accountType, usedFor)
		if err != nil {
			return err
		}
	case signmethod.SignMethodType_Mobile:
		err := sms.SendCode(ctx, appID, langID, account, accountType, usedFor)
		if err != nil {
			return err
		}
	}
	return nil
}

func VerifyCode(
	ctx context.Context,
	appID string,
	account string,
	code string,
	accountType signmethod.SignMethodType,
	usedFor usedfor.UsedFor,
) error {
	if accountType == signmethod.SignMethodType_Google {
		verified, err := google.VerifyCode(account, code)
		if err != nil {
			return err
		}
		if !verified {
			return fmt.Errorf("invalid code: %v", err)
		}
		return nil
	}
	userCode := verifycode.UserCode{
		AppID:       uuid.MustParse(appID),
		Account:     account,
		AccountType: accountType.String(),
		UsedFor:     usedFor.String(),
		Code:        code,
	}

	err := verifycode.VerifyCodeCache(ctx, &userCode)
	if err != nil {
		return fmt.Errorf("invalid code: %v", err)
	}

	return nil
}

func VerifyGoogleRecaptchaV3(token string) error {
	return google.VerifyGoogleRecaptchaV3(token)
}
