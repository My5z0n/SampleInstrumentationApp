package api

import (
	"context"
	"github.com/My5z0n/SampleInstrumentationApp/MessageHandler"
	"github.com/My5z0n/SampleInstrumentationApp/Utils"
	"go.opentelemetry.io/otel/trace"
	"log"
)

func GetUserInfo(span trace.Span, ctx context.Context, msg map[string]any, f *MessageHandler.Factory) {
	defer span.End()

	tmp := msg["UserName"].(string)

	if tmp == "Jeff" {
		return
	}

	hdlProductDetails := f.GetMessageHandler(Utils.GetUserInfoResponseQueueName)
	//TODO

	hdlProductDetails.SendMsg(msg, ctx)

}
func ConfirmUserOrder(span trace.Span, ctx context.Context, msg map[string]any, f *MessageHandler.Factory) {
	defer span.End()
	log.Printf("ORDER ENDED")

}
