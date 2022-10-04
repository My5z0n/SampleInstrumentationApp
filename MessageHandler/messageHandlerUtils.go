package MessageHandler

import (
	"context"
	"go.opentelemetry.io/otel/trace"
)

func MsgRcv(handler func(trace.Span, context.Context, map[string]any), queueName string) {
	msgHandler := MessageHandler{}
	msgHandler.Create(queueName)
	inputChan := msgHandler.RegisterConsumer()

	for msg := range inputChan {
		if span, ok := msg["OTELSPAN"].(trace.Span); ok {
			if c, ok := msg["CONTEXT"].(context.Context); ok {
				if msgBody, ok := msg["msg"].(map[string]any); ok {
					go handler(span, c, msgBody)
				}
			}

		}
	}
}
