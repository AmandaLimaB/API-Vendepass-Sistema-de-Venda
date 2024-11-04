package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/AmandaLimaB/API-Vendepass-Sistema-de-Venda/tree/main/server/models"
	"github.com/AmandaLimaB/API-Vendepass-Sistema-de-Venda/tree/main/server/router"
)

func loadConfig() models.Config {
    file, err := os.Open("data/config.json")
    if err != nil {
        log.Fatalf("Erro ao carregar configuração: %v", err)
    }
    defer file.Close()

    var config models.Config
    if err := json.NewDecoder(file).Decode(&config); err != nil {
        log.Fatalf("Erro ao decodificar configuração: %v", err)
    }
    return config
}

func main() {
    config := loadConfig()
    fmt.Printf("Iniciando %s na porta %d\n", config.CompanyName, config.Port)
    router.Initialize(config.ExternalCompanies, config.Port)
}
