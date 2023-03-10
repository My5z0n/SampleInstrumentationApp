package MessageHandler

import (
	"context"
	"github.com/My5z0n/SampleInstrumentationApp/Utils"
	"go.opentelemetry.io/otel/trace"
	"log"
	"time"
)

func ReceiveIDMsg(span trace.Span, ctx context.Context, msg map[string]any, f *Factory, queueName string) {

	msgID := msg["QID"].(string)

	for true {
		if waitChan, okResponse := f.GetWaitingResponseChan(msgID, queueName); okResponse {
			waitChan <- MessageResponse{
				Span: span,
				Ctx:  ctx,
				Msg:  msg,
			}
			return
		} else {
			log.Printf("Unable to process message, QueueName: %v, ID: %v", queueName, msgID)
			time.Sleep(time.Second * 5)
		}
	}

}

func MsgRcv(handler func(span trace.Span, ctx context.Context, msg map[string]any, f *Factory), queueName string, f *Factory, modeID Utils.ReceiveMessageMode) {
	msgHandler := f.GetMessageHandler(queueName)
	inputChan := msgHandler.RegisterConsumer()

	for msg := range inputChan {
		if span, ok := msg["OTELSPAN"].(trace.Span); ok {
			if c, ok := msg["CONTEXT"].(context.Context); ok {
				if msgBody, ok := msg["msg"].(map[string]any); ok {
					if modeID == Utils.All {
						go handler(span, c, msgBody, f)
					} else {
						go ReceiveIDMsg(span, c, msgBody, f, queueName)
					}
				}
			}

		}
	}
}
