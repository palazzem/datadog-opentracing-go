package main

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"

	"net/http"
	"strconv"
)

func main() {
	// define the application router
	router := gin.Default()
	router.Use(TracingMiddleware())
	router.GET("/account/:id", getAccount)

	router.Run(":3000")
}

func getAccount(ctx *gin.Context) {
	id := ctx.Param("id")

	// add some context to the Trace (if the Span is available)
	if span := opentracing.SpanFromContext(ctx.Request.Context()); span != nil {
		span.SetTag("account_id", id)
	}

	ctx.String(http.StatusOK, "Account details for: %s", id)
}

func TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// create a "root" Span from the given Context
		// and finish it when the request ends
		span, ctx := opentracing.StartSpanFromContext(c, "gin.request")
		defer span.Finish()

		// propagate the trace in the Gin Context and process the request
		c.Request = c.Request.WithContext(ctx)
		c.Next()

		// add useful tags to your Trace
		span.SetTag("http.method", c.Request.Method)
		span.SetTag("http.status_code", strconv.Itoa(c.Writer.Status()))
		span.SetTag("http.url", c.Request.URL.Path)

		// add Datadog tag to distinguish stats for different endpoints
		span.SetTag("resource.name", c.HandlerName())
	}
}
