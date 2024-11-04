package handler

import (
	"net/http"
	"github.com/AmandaLimaB/API-Vendepass-Sistema-de-Venda/tree/main/server/models"
	"github.com/AmandaLimaB/API-Vendepass-Sistema-de-Venda/tree/main/server/repository"
	"github.com/gin-gonic/gin"
)

type ClientRegistrationRequest struct {
	CPF   string `json:"cpf" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func RegisterClientHandler(c *gin.Context) {
	var req ClientRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CPF e senha são obrigatórios"})
		return
	}

	client := models.Client{
		CPF:      req.CPF,
		Password: req.Password,
	}

	if err := repository.SaveClient(client); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao cadastrar cliente"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cliente cadastrado com sucesso"})
}