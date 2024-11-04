package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/AmandaLimaB/API-Vendepass-Sistema-de-Venda/tree/main/server/models"
	"github.com/AmandaLimaB/API-Vendepass-Sistema-de-Venda/tree/main/server/repository"
)

// ListFlightsHandler lista os voos locais e os das outras companhias
func ListFlightsHandler(c *gin.Context, externalCompanies []string) {
	var allFlights []models.Flight
	fmt.Println(externalCompanies)
	// Carregar voos locais primeiro
	localFlights, err := repository.LoadFlights()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao carregar voos locais"})
		return
	}
	allFlights = append(allFlights, localFlights...)

	// Verificar se a requisição é externa para incluir voos das outras companhias
	if c.GetHeader("X-Internal-Request") != "true" {
		for _, url := range externalCompanies {
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/list-flights", url), nil)
			if err != nil {
				fmt.Printf("Erro ao criar requisição para %s: %v\n", url, err)
				continue
			}
			req.Header.Set("X-Internal-Request", "true") // Define o cabeçalho para identificar a requisição como interna

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Printf("Erro ao requisitar voos da companhia em %s: %v\n", url, err)
				continue
			}
			defer resp.Body.Close()

			var externalFlights []models.Flight
			if err := json.NewDecoder(resp.Body).Decode(&externalFlights); err != nil {
				fmt.Printf("Erro ao decodificar resposta da companhia em %s: %v\n", url, err)
				continue
			}
			allFlights = append(allFlights, externalFlights...)
		}
	}

	// Retornar todos os voos (locais + externos)
	c.JSON(http.StatusOK, allFlights)
}
