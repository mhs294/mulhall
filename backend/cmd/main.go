package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.SetTrustedProxies(nil)

	// CORS Middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	router.GET("/teams", getTeams)
	router.POST("/teams/:id", chooseTeam)

	router.Run("0.0.0.0:8080") // Port must match EXPOSE command in Dockerfile
}

type Team struct {
	ID        int    `json:"id"`
	Shorthand string `json:"shorthand"`
	Location  string `json:"location"`
	Name      string `json:"name"`
}

// TODO - host in database
var teams = map[int]Team{
	1: {
		ID:        1,
		Shorthand: "ARI",
		Location:  "Arizona",
		Name:      "Cardinals",
	},
	2: {
		ID:        2,
		Shorthand: "GB",
		Location:  "Green Bay",
		Name:      "Packers",
	},
}

func getTeams(ctx *gin.Context) {
	t := make([]Team, len(teams))
	for i := 1; i <= len(teams); i++ {
		t = append(t, teams[i])
	}
	ctx.JSON(http.StatusOK, t)
}

func chooseTeam(ctx *gin.Context) {
	id := ctx.Param("id")
	teamID, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.String(http.StatusBadRequest, "need an int value for team ID")
		return
	}

	team := teams[int(teamID)]
	ctx.String(http.StatusOK, fmt.Sprintf("You chose the %s %s!", team.Location, team.Name))
}
