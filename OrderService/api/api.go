package api

import (
	"context"
	"github.com/My5z0n/SampleInstrumentationApp/MessageHandler"
	"github.com/My5z0n/SampleInstrumentationApp/Utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("OrderService")

func CreateOrderHandler(span trace.Span, ctx context.Context, msg map[string]any, f *MessageHandler.Factory) {
	defer span.End()
	hdlProductDetails := f.GetMessageHandler(Utils.ConfirmProductDetailsQueueName)
	//TODO

	hdlProductDetails.SendMsg(msg, ctx)
}
func ProcessOrderHandler(span trace.Span, ctx context.Context, msg map[string]any, f *MessageHandler.Factory) {
	defer span.End()
	hdlProductDetails := f.GetMessageHandler(Utils.ProcessPaymentQueueName)
	//TODO

	hdlProductDetails.SendMsg(msg, ctx)
}
func ProcessReturnedPaymentHandler(span trace.Span, ctx context.Context, msg map[string]any, f *MessageHandler.Factory) {
	defer span.End()
	hdlProductDetails := f.GetMessageHandler(Utils.ConfirmUserOrderQueueName)
	//TODO

	hdlProductDetails.SendMsg(msg, ctx)
}
