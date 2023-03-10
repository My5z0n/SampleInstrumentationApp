package MessageHandler

import (
	"context"
	"fmt"
	"github.com/My5z0n/SampleInstrumentationApp/Utils"
	"go.opentelemetry.io/otel/trace"
	"log"
	"sync"
)

type MessageResponse struct {
	Span trace.Span
	Ctx  context.Context
	Msg  map[string]any
}

type Factory struct {
	handlerMapper        map[string]MessageHandler
	receiveIDMapper      map[string]map[string]chan MessageResponse
	config               Utils.Config
	messageHandlerMutex  sync.Mutex
	waitingResponseMutex sync.Mutex
}

func GetFactory(cfg Utils.Config) Factory {
	return Factory{
		handlerMapper:   make(map[string]MessageHandler),
		receiveIDMapper: make(map[string]map[string]chan MessageResponse),
		config:          cfg,
	}
}

func (f *Factory) GetMessageHandler(queueName string) MessageHandler {
	f.messageHandlerMutex.Lock()
	defer f.messageHandlerMutex.Unlock()
	if val, ok := f.handlerMapper[queueName]; ok {
		return val
	} else {
		connectionString := fmt.Sprintf("amqp://%s:%s@%s:5672",
			f.config.ServiceCredential["rabbitmq"].Username,
			f.config.ServiceCredential["rabbitmq"].Password,
			f.config.URLMapper["rabbitmq"])
		log.Printf("Try connection to: %v", connectionString)

		tmp := MessageHandler{}
		tmp.CreateConnection(queueName, connectionString)
		f.handlerMapper[queueName] = tmp
		return tmp
	}

}
func (f *Factory) SetWaitingResponse(id string, queueName string) chan MessageResponse {
	f.waitingResponseMutex.Lock()
	defer f.waitingResponseMutex.Unlock()

	var response = make(chan MessageResponse)

	if _, ok := f.receiveIDMapper[queueName]; !ok {
		f.receiveIDMapper[queueName] = make(map[string]chan MessageResponse)
	}
	chanMap := f.receiveIDMapper[queueName]

	if _, ok := chanMap[id]; ok {
		log.Printf("Chan queue %v with response id %v already exists", queueName, id)
		return nil
	}
	chanMap[id] = response

	return response
}
func (f *Factory) GetWaitingResponseChan(id string, queueName string) (chan MessageResponse, bool) {
	f.waitingResponseMutex.Lock()
	defer f.waitingResponseMutex.Unlock()

	var responseChan chan MessageResponse
	var status = false
	if mapper, ok := f.receiveIDMapper[queueName]; ok {
		responseChan, status = mapper[id]
		delete(mapper, id)
	}

	return responseChan, status
}
