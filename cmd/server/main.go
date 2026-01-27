package main

import (
    "api-turnos/internal/db"
    "api-turnos/internal/routes"
    "api-turnos/internal/websockets"
    "log"
	"api-turnos/internal/config"

    _ "github.com/denisenkom/go-mssqldb" // Importante: Driver SQL Server
)

func main() {
    config := config.Load();
    err := db.ConnectDB(config.DBHost, config.DBPort,config.DBUser, config.DBPass,
		config.DBName, config.DBEncrypt) // Asumo que tienes esta func en tu db/dg.go
    if err != nil {
        log.Fatalf("No se pudo conectar a la BD: %v", err)
    }
	log.Println("Conexión a la base de datos establecida")

    // 2. Inicializar Hub de WebSockets
    hub := websockets.NewHub()
    go hub.Run() // ¡IMPORTANTE! Correr en una goroutine aparte

    // 3. Configurar Router e inyectar dependencias
    r := routes.SetupRouter(db.GetDb(), hub)

    // 4. Iniciar Servidor
    log.Println("Servidor corriendo en el puerto :9090")
    if err := r.Run(":9090"); err != nil {
        log.Fatalf("Error iniciando servidor: %v", err)
    }
}