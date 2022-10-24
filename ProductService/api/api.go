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

func OrderDetails(span trace.Span, ctx context.Context, msg map[string]any) {
	defer span.End()
	hdlProductDetails := MessageHandler.GetMessageHandler(Utils.ProcessConfirmedOrderQueueName)
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

	hdlProductDetails := MessageHandler.GetMessageHandler(Utils.BigDataProductRequestQueueName)
	//TODO

	hdlProductDetails.SendMsg(map[string]any{
		"productname": inputModel.ProductName,
	}, c.Request.Context())
}
