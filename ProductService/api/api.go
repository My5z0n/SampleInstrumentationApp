package api

import (
	"context"
	"github.com/My5z0n/SampleInstrumentationApp/MessageHandler"
	"github.com/My5z0n/SampleInstrumentationApp/ProductService/model"
	"github.com/My5z0n/SampleInstrumentationApp/Utils"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
)

var tracer = otel.Tracer("ProductService")

var MainConfig Utils.Config
var MsgHdlFactory MessageHandler.Factory

func SetSetting(cfg Utils.Config, msgHdlFactory MessageHandler.Factory) {
	MainConfig = cfg
	MsgHdlFactory = msgHdlFactory
}

func OrderDetails(span trace.Span, ctx context.Context, msg map[string]any, f MessageHandler.Factory) {
	defer span.End()
	hdlProductDetails := f.GetMessageHandler(Utils.ProcessConfirmedOrderQueueName)
	//TODO

	hdlProductDetails.SendMsg(msg, ctx)

}
func ProductDetails(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())

	var inputModel model.ProductDetailsModel
	err := c.ShouldBindUri(&inputModel)
	if err != nil {
		log.Printf("Unable to bind model: %s", err)
		return
	}
	span.SetAttributes(attribute.String("app.productname", inputModel.ProductName))
	c.JSON(http.StatusOK, nil)

	hdlProductDetails := MsgHdlFactory.GetMessageHandler(Utils.BigDataProductRequestQueueName)
	//TODO

	hdlProductDetails.SendMsg(map[string]any{
		"productname": inputModel.ProductName,
	}, c.Request.Context())
}
