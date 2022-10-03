package Utils

import "github.com/My5z0n/SampleInstrumentationApp/MessageHandler"

var factoryMap map[string]MessageHandler.MessageHandler

func GetMessageHandler(handlerName string) MessageHandler.MessageHandler {

	if val, ok := factoryMap[handlerName]; ok {
		return val
	} else {

		tmp := MessageHandler.MessageHandler{}
		tmp.Create(handlerName)
		factoryMap[handlerName] = tmp
		return tmp
	}

}
