package main

import (
	"github.com/gin-gonic/gin"

	"net/http"
)

func main() {
	// define the application router
	router := gin.Default()
	router.GET("/account/:id", getAccount)

	router.Run(":3000")
}

func getAccount(ctx *gin.Context) {
	id := ctx.Param("id")
	ctx.String(http.StatusOK, "Account details for: %s", id)
}
