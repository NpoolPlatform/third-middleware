package template

import (
	"strings"

	constant "github.com/NpoolPlatform/third-middleware/pkg/const"
)

type TemplateVars struct {
	Username *string
	Message  *string
	Amount   *string
	CoinUnit *string
	Date     *string
	Time     *string
	Address  *string
}

func FillTemplate(template string, vars *TemplateVars) string {
	if vars.Username != nil {
		template = strings.ReplaceAll(template, constant.NotifTemplateVarName, *vars.Username)
	}
	if vars.Message != nil {
		template = strings.ReplaceAll(template, constant.NotifTemplateVarMessage, *vars.Message)
	}
	if vars.Amount != nil {
		template = strings.ReplaceAll(template, constant.NotifTemplateVarAmount, *vars.Amount)
	}
	if vars.CoinUnit != nil {
		template = strings.ReplaceAll(template, constant.NotifTemplateVarCoinUnit, *vars.CoinUnit)
	}
	if vars.Date != nil {
		template = strings.ReplaceAll(template, constant.NotifTemplateVarDate, *vars.Date)
	}
	if vars.Time != nil {
		template = strings.ReplaceAll(template, constant.NotifTemplateVarTime, *vars.Time)
	}
	if vars.Address != nil {
		template = strings.ReplaceAll(template, constant.NotifTemplateVarAddress, *vars.Address)
	}
	return template
}
