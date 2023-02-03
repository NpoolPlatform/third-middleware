package notify

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	trace1 "go.opentelemetry.io/otel/trace"

	npool "github.com/NpoolPlatform/message/npool/third/mw/v1/notif"
)

func trace(span trace1.Span, in *npool.NotifEmailRequest, index int) trace1.Span {
	span.SetAttributes(
		attribute.String(fmt.Sprintf("AppID.%v", index), in.GetAppID()),
		attribute.String(fmt.Sprintf("FromAccount.%v", index), in.GetFromAccount()),
		attribute.String(fmt.Sprintf("UsedFor.%v", index), in.GetUsedFor().String()),
		attribute.String(fmt.Sprintf("ReceiverAccount.%v", index), in.GetReceiverAccount()),
		attribute.String(fmt.Sprintf("LangID.%v", index), in.GetLangID()),
		attribute.String(fmt.Sprintf("SenderName.%v", index), in.GetSenderName()),
		attribute.String(fmt.Sprintf("ReceiverName.%v", index), in.GetReceiverName()),
	)
	return span
}

func Trace(span trace1.Span, in *npool.NotifEmailRequest) trace1.Span {
	return trace(span, in, 0)
}
