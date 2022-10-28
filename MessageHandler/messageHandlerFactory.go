package MessageHandler

import (
	"fmt"
	"github.com/My5z0n/SampleInstrumentationApp/Utils"
	"sync"
)

type Factory struct {
	handlerMapper map[string]MessageHandler
	config        Utils.Config
	mu            sync.Mutex
}

func GetFactory(cfg Utils.Config) Factory {
	return Factory{
		handlerMapper: make(map[string]MessageHandler),
		config:        cfg,
	}
}

func (f *Factory) GetMessageHandler(queueName string) MessageHandler {
	f.mu.Lock()
	defer f.mu.Unlock()
	if val, ok := f.handlerMapper[queueName]; ok {
		return val
	} else {
		connectionString := fmt.Sprintf("amqp://%s:%s@%s:5672/",
			f.config.ServiceCredential["rabbitmq"].Username,
			f.config.ServiceCredential["rabbitmq"].Password,
			f.config.URLMapper["rabbitmq"])

		tmp := MessageHandler{}
		tmp.CreateConnection(queueName, connectionString)
		f.handlerMapper[queueName] = tmp
		return tmp
	}

}
