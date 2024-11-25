package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/mhs294/mulhall/db"
	"github.com/mhs294/mulhall/types"
	"github.com/mhs294/mulhall/views"
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

	// Static content (CSS/JS)
	router.Static("/static", "./static")

	// Middleware
	router.Use(corsMiddleware)

	router.GET("/", indexPage)
	router.GET("/teams", getTeams)
	router.POST("/teams/:id", chooseTeam)

	router.Run("0.0.0.0:8080") // Port must match EXPOSE command in Dockerfile
}

func corsMiddleware(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Credentials", "true")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	ctx.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	if ctx.Request.Method == http.MethodOptions {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}

	ctx.Next()
}
func render(ctx *gin.Context, status int, template templ.Component) error {
	ctx.Status(status)
	return template.Render(ctx.Request.Context(), ctx.Writer)
}

func indexPage(ctx *gin.Context) {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	teams := teamRepo.GetAllTeams()
	render(ctx, http.StatusOK, views.Index(teams))
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
