package main

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.GET("/", home)
	router.POST("/choose/:val", choose)

	router.Run("0.0.0.0:8080") // Port must match EXPOSE command in Dockerfile
}

func home(ctx *gin.Context) {
	choices := make([]int, 4)
	for i := range choices {
		choices[i] = rand.IntN(99) + 1
	}

	ctx.JSON(http.StatusOK, choices)
}

func choose(ctx *gin.Context) {
	val := ctx.Param("val")
	intVal, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		ctx.String(http.StatusBadRequest, "You can't choose something that isn't a number!")
		return
	}

	ctx.String(http.StatusOK, fmt.Sprintf("You chose %d!", intVal))
}
