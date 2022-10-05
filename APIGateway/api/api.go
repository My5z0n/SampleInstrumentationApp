package api

import (
	"encoding/json"
	"fmt"
	"github.com/My5z0n/SampleInstrumentationApp/APIGateway/model"
	"github.com/My5z0n/SampleInstrumentationApp/MessageHandler"
	"github.com/My5z0n/SampleInstrumentationApp/Utils"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
	"io"
	"log"
	"net/http"
)

var tracer = otel.Tracer("APIGateway")

func addAPIAtributes(c *gin.Context) {
	span := oteltrace.SpanFromContext(c.Request.Context())
	span.SetAttributes(attribute.String("firedog.user.agent", c.GetHeader("User-Agent")))

}

func GetUserInfo(c *gin.Context) {

	addAPIAtributes(c)
	span := oteltrace.SpanFromContext(c.Request.Context())

	var inputModel model.GetUserInfoModelInput
	err := c.ShouldBindUri(&inputModel)
	if err != nil {
		log.Printf("Unable to bind model: %s", err)
		return
	}

	span.SetAttributes(attribute.String("firedog.user.name", inputModel.User))

	targetURL := fmt.Sprintf("http://customerservice:8801/api/userinfo/%s", inputModel.User)

	res, err := otelhttp.Get(c.Request.Context(), targetURL)
	res.Request.Context()

	resBody, err := io.ReadAll(res.Body)

	var dat map[string]interface{}

	if err := json.Unmarshal(resBody, &dat); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, dat)

}
func CreateOrder(c *gin.Context) {
	addAPIAtributes(c)

	var orderModel model.CreateOrderModel
	err := c.ShouldBindJSON(&orderModel)
	if err != nil {
		log.Printf("Unable to bind model: %s", err)
		return
	}

	hdl := MessageHandler.GetMessageHandler(Utils.CreateOrderQueueName)
	hdl.SendMsg(map[string]any{
		"ProductName": orderModel.ProductName,
	}, c.Request.Context())

	c.JSON(http.StatusAccepted, nil)

}
func GetProductDetails(c *gin.Context) {
	addAPIAtributes(c)

	var productModel model.ProductDetailsModel
	err := c.ShouldBindJSON(&productModel)
	if err != nil {
		log.Printf("Unable to bind model: %s", err)
		return
	}

	targetURL := fmt.Sprintf("http://productservice:8802/api/getproductdetails/%s", productModel.ProductName)

	res, err := otelhttp.Get(c.Request.Context(), targetURL)
	res.Request.Context()

	resBody, err := io.ReadAll(res.Body)

	var dat map[string]interface{}

	if err := json.Unmarshal(resBody, &dat); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, dat)
}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, "Pong!")
}
