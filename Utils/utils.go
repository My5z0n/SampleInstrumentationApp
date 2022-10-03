package Utils

import (
	"context"
	"fmt"
	"github.com/My5z0n/SampleInstrumentationApp/MessageHandler"
	"go.opentelemetry.io/otel/trace"
	"log"
	"math/rand"
	"os"
	"time"
)

const TraceparentHeader = "traceparent"

// Queue Name
const CreateOrderQueueName = "CreateOrderQ"
const ConfirmProductDetailsQueueName = "ConfirmProductDetailsQ"
const ProcessConfirmedOrderQueueName = "ProcessConfirmedOrderQ"
const ProcessPaymentQueueName = "ProcessPaymentQ"
const ProcessReturnedPaymentQueueName = "ProcessReturnedPaymentQ"
const ConfirmUserOrderQueueName = "ConfirmUserOrderQ"
const BigDataProductRequestQueueName = "BigDataProductRequestQ"

func GetEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func GetRandomString(lenstr int) string {
	rand.Seed(time.Now().UnixNano())
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	randomString := ""

	for i := 0; i < lenstr; i++ {
		randomchar := string(charset[rand.Intn(len(charset))])
		randomString = fmt.Sprintf("%s%s", randomString, randomchar)
	}
	return randomString
}
func MsgRcv(handler func(trace.Span, context.Context, map[string]any), queueName string) {
	msgHandler := MessageHandler.MessageHandler{}
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
