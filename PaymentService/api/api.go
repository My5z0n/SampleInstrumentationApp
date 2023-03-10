package api

import (
	"context"
	"github.com/My5z0n/SampleInstrumentationApp/MessageHandler"
	"github.com/My5z0n/SampleInstrumentationApp/Utils"
	"go.opentelemetry.io/otel/trace"
)

func ProcessPaymentHandler(span trace.Span, ctx context.Context, msg map[string]any, f *MessageHandler.Factory) {
	defer span.End()
	hdlProductDetails := f.GetMessageHandler(Utils.ProcessReturnedPaymentQueueName)
	//TODO

	hdlProductDetails.SendMsg(msg, ctx)
}
