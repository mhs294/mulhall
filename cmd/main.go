package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mhs294/mulhall/db"
	"github.com/mhs294/mulhall/types"
)

var mongoDBConnStr string
var teamRepo *db.TeamRepository

func init() {
	log.Println("init started")

	mongoDBConnStr = os.Getenv("MULHALL_DB_CONN_STR")
	var err error
	teamRepo, err = db.NewTeamRepository(mongoDBConnStr)
	if err != nil {
		panic(err)
	}

	log.Println("init complete")
}

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

func getTeams(ctx *gin.Context) {
	teams := teamRepo.GetAllTeams()
	ctx.JSON(http.StatusOK, teams)
}

func chooseTeam(ctx *gin.Context) {
	id := ctx.Param("id")
	team := teamRepo.GetTeam(id)
	if team == (types.Team{}) {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.String(http.StatusOK, fmt.Sprintf("You chose the %s %s!", team.Location, team.Name))
}
