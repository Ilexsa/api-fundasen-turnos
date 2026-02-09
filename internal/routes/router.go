package routes

import (
	"api-turnos/internal/handlers"
	repositories "api-turnos/internal/repository"
	"api-turnos/internal/services"
	"api-turnos/internal/websockets"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetupRouter(db *sql.DB, hub *websockets.Hub, pacientesHub *websockets.Hub) *gin.Engine {
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
	pantallaRepo := repositories.NewPantallaRepository(db)
	pantallaService := services.PantallaService(pantallaRepo)
	pantallaHandler := handlers.NewPantallaHandle(pantallaService, hub)
	pacienteRepo := repositories.NewPacienteRepository(db)
	pacienteService := services.NewPacienteService(pacienteRepo)
	pacienteHandler := handlers.NewPacienteHandler(pacienteService, pacientesHub)

	// --- Definición de Rutas ---
	api := r.Group("/api")
	{
		// Rutas de Pantallas
		api.POST("/turnos/llamar", pantallaHandler.Receive)
		api.GET("/turnos/estado/*ubicacion", pantallaHandler.GetEstadoActual)

		// Rutas de Pacientes
		api.POST("/pacientes", pacienteHandler.AgregarPaciente)
		api.GET("/pacientes/llamar/:id/:nombre_caja", pacienteHandler.LlamarPaciente)
	}

	r.GET("/ws", pantallaHandler.ConnectWS)
	r.GET("/ws/pacientes/:nombre_caja", pacienteHandler.ConnectWS)

	return r
}
