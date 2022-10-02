package main

import (
	"context"
	"github.com/My5z0n/SampleInstrumentationApp/APIGateway/model"
	"github.com/My5z0n/SampleInstrumentationApp/Utils"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
)

func initTracer() *sdktrace.TracerProvider {
	exporter, err := stdout.New(stdout.WithPrettyPrint())
	if err != nil {
		log.Fatal(err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("CustomerService"),
			)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}

func handler(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())

	var inputModel model.GetUserInfoModelInput

	if err := c.ShouldBindUri(&inputModel); err != nil {
		log.Printf("Unable to bind model: %s", err)
		return
	}
	span.SetAttributes(attribute.String("app.username", inputModel.User))
	span.SetAttributes(attribute.String("app.userdetails", inputModel.User))

	c.JSON(http.StatusOK, gin.H{
		"message": Utils.GetRandomString(10),
	})

}

var tracer = otel.Tracer("CustomerService")

func main() {
	tp := initTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	r := gin.Default()
	r.Use(otelgin.Middleware("CustomerService"))

	r.GET("/api/userinfo/:user", handler)

	r.Run(":8800") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
