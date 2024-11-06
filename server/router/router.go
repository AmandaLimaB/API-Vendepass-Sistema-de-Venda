package router

import (
    "fmt"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

func Initialize(externalCompanies []string, port int) {
    r := gin.Default()
    r.Use(cors.Default())
    InitializeRoutes(r, externalCompanies)
    r.Run(fmt.Sprintf(":%d", port)) // Escuta na porta configurada
}
