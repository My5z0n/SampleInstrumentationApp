package api

import (
	"context"
	"fmt"
	"github.com/My5z0n/SampleInstrumentationApp/MessageHandler"
	"github.com/My5z0n/SampleInstrumentationApp/ProductService/model"
	"github.com/My5z0n/SampleInstrumentationApp/Utils"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
	"strconv"
)

var tracer = otel.Tracer("ProductService")

var MainConfig Utils.Config
var MsgHdlFactory *MessageHandler.Factory

func SetSetting(cfg Utils.Config, msgHdlFactory *MessageHandler.Factory) {
	MainConfig = cfg
	MsgHdlFactory = msgHdlFactory
}

func OrderDetails(span trace.Span, ctx context.Context, msg map[string]any, f *MessageHandler.Factory) {
	defer span.End()
	hdlProductDetails := f.GetMessageHandler(Utils.ProcessConfirmedOrderQueueName)
	//TODO

	hdlProductDetails.SendMsg(msg, ctx)

}
func ProductDetails(span trace.Span, ctx context.Context, msg map[string]any, f *MessageHandler.Factory) {
	defer span.End()

	hdlProductDetails := f.GetMessageHandler(Utils.GetProductDetailsResponseQueueName)
	//TODO

	productName := msg["ProductName"].(string)
	productID, _ := strconv.Atoi(productName[7:])

	span.SetAttributes(attribute.String("app.productName", productName))
	if productID < 30 {
		//PL warehouse
		span.SetAttributes(attribute.String("app.productWarehouse", "PL"))
	} else {
		span.SetAttributes(attribute.String("app.productWarehouse", "DE"))
	}

	if productID < 20 {
		fmt.Print("Product ERROR")
		return
	}
	hdlProductDetails.SendMsg(msg, ctx)

}
func OldProductDetails(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())

	var inputModel model.ProductDetailsModel
	err := c.ShouldBindUri(&inputModel)
	if err != nil {
		log.Printf("Unable to bind model: %s", err)
		return
	}
	span.SetAttributes(attribute.String("app.productname", inputModel.ProductName))

	c.JSON(http.StatusOK, nil)

}
