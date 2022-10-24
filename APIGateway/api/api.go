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

var URLMapper map[string]string

func addAPIAtributes(c *gin.Context) {
	span := oteltrace.SpanFromContext(c.Request.Context())
	span.SetAttributes(attribute.String("firedog.test1", c.GetHeader("User-Agent")))

}

func SetURLs() {

	URLMapper = make(map[string]string)
	switch Utils.GetEnv("NotDevelopment", "False") {
	default:
	case "False":
		URLMapper["customerservice"] = "localhost"
		URLMapper["productservice"] = "localhost"
	case "True":
		URLMapper["customerservice"] = "customerservice"
		URLMapper["productservice"] = "productservice"
	}

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

	span.SetAttributes(attribute.String("firedog.test2", inputModel.User))

	targetURL := fmt.Sprintf("http://%s:8801/api/userinfo/%s", URLMapper["customerservice"], inputModel.User)

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

	c.JSON(http.StatusAccepted, "OK")

}
func GetProductDetails(c *gin.Context) {
	addAPIAtributes(c)

	var productModel model.ProductDetailsModel
	err := c.ShouldBindJSON(&productModel)
	if err != nil {
		log.Printf("Unable to bind model: %s", err)
		return
	}

	targetURL := fmt.Sprintf("http://%s:8802/api/getproductdetails/%s", URLMapper["productservice"], productModel.ProductName)

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
	c.JSON(http.StatusOK, "Pong!")
}
