package handlers

import (
	"api-turnos/internal/models"
	"api-turnos/internal/services"
	"api-turnos/internal/websockets"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type PantallaHandlers struct {
	service services.PantallaService
	hub     *websockets.Hub
}

func NewPantallaHandle(service services.PantallaService, hub *websockets.Hub) *PantallaHandlers {
	return &PantallaHandlers{
		service: service,
		hub:     hub,
	}
}

const PANTALLA_ROOM = "pantallas"

func (h *PantallaHandlers) Receive(c *gin.Context) {
	var req models.Respuesta
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON INVALIDO"})
		log.Printf("Error JSON invalido: %v", err)
		return
	}

	turno, err := h.service.AgregarTurno(req)
	if err != nil {
		log.Printf("Error en servicio: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error procesando turno"})
		return
	}

	// Enviar mensaje a la sala general de pantallas
	msg := websockets.Message{
		Room: PANTALLA_ROOM,
		Data: turno,
	}
	h.hub.Broadcast <- msg

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"mensaje": "Turno generado y enviado a pantallas",
	})
}

func (h *PantallaHandlers) ConnectWS(c *gin.Context) {
	h.hub.HandleConnections(c, PANTALLA_ROOM)
}

func (h *PantallaHandlers) GetEstadoActual(c *gin.Context) {
	ubicacionParam := strings.TrimPrefix(c.Param("ubicacion"), "/")
	turnos, err := h.service.ObtenerUltimosTurnos(10, ubicacionParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error DB"})
		return
	}
	c.JSON(http.StatusOK, turnos)
}
