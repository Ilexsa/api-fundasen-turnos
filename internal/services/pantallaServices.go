package services

import(
	"api-turnos/internal/models"
	"api-turnos/internal/repository"
)
type PantallaService interface{
		AgregarTurno(data models.Respuesta) (models.TurnoRequest, error)
		ObtenerUltimosTurnos(limit int) ([]models.TurnoRequest, error)
}

type pantallaService struct{
	repo repositories.PantallaReposotory
}

func (s *pantallaService) AgregarTurno(data models.Respuesta) (models.TurnoRequest, error){
	return s.repo.AgregarTurno(data)
}

func (s *pantallaService) ObtenerUltimosTurnos(limit int) ([]models.TurnoRequest, error){
	return s.repo.ObtenerUltimosTurnos(limit)
}