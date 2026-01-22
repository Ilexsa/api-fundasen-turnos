package handlers

import (
	"api-turnos/internal/models"
	"api-turnos/internal/services"
	"api-turnos/internal/websockets"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PantallaHandlers struct {
    service services.PantallaService
    hub     *websockets.Hub
}

// Inyectamos tanto el Service como el Hub
func NewPantallaHandle(service services.PantallaService, hub *websockets.Hub) *PantallaHandlers {
    return &PantallaHandlers{
        service: service,
        hub:     hub,
    }
}
func (h *PantallaHandlers) Receive (c *gin.Context){
	var req models.Respuesta
	if err :=c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON INVALIDO"})
		log.Printf("Error JSON invalido: %v" ,err)
		return
	}

	turno, err := h.service.AgregarTurno(req)
    if err != nil {
        log.Printf("Error en servicio: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error procesando turno"})
        return
    }

    h.hub.Broadcast <- turno



    c.JSON(http.StatusOK, gin.H{
        "status": "ok",
        "mensaje": "Turno generado y enviado a pantallas",
    })
}

func (h *PantallaHandlers) ConnectWS(c *gin.Context) {
    h.hub.HandleConnections(c)
}

func (h *PantallaHandlers) GetEstadoActual(c *gin.Context) {
    // Pedimos los últimos 10 para llenar la pantalla
    turnos, err := h.service.ObtenerUltimosTurnos(10)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error DB"})
        return
    }
    c.JSON(http.StatusOK, turnos)
}