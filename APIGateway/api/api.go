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
	"math/rand"
	"net/http"
)

var tracer = otel.Tracer("APIGateway")

var MainConfig Utils.Config
var MsgHdlFactory MessageHandler.Factory
var regions = []string{"eu-central-1", "eu-west-3", "us-east-1", "us-west-2", "eu-north-1"}

func SetSetting(cfg Utils.Config, msgHdlFactory MessageHandler.Factory) {
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

	targetURL := fmt.Sprintf("http://%s:8081/api/customer-userinfo/%s",
		MainConfig.URLMapper["customerservice"], inputModel.User)

	r := rand.Intn(len(regions))
	region := regions[r]
	span.SetAttributes(attribute.String("AWS.region", region))

	if c.GetHeader("experiment") == "true" {
		if region == "eu-central-1" {
			c.JSON(http.StatusInternalServerError, nil)
			return
		}
	}

	res, err := otelhttp.Get(c.Request.Context(), targetURL)
	if err != nil {
		log.Printf("Error during customerservice request: %v", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	if code := res.StatusCode; code == 200 {
		resBody, _ := io.ReadAll(res.Body)

		var dat map[string]interface{}

		if err := json.Unmarshal(resBody, &dat); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, dat)

	} else {
		c.JSON(http.StatusInternalServerError, nil)
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
	Utils.AddAPIAttributes(c)

	var productModel model.ProductDetailsModel
	err := c.ShouldBindJSON(&productModel)
	if err != nil {
		log.Printf("Unable to bind model: %s", err)
		return
	}

	targetURL := fmt.Sprintf("http://%s:8080/api/getproductdetails/%s", MainConfig.URLMapper["productservice"], productModel.ProductName)

	res, err := otelhttp.Get(c.Request.Context(), targetURL)
	if code := res.StatusCode; code == 200 {

		resBody, _ := io.ReadAll(res.Body)

		var dat map[string]interface{}

		if err := json.Unmarshal(resBody, &dat); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, dat)
	} else {
		c.JSON(http.StatusInternalServerError, nil)
	}
}

func Ping(c *gin.Context) {
	Utils.AddAPIAttributes(c)
	c.JSON(http.StatusOK, "Pong!")
}
