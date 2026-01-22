package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init(){
	// Intenta cargar .env para desarrollo; en producción se suelen usar variables de entorno.
	if err := godotenv.Load(); err != nil {
		log.Println(".env no cargado (ok si estas en prod):", err)
	}
}

func Get(key, def string) string{
	value, ok := os.LookupEnv(key)
	if ok{
		return value
	}

	log.Println(key, ":valor default")
	return  def
}

func MustGet(key string) string {
	value, ok := os.LookupEnv(key)
	if ok && value != "" {
		return value
	}
	log.Fatalln("Falta variable de entorno:", key)
	return ""
}