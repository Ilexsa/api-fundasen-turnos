package websockets

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Permitir conexiones desde cualquier origen (CORS)
	},
}

// Subscription define la conexión y la sala a la que pertenece
type Subscription struct {
	Conn *websocket.Conn
	Room string
}

// Message define el mensaje y la sala destino
type Message struct {
	Data interface{}
	Room string
}

// Hub gestiona a los clientes conectados y los mensajes por sala
type Hub struct {
	Clients    map[string]map[*websocket.Conn]bool // Mapa de salas -> clientes
	Broadcast  chan Message
	Register   chan Subscription
	Unregister chan Subscription
	Mutex      sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan Message),
		Register:   make(chan Subscription),
		Unregister: make(chan Subscription),
		Clients:    make(map[string]map[*websocket.Conn]bool),
	}
}

// Run inicia el bucle principal del Hub (debe correr en una goroutine)
func (h *Hub) Run() {
	for {
		select {
		case subscription := <-h.Register:
			h.Mutex.Lock()
			if _, ok := h.Clients[subscription.Room]; !ok {
				h.Clients[subscription.Room] = make(map[*websocket.Conn]bool)
			}
			h.Clients[subscription.Room][subscription.Conn] = true
			h.Mutex.Unlock()

		case subscription := <-h.Unregister:
			h.Mutex.Lock()
			if clients, ok := h.Clients[subscription.Room]; ok {
				if _, ok := clients[subscription.Conn]; ok {
					delete(clients, subscription.Conn)
					subscription.Conn.Close()
					if len(clients) == 0 {
						delete(h.Clients, subscription.Room)
					}
				}
			}
			h.Mutex.Unlock()

		case message := <-h.Broadcast:
			jsonMessage, err := json.Marshal(message.Data)
			if err != nil {
				log.Printf("Error serializando mensaje: %v", err)
				continue
			}

			h.Mutex.Lock()
			if clients, ok := h.Clients[message.Room]; ok {
				for client := range clients {
					err := client.WriteMessage(websocket.TextMessage, jsonMessage)
					if err != nil {
						log.Printf("Error enviando mensaje a la sala %s: %v", message.Room, err)
						client.Close()
						delete(clients, client)
					}
				}
			}
			h.Mutex.Unlock()
		}
	}
}

func (h *Hub) HandleConnections(c *gin.Context, room string) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error al conectar: %v", err)
		return
	}

	subscription := Subscription{Conn: ws, Room: room}
	h.Register <- subscription

	// Escuchar mensajes de cierre del cliente para desuscribirlo correctamente
	go func() {
		defer func() {
			h.Unregister <- subscription
		}()
		for {
			_, _, err := ws.ReadMessage()
			if err != nil {
				break
			}
		}
	}()
}
