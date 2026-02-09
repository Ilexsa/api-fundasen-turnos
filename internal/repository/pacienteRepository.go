package repositories

import (
	"api-turnos/internal/models"
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type PacienteReposotory interface {
	AgregarPaciente(data models.Paciente) (models.Paciente, error)
	ObtenerPaciente(idPaciente string) (models.Paciente, error)
}

type dbPaciente struct {
	database *sql.DB
}

func NewPacienteRepository(database *sql.DB) PacienteReposotory {
	return &dbPaciente{database: database}
}

func (repo *dbPaciente) AgregarPaciente(data models.Paciente) (models.Paciente, error) {

	queryObtenerUltimoId := `SELECT secuencia from SEG_Secuencia where id_secuencia = 27`
	query := `
		INSERT INTO Cme_Paciente (IdPaciente, Nombre, Fecha_Ing, Ruc, Extranjero, TipoIdentificacion, 
		FechaNacimiento, Sexo, Estado, Direccion, Telefono1, Email)
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, @p11, @p12)
	`
	var ultimoId int
	errId := repo.database.QueryRow(queryObtenerUltimoId).Scan(&ultimoId)
	if errId != nil {
		return models.Paciente{}, fmt.Errorf("error getting last id: %v", errId)
	}
	ultimoId = ultimoId + 1
	err := repo.database.QueryRow(query,
		"C"+strconv.Itoa(ultimoId),
		data.Nombre,
		time.Now(),
		data.Dni,
		data.Extranjero,
		data.TipoIdentificacion,
		data.FechaNacimiento,
		data.Sexo,
		data.Estado,
		data.Direccion,
		data.Telefono,
		data.Email,
	)
	queryUpdateSecuencia := `UPDATE SEG_Secuencia SET secuencia = @p1 WHERE id_secuencia = 27`
	errUpdate := repo.database.QueryRow(queryUpdateSecuencia, ultimoId)
	if errUpdate != nil {
		return models.Paciente{}, fmt.Errorf("error updating secuencia: %v", errUpdate)
	}
	if err != nil {
		return models.Paciente{}, fmt.Errorf("error inserting paciente: %v", err)
	}
	return data, nil
}

func (repo *dbPaciente) ObtenerPaciente(idPaciente string) (models.Paciente, error) {

	query := `
		SELECT Id_Paciente, Nombre, Ruc, Extranjero, TipoIdentificacion, Fecha_Nac, Sexo, EstCiv, Direccion, Telefono1, Email
		FROM Cme_Paciente
		WHERE Id_Paciente = @p1
	`
	var paciente models.Paciente
	err := repo.database.QueryRow(query, idPaciente).Scan(
		&paciente.IdPaciente,
		&paciente.Nombre,
		&paciente.Dni,
		&paciente.Extranjero,
		&paciente.TipoIdentificacion,
		&paciente.FechaNacimiento,
		&paciente.Sexo,
		&paciente.Estado,
		&paciente.Direccion,
		&paciente.Telefono,
		&paciente.Email,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Paciente{}, fmt.Errorf("paciente not found with id %s", idPaciente)
		}
		return models.Paciente{}, fmt.Errorf("error fetching paciente: %w", err)
	}

	return paciente, nil
}
