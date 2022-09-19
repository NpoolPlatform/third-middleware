package contact

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	trace1 "go.opentelemetry.io/otel/trace"

	npool "github.com/NpoolPlatform/message/npool/third/mw/v1/contact"
)

func trace(span trace1.Span, in *npool.ContactViaEmailRequest, index int) trace1.Span {
	span.SetAttributes(
		attribute.String(fmt.Sprintf("AppID.%v", index), in.GetAppID()),
		attribute.String(fmt.Sprintf("UsedFor.%v", index), in.GetUsedFor().String()),
		attribute.String(fmt.Sprintf("Sender.%v", index), in.GetSender()),
		attribute.String(fmt.Sprintf("Subject.%v", index), in.GetSubject()),
		attribute.String(fmt.Sprintf("Body.%v", index), in.GetBody()),
		attribute.String(fmt.Sprintf("SenderName.%v", index), in.GetSenderName()),
	)
	return span
}

func Trace(span trace1.Span, in *npool.ContactViaEmailRequest) trace1.Span {
	return trace(span, in, 0)
}
