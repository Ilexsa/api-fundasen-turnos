package handlers

import (
	"api-turnos/internal/models"
	"api-turnos/internal/services"
	"api-turnos/internal/websockets"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PacienteHandlers struct {
	service services.PacienteService
	hub     *websockets.Hub
}

func NewPacienteHandler(service services.PacienteService, hub *websockets.Hub) *PacienteHandlers {
	return &PacienteHandlers{
		service: service,
		hub:     hub,
	}
}

func (h *PacienteHandlers) AgregarPaciente(c *gin.Context) {
	var req models.Paciente
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON INVALIDO"})
		log.Printf("Error JSON invalido: %v", err)
		return
	}

	paciente, err := h.service.AgregarPaciente(req)
	if err != nil {
		log.Printf("Error en servicio: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando paciente"})
		return
	}

	c.JSON(http.StatusOK, paciente)
}

func (h *PacienteHandlers) LlamarPaciente(c *gin.Context) {
	idPaciente := c.Param("id")
	nombreCaja := c.Param("nombre_caja")

	if idPaciente == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de paciente requerido"})
		return
	}

	paciente, err := h.service.ObtenerPaciente(idPaciente)
	if err != nil {
		log.Printf("Error obteniendo paciente: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Paciente no encontrado"})
		return
	}

	// Enviar mensaje a la sala específica
	msg := websockets.Message{
		Room: nombreCaja,
		Data: paciente,
	}
	h.hub.Broadcast <- msg

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"mensaje": "Paciente llamado y enviado a pantallas",
		"data":    paciente,
	})
}

func (h *PacienteHandlers) ConnectWS(c *gin.Context) {
	nombreCaja := c.Param("nombre_caja")
	if nombreCaja == "" {
		log.Println("Advertencia: Conexión WS Pacientes sin nombre_caja, usando 'default'")
		nombreCaja = "default"
	}
	h.hub.HandleConnections(c, nombreCaja)
}
