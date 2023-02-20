package api

import (
	"context"
	"github.com/My5z0n/SampleInstrumentationApp/CustomerService/model"
	"github.com/My5z0n/SampleInstrumentationApp/MessageHandler"
	"github.com/My5z0n/SampleInstrumentationApp/Utils"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
)

func GetUserHandler(c *gin.Context) {
	Utils.AddAPIAttributes(c)

	span := trace.SpanFromContext(c.Request.Context())

	var inputModel model.GetUserInfoModelInput

	if err := c.ShouldBindUri(&inputModel); err != nil {
		log.Printf("Unable to bind model: %s", err)
		return
	}

	span.SetAttributes(attribute.String("app.username", inputModel.User))

	c.JSON(http.StatusOK, gin.H{
		"message": Utils.GetRandomString(10),
	})
}
func ConfirmUserOrder(span trace.Span, ctx context.Context, msg map[string]any, f MessageHandler.Factory) {
	defer span.End()
	log.Printf("ORDER ENDED")

}
