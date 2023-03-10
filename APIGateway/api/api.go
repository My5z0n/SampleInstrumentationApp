package api

import (
	"fmt"
	"github.com/My5z0n/SampleInstrumentationApp/APIGateway/model"
	"github.com/My5z0n/SampleInstrumentationApp/MessageHandler"
	"github.com/My5z0n/SampleInstrumentationApp/Utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var tracer = otel.Tracer("APIGateway")

var MainConfig Utils.Config
var MsgHdlFactory *MessageHandler.Factory
var regions = []string{"eu-central-1", "eu-west-3", "us-east-1", "us-west-2", "eu-north-1"}

func SetSetting(cfg Utils.Config, msgHdlFactory *MessageHandler.Factory) {
	MainConfig = cfg
	MsgHdlFactory = msgHdlFactory
}

func GetUserInfo(c *gin.Context) {
	Utils.AddAPIAttributes(c)

	span := oteltrace.SpanFromContext(c.Request.Context())

	var inputModel model.GetUserInfoModelInput
	err := c.ShouldBindUri(&inputModel)
	if err != nil {
		log.Printf("Unable to bind model: %s", err)
		return
	}

	r := rand.Intn(len(regions))
	region := regions[r]
	span.SetAttributes(attribute.String("AWS.region", region))

	if c.GetHeader("experiment") == "true" {
		if region == "eu-central-1" {
			c.JSON(http.StatusInternalServerError, nil)
			return
		}
	}

	qid := fmt.Sprintf("%v", uuid.New())
	span.SetAttributes(attribute.String("GetProductDetailsQueue.qid", qid))
	msgResponse := MsgHdlFactory.SetWaitingResponse(qid, Utils.GetUserInfoResponseQueueName)
	hdl := MsgHdlFactory.GetMessageHandler(Utils.GetUserInfoQueueName)
	hdl.SendMsg(map[string]any{
		"UserName": inputModel.User,
		"QID":      qid,
	}, c.Request.Context())

	select {
	case msg := <-msgResponse:
		{
			close(msgResponse)
			defer msg.Span.End()
			c.JSON(http.StatusAccepted, "GetUser - OK Response")
		}
	case <-time.After(3 * time.Second):
		{
			c.JSON(http.StatusInternalServerError, "GetUser - Error Response")
		}
	}

}
func CreateOrder(c *gin.Context) {
	Utils.AddAPIAttributes(c)

	var orderModel model.CreateOrderModel
	err := c.ShouldBindJSON(&orderModel)
	if err != nil {
		log.Printf("Unable to bind model: %s", err)
		return
	}

	hdl := MsgHdlFactory.GetMessageHandler(Utils.CreateOrderQueueName)
	hdl.SendMsg(map[string]any{
		"ProductName": orderModel.ProductName,
	}, c.Request.Context())

	c.JSON(http.StatusAccepted, "OK")

}
func GetProductDetails(c *gin.Context) {
	span := oteltrace.SpanFromContext(c.Request.Context())
	Utils.AddAPIAttributes(c)
	var productModel model.ProductDetailsModel
	err := c.ShouldBindUri(&productModel)
	if err != nil {
		log.Printf("Unable to bind model: %s", err)
		return
	}

	qid := fmt.Sprintf("%v", uuid.New())
	span.SetAttributes(attribute.String("GetProductDetailsQueue.qid", qid))
	msgResponse := MsgHdlFactory.SetWaitingResponse(qid, Utils.GetProductDetailsResponseQueueName)
	hdl := MsgHdlFactory.GetMessageHandler(Utils.GetProductDetailsQueueName)
	hdl.SendMsg(map[string]any{
		"ProductName": productModel.ProductName,
		"QID":         qid,
	}, c.Request.Context())

	select {
	case msg := <-msgResponse:
		{
			close(msgResponse)
			defer msg.Span.End()
			c.JSON(http.StatusAccepted, "ProductDetails - OK Response")
		}
		//TODO Change timer after increase of max time in OtelCollector wait
	case <-time.After(3 * time.Second):
		{
			c.JSON(http.StatusAccepted, "ProductDetails - ERROR Response")
		}
	}

}

func Ping(c *gin.Context) {
	Utils.AddAPIAttributes(c)
	c.JSON(http.StatusOK, "Pong!")
}
