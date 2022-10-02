package api

import (
	"context"
	"github.com/My5z0n/SampleInstrumentationApp/APIGateway/model"
	"github.com/My5z0n/SampleInstrumentationApp/MessageHandler"
	"github.com/My5z0n/SampleInstrumentationApp/Utils"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
)

var tracer = otel.Tracer("fiber-server")

func GetUserInfo(c *gin.Context) {

	hdl := MessageHandler.MessageHandler{}
	hdl.Create(Utils.CreateOrderQueueName)

	var inputModel model.GetUserInfoModelInput
	err := c.ShouldBindUri(&inputModel)
	if err != nil {
		log.Printf("Unable to bind model: %s", err)
		return
	}

	//res, err := otelhttp.Get(c.Request.Context(), "http://localhost:8800/ping")
	//res.Request.Context()

	//resBody, err := io.ReadAll(res.Body)

	//var dat map[string]interface{}

	//if err := json.Unmarshal(resBody, &dat); err != nil {
	//	panic(err)
	//}

	//getUser(c.Request.Context(), "xD")
	//tc := propagation.TraceContext{}
	//x := propagation.MapCarrier{}
	//tc.Inject(c.Request.Context(), x)
	//fmt.Println("MESYCZ: ", x)

	// Register the TraceContext propagator globally.
	//otel.SetTextMapPropagator(tc)
	hdl.SendMsg(map[string]any{
		"message": inputModel.User,
	}, c.Request.Context())

	//newMSG := map[string]any{
	//	"message": dat["message"],
	//}
	c.JSON(http.StatusOK, inputModel.User)

}
func getUser(ctx context.Context, id string) string {
	_, span := tracer.Start(ctx, "getUser", oteltrace.WithAttributes(attribute.String("id", id)))
	defer span.End()
	if id == "123" {
		return "otelfiber tester"
	}
	return "unknown"
}
