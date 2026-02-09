package services

import (
	"api-turnos/internal/models"
	repositories "api-turnos/internal/repository"
)

type PacienteService interface {
	AgregarPaciente(data models.Paciente) (models.Paciente, error)
	ObtenerPaciente(idPaciente string) (models.Paciente, error)
}

type pacienteService struct {
	repo repositories.PacienteReposotory
}

func NewPacienteService(repo repositories.PacienteReposotory) PacienteService {
	return &pacienteService{repo: repo}
}

func (s *pacienteService) AgregarPaciente(data models.Paciente) (models.Paciente, error) {
	return s.repo.AgregarPaciente(data)
}

func (s *pacienteService) ObtenerPaciente(idPaciente string) (models.Paciente, error) {
	return s.repo.ObtenerPaciente(idPaciente)
}
