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
        return true; // Permitir conexiones desde cualquier origen (CORS)
    },
}

// Hub gestiona a los clientes conectados y los mensajes
type Hub struct {
    Clients    map[*websocket.Conn]bool
    Broadcast  chan interface{} // Canal para recibir datos de cualquier tipo
    Register   chan *websocket.Conn
    Unregister chan *websocket.Conn
    Mutex      sync.Mutex
}

func NewHub() *Hub {
    return &Hub{
        Broadcast:  make(chan interface{}),
        Register:   make(chan *websocket.Conn),
        Unregister: make(chan *websocket.Conn),
        Clients:    make(map[*websocket.Conn]bool),
    }
}

// Run inicia el bucle principal del Hub (debe correr en una goroutine)
func (h *Hub) Run() {
    for {
        select {
        case client := <-h.Register:
            h.Mutex.Lock()
            h.Clients[client] = true
            h.Mutex.Unlock()

        case client := <-h.Unregister:
            h.Mutex.Lock()
            if _, ok := h.Clients[client]; ok {
                delete(h.Clients, client)
                client.Close()
            }
            h.Mutex.Unlock()

        case message := <-h.Broadcast:
            jsonMessage, err := json.Marshal(message)
            if err != nil {
                log.Printf("Error serializando mensaje: %v", err)
                continue
            }

            h.Mutex.Lock()
            for client := range h.Clients {
                err := client.WriteMessage(websocket.TextMessage, jsonMessage)
                if err != nil {
                    log.Printf("Error enviando mensaje: %v", err)
                    client.Close()
                    delete(h.Clients, client)
                }
            }
            h.Mutex.Unlock()
        }
    }
}

// HandleConnections es el endpoint al que se conectan las Pantallas (Frontend)
func (h *Hub) HandleConnections(c *gin.Context) {
    ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Fatal(err)
        return
    }
    h.Register <- ws
}