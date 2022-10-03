package google

import (
	"fmt"
	"time"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	constant "github.com/NpoolPlatform/third-middleware/pkg/const"
	"github.com/go-resty/resty/v2"
)

type verifyResponse struct {
	Success     bool      `json:"success"`
	Score       float64   `json:"score"`
	Action      string    `json:"action"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

func VerifyGoogleRecaptchaV3(recaptchaToken string) error {
	hostname := config.GetStringValueWithNameSpace("", config.KeyHostname)
	recaptchaURL := config.GetStringValueWithNameSpace(hostname, constant.GoogleRecaptchaV3URL)
	recaptchaSecret := config.GetStringValueWithNameSpace(hostname, constant.GoogleRecaptchaV3Secret)

	if recaptchaURL == "" || recaptchaSecret == "" {
		return fmt.Errorf("invalid recaptcha parameter")
	}

	url := fmt.Sprintf("%v?secret=%v&response=%v", recaptchaURL, recaptchaSecret, recaptchaToken)
	resp, err := resty.New().R().SetResult(&verifyResponse{}).Post(url)
	if err != nil {
		return fmt.Errorf("fail verify google recaptcha v3: %v", err)
	}

	vResp, ok := resp.Result().(*verifyResponse)
	if !ok {
		return fmt.Errorf("fail get response")
	}

	if !vResp.Success {
		return fmt.Errorf("fail verify google recaptcha v3: %v", vResp.ErrorCodes)
	}

	return nil
}
