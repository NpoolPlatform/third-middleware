package notif

import (
	"strings"

	constant "github.com/NpoolPlatform/third-middleware/pkg/const"
)

func ReplaceVariable(content string, userName, message *string) string {
	if userName != nil {
		content = strings.ReplaceAll(content, constant.NotifTemplateVarName, *userName)
	}
	if message != nil {
		content = strings.ReplaceAll(content, constant.NotifTemplateVarMessage, *message)
	}
	return content
}
