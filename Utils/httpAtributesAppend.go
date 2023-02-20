package Utils

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func AddAPIAttributes(c *gin.Context) {
	span := oteltrace.SpanFromContext(c.Request.Context())
	span.SetAttributes(attribute.String("api.test", "AddAPIAttributes"))
	//span.SetAttributes(attribute.String("user.agent", c.GetHeader("User-Agent")))
	//span.SetAttributes(attribute.String("ip", c.ClientIP()))

}
