package Utils

import (
	"fmt"
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

const GetProductDetailsQueueName = "GetProductDetailsQ"
const GetProductDetailsResponseQueueName = "GetProductDetailsQR"

const GetUserInfoQueueName = "GetUserInfoQ"
const GetUserInfoResponseQueueName = "GetUserInfoQR"

type ReceiveMessageMode int64

const (
	All ReceiveMessageMode = iota
	IDOnly
)

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
