package verify

import (
	"github.com/NpoolPlatform/third-middleware/pkg/google"
)

func VerifyGoogleRecaptchaV3(token string) error {
	return google.VerifyGoogleRecaptchaV3(token)
}
