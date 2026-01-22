package routes

import (
    "api-turnos/internal/handlers"
    "api-turnos/internal/repository"
    "api-turnos/internal/services"
    "api-turnos/internal/websockets"
    "database/sql"

    "github.com/gin-gonic/gin"
)

func SetupRouter(db *sql.DB, hub *websockets.Hub) *gin.Engine {
    r := gin.Default()

    // Configuración de CORS (Importante para que las pantallas no den error)
    r.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    })

    // --- Inyección de Dependencias ---
    // 1. Repository
    pantallaRepo := repositories.NewPantallaRepository(db)
    
    // 2. Service
    pantallaService := services.PantallaService(pantallaRepo)
    
    // 3. Handler (Recibe Service y Hub)
    pantallaHandler := handlers.NewPantallaHandle(pantallaService, hub)

    // --- Definición de Rutas ---
    api := r.Group("/api")
    {
        api.POST("/turnos/llamar", pantallaHandler.Receive)
		api.GET("/turnos/estado", pantallaHandler.GetEstadoActual)
    }

    r.GET("/ws", pantallaHandler.ConnectWS)

    return r
}