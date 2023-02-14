package template

import (
	"strings"

	constant "github.com/NpoolPlatform/third-middleware/pkg/const"
)

func ReplaceVariable(
	content string,
	userName, message, amount, coinUnit, date, time, address *string,
) string {
	if userName != nil {
		content = strings.ReplaceAll(content, constant.NotifTemplateVarName, *userName)
	}
	if message != nil {
		content = strings.ReplaceAll(content, constant.NotifTemplateVarMessage, *message)
	}
	if amount != nil {
		content = strings.ReplaceAll(content, constant.NotifTemplateVarAmount, *amount)
	}
	if coinUnit != nil {
		content = strings.ReplaceAll(content, constant.NotifTemplateVarCoinUnit, *coinUnit)
	}
	if date != nil {
		content = strings.ReplaceAll(content, constant.NotifTemplateVarDate, *date)
	}
	if time != nil {
		content = strings.ReplaceAll(content, constant.NotifTemplateVarTime, *time)
	}
	if address != nil {
		content = strings.ReplaceAll(content, constant.NotifTemplateVarAddress, *address)
	}
	return content
}
