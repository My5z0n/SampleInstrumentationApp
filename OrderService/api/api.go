package api

import (
	"context"
	"fmt"
	"github.com/My5z0n/SampleInstrumentationApp/MessageHandler"
	"github.com/My5z0n/SampleInstrumentationApp/Utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"math/rand"
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
	hdlProcessPayment := f.GetMessageHandler(Utils.ProcessPaymentQueueName)

	productName := msg["ProductName"].(string)
	coupon := msg["Coupon"].(string)

	price := 0
	if coupon != "TEST-COUPON" {
		price = rand.Intn(1000) + 1
	}
	fmt.Printf("Price %d\n", price)

	span.SetAttributes(attribute.Int("app.priceForOrder", price))
	span.SetAttributes(attribute.String("app.orderedProductName", productName))
	span.SetAttributes(attribute.String("app.orderCoupon", coupon))

	if coupon == "TEST-COUPON" {
		hdlConfirmUserOrder := f.GetMessageHandler(Utils.ConfirmUserOrderQueueName)
		hdlConfirmUserOrder.SendMsg(msg, ctx)
	} else {
		hdlProcessPayment.SendMsg(msg, ctx)
	}

}
func ProcessReturnedPaymentHandler(span trace.Span, ctx context.Context, msg map[string]any, f *MessageHandler.Factory) {
	defer span.End()
	hdlProductDetails := f.GetMessageHandler(Utils.ConfirmUserOrderQueueName)
	//TODO

	hdlProductDetails.SendMsg(msg, ctx)
}
