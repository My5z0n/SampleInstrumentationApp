package MessageHandler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/My5z0n/SampleInstrumentationApp/Utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"log"
	"time"
)

type MessageHandler struct {
	queue            amqp.Queue
	chanel           *amqp.Channel
	connection       *amqp.Connection
	QueueName        string
	ConnectionString string
	ExternalChan     chan map[string]any
	internalChan     <-chan amqp.Delivery
}

var tr = otel.Tracer("MessageHandler")

func (rcv *MessageHandler) CreateConnection(queueName string, con string) {

	if con == "" {
		log.Panicf("Empty connection string in Rabbitmq Message Handler. Unable to connect.")
	} else {
		rcv.ConnectionString = con
	}

	rcv.QueueName = queueName

	maxPoolTime := time.Now().Add(time.Second * 120)
	for true {
		conn, err := amqp.Dial(rcv.ConnectionString)
		if err == nil {
			rcv.connection = conn
			break
		} else if !time.Now().After(maxPoolTime) {
			log.Println("Failed to connect to RabbitMQ retrying...")
			time.Sleep(time.Second * 5)
			continue
		} else {
			Utils.FailOnError(err, "Failed to connect to RabbitMQ")
		}
	}
	ch, err := rcv.connection.Channel()
	Utils.FailOnError(err, "Failed to open a channel")
	rcv.chanel = ch

	q, err := ch.QueueDeclare(
		rcv.QueueName, // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	Utils.FailOnError(err, "Failed to declare a queue")
	rcv.queue = q

}

func (rcv *MessageHandler) SendMsg(message map[string]any, ctx context.Context) {
	newCtx, sp := tr.Start(ctx, fmt.Sprintf("%s send", rcv.queue.Name))
	defer sp.End()

	rcv.handleMsgSpanAttributes(sp)

	tc := propagation.TraceContext{}
	carrier := propagation.MapCarrier{}
	tc.Inject(newCtx, carrier)

	headers := map[string]any{
		Utils.TraceparentHeader: carrier[Utils.TraceparentHeader],
	}
	body, err := json.Marshal(message)

	//Publish
	err = rcv.publish(body, headers)

	//Try to retry connection if error
	if err != nil {
		rcv.CreateConnection(rcv.QueueName, rcv.ConnectionString)
		err = rcv.publish(body, headers)
		Utils.FailOnError(err, "Failed to publish a message")
	}

}

func (rcv *MessageHandler) MockSendMsg(message map[string]any, ctx context.Context) {
	newCtx, sp := tr.Start(ctx, fmt.Sprintf("%s send", rcv.queue.Name))
	defer sp.End()

	rcv.handleMsgSpanAttributes(sp)

	tc := propagation.TraceContext{}
	carrier := propagation.MapCarrier{}
	tc.Inject(newCtx, carrier)

}

func (rcv *MessageHandler) RegisterConsumer() <-chan map[string]any {
	msgs, err := rcv.chanel.Consume(
		rcv.queue.Name, // queue
		"",             // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	Utils.FailOnError(err, "Failed to register a consumer")

	rcv.internalChan = msgs

	var retch chan map[string]any

	if rcv.ExternalChan == nil {
		retch = make(chan map[string]any)
		rcv.ExternalChan = retch

		go func() {
			var ret map[string]any
			for {
				log.Printf("Listening at message chanel %s", rcv.QueueName)
				for rawMsg := range rcv.internalChan {

					if err := json.Unmarshal(rawMsg.Body, &ret); err != nil {
						return
					}
					tc := propagation.TraceContext{}

					tmpCarrier := propagation.MapCarrier{
						Utils.TraceparentHeader: rawMsg.Headers[Utils.TraceparentHeader].(string),
					}

					ctx := tc.Extract(context.Background(), tmpCarrier)
					c, sp := tr.Start(ctx, fmt.Sprintf("%s receive", rcv.queue.Name))
					rcv.handleMsgSpanAttributes(sp)

					retch <- map[string]any{
						"msg":      ret,
						"OTELSPAN": sp,
						"CONTEXT":  c,
					}
				}

				fmt.Println("Message chanel unexpectedly closed. Reconnecting...")
				rcv.reconnectConsumer()
			}

		}()

	} else {
		retch = rcv.ExternalChan
	}

	return retch
}

func (rcv *MessageHandler) publish(body []byte, headers map[string]any) error {
	return rcv.chanel.PublishWithContext(context.Background(),
		"",             // exchange
		rcv.queue.Name, // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Headers:     headers,
		})
}

func (rcv *MessageHandler) reconnectConsumer() {
	rcv.CreateConnection(rcv.QueueName, rcv.ConnectionString)
	_ = rcv.RegisterConsumer()
}
func (rcv *MessageHandler) handleMsgSpanAttributes(sp trace.Span) {

	sp.SetAttributes(attribute.String("messaging.system", "rabbitmq"))
	sp.SetAttributes(attribute.String("messaging.destination", rcv.queue.Name))
	sp.SetAttributes(attribute.String("messaging.destination_kind", "queue"))
	sp.SetAttributes(attribute.String("messaging.messaging.protocol", "AMQP"))
	sp.SetAttributes(attribute.String("messaging.protocol_version", "0.9.1"))
	sp.SetAttributes(attribute.String("messaging.url", rcv.ConnectionString))

}
