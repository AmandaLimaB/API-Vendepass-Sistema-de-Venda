package router

import (
	"github.com/AmandaLimaB/API-Vendepass-Sistema-de-Venda/tree/main/server/handler"
	"github.com/gin-gonic/gin"
)

func initializeRoutes(r *gin.Engine) {
	r.GET("/list-flights", handler.ListFlightsHandler)
}
