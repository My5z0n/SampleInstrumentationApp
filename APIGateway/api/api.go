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
var azones = []string{"us-east-1a", "us-east-1b", "us-east-1c", "us-east-1d", "us-east-1e"}

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

	r := rand.Intn(len(azones))
	region := azones[r]
	span.SetAttributes(attribute.String("AWS.azone", region))

	qid := fmt.Sprintf("%v", uuid.New())
	span.SetAttributes(attribute.String("GetProductDetailsQueue.qid", qid))
	msgResponse := MsgHdlFactory.SetWaitingResponse(qid, Utils.GetUserInfoResponseQueueName)
	hdl := MsgHdlFactory.GetMessageHandler(Utils.GetUserInfoQueueName)

	if c.GetHeader("experiment") == "true" {
		if region == "us-east-1b" {
			hdl.MockSendMsg(map[string]any{
				"UserName": inputModel.User,
				"QID":      qid,
			}, c.Request.Context())
		} else {
			hdl.SendMsg(map[string]any{
				"UserName": inputModel.User,
				"QID":      qid,
			}, c.Request.Context())
		}
	} else {
		hdl.SendMsg(map[string]any{
			"UserName": inputModel.User,
			"QID":      qid,
		}, c.Request.Context())
	}

	select {
	case msg := <-msgResponse:
		{
			close(msgResponse)
			defer msg.Span.End()
			c.JSON(http.StatusOK, "GetUser - OK Response")
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
		"Coupon":      orderModel.Coupon,
	}, c.Request.Context())

	c.JSON(http.StatusAccepted, "OK")

}
func GetProductDetails(c *gin.Context) {
	span := oteltrace.SpanFromContext(c.Request.Context())
	Utils.AddAPIAttributes(c)
	var productModel = c.Param("productname")

	qid := fmt.Sprintf("%v", uuid.New())
	span.SetAttributes(attribute.String("GetProductDetailsQueue.qid", qid))
	msgResponse := MsgHdlFactory.SetWaitingResponse(qid, Utils.GetProductDetailsResponseQueueName)
	hdl := MsgHdlFactory.GetMessageHandler(Utils.GetProductDetailsQueueName)
	hdl.SendMsg(map[string]any{
		"ProductName": productModel,
		"QID":         qid,
	}, c.Request.Context())

	select {
	case msg := <-msgResponse:
		{
			close(msgResponse)
			defer msg.Span.End()
			c.JSON(http.StatusOK, "ProductDetails - OK Response")
		}
		//TODO Change timer after increase of max time in OtelCollector wait
	case <-time.After(3 * time.Second):
		{
			c.JSON(http.StatusInternalServerError, "ProductDetails - ERROR Response")
		}
	}

}

func Ping(c *gin.Context) {
	Utils.AddAPIAttributes(c)
	c.JSON(http.StatusOK, "Pong!")
}
