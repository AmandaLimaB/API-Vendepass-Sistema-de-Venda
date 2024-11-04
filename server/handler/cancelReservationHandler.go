package handler

import (
	"net/http"

	"github.com/AmandaLimaB/API-Vendepass-Sistema-de-Venda/tree/main/server/models"	
	"github.com/AmandaLimaB/API-Vendepass-Sistema-de-Venda/tree/main/server/repository"
	"github.com/gin-gonic/gin"
)

// CancelReservationHandler cancela uma reserva de assento
func CancelReservationHandler(c *gin.Context, externalCompanies []string) {
	var req models.CancelReservationRequest

	// Bind JSON da requisição
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados de entrada inválidos"})
		return
	}

	// Realiza o cancelamento localmente ou em outra companhia
	err := repository.CancelReservation(req.FlightId, req.SeatID, req.CPF, externalCompanies)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"message": "Cancelamento não realizado"})
		return
	}

	// Retorna sucesso
	c.JSON(http.StatusOK, gin.H{"message": "Reserva cancelada com sucesso"})
}
