package router

import (
    "fmt"
    "github.com/gin-gonic/gin"
)

func Initialize(externalCompanies []string, port int) {
    r := gin.Default()
    InitializeRoutes(r, externalCompanies)
    r.Run(fmt.Sprintf(":%d", port)) // Escuta na porta configurada
}
