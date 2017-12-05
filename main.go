package main

import (
	datadog "github.com/DataDog/dd-trace-go/opentracing"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"

	"context"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func main() {
	// configure Datadog Tracer
	config := datadog.NewConfiguration()
	config.ServiceName = "api-intake"

	tracer, closer, _ := datadog.NewTracer(config)
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	// define the application router
	router := gin.Default()
	router.Use(TracingMiddleware())
	router.GET("/account/:id", getAccount)

	router.Run(":3000")
}

func getAccount(ctx *gin.Context) {
	var wg sync.WaitGroup

	id := ctx.Param("id")

	// add some context to the Trace (if the Span is available)
	span := opentracing.SpanFromContext(ctx.Request.Context())
	if span != nil {
		span.SetTag("account_id", id)
	}

	template := opentracing.StartSpan("template.rendering", opentracing.ChildOf(span.Context()))
	defer template.Finish()

	context := opentracing.ContextWithSpan(ctx.Request.Context(), template)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go dbCall(context, &wg)
	}

	wg.Wait()
	ctx.String(http.StatusOK, "Account details for: %s", id)
}

func dbCall(context context.Context, wg *sync.WaitGroup) {
	parent := opentracing.SpanFromContext(context)
	db := opentracing.StartSpan("db.query", opentracing.ChildOf(parent.Context()))
	db.SetTag("resource.name", "SELECT * FROM users;")
	db.SetTag("service.name", "db-cluster")
	defer db.Finish()

	time.Sleep(20 * time.Microsecond)
	wg.Done()
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
