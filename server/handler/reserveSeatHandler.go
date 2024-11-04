package handler

import (
	"net/http"

	"github.com/AmandaLimaB/API-Vendepass-Sistema-de-Venda/tree/main/server/models"
	"github.com/AmandaLimaB/API-Vendepass-Sistema-de-Venda/tree/main/server/repository"
	"github.com/gin-gonic/gin"
)

func ReserveSeatHandler(c *gin.Context, externalCompanies []string) {
	var req models.ReserveSeatRequest

	// Bind JSON da requisição
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados de entrada inválidos"})
		return
	}

	// Realiza a reserva localmente ou em outra companhia
	err := repository.ReserveSeat(req.FlightId, req.SeatID, req.CPF, externalCompanies)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"message": "Não foi possível reservar o assento"})
		return
	}

	// Retorna sucesso
	c.JSON(http.StatusOK, gin.H{"message": "Assento reservado com sucesso"})
}