package router

import (
	"github.com/AmandaLimaB/API-Vendepass-Sistema-de-Venda/tree/main/server/handler"
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(r *gin.Engine, externalCompanies []string) {
	r.GET("/list-flights", func(c *gin.Context) {
		handler.ListFlightsHandler(c, externalCompanies)
	})
	r.POST("/reserve-seat", func(c *gin.Context) {
		handler.ReserveSeatHandler(c, externalCompanies)
	})
	r.POST("/register-client", handler.RegisterClientHandler)
	r.POST("/cancel-reservation", func(c *gin.Context) {
		handler.CancelReservationHandler(c, externalCompanies)
	})
}
