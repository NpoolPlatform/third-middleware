package verify

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	trace1 "go.opentelemetry.io/otel/trace"

	npool "github.com/NpoolPlatform/message/npool/third/mw/v1/verify"
)

func trace(span trace1.Span, in *npool.SendCodeRequest, index int) trace1.Span {
	span.SetAttributes(
		attribute.String(fmt.Sprintf("AppID.%v", index), in.GetAppID()),
		attribute.String(fmt.Sprintf("Account.%v", index), in.GetAccount()),
		attribute.String(fmt.Sprintf("AccountType.%v", index), in.GetAccountType().String()),
		attribute.String(fmt.Sprintf("UsedFor.%v", index), in.GetUsedFor().String()),
	)
	return span
}

func Trace(span trace1.Span, in *npool.SendCodeRequest) trace1.Span {
	return trace(span, in, 0)
}

func traceVerify(span trace1.Span, in *npool.VerifyCodeRequest, index int) trace1.Span {
	span.SetAttributes(
		attribute.String(fmt.Sprintf("AppID.%v", index), in.GetAppID()),
		attribute.String(fmt.Sprintf("Account.%v", index), in.GetAccount()),
		attribute.String(fmt.Sprintf("AccountType.%v", index), in.GetAccountType().String()),
		attribute.String(fmt.Sprintf("UsedFor.%v", index), in.GetUsedFor().String()),
		attribute.String(fmt.Sprintf("Code.%v", index), in.GetCode()),
	)
	return span
}

func TraceVerify(span trace1.Span, in *npool.VerifyCodeRequest) trace1.Span {
	return traceVerify(span, in, 0)
}
