package repositories

import (
	"api-turnos/internal/models"
	"database/sql"
	"log"
)

type PantallaReposotory interface{
	AgregarTurno(data models.Respuesta) (models.TurnoRequest, error)
	ObtenerUltimosTurnos(limit int) ([]models.TurnoRequest, error)
}

type dbPantalla struct {
	database *sql.DB
}

func NewPantallaRepository (database *sql.DB) PantallaReposotory{
	return &dbPantalla{database: database}
}

func (repo *dbPantalla) AgregarTurno(data models.Respuesta) (models.TurnoRequest, error){
	query :=`EXEC sp_GestionarTurnoPantalla @p1, @p2`
	var resultado models.TurnoRequest
	err := repo.database.QueryRow(query, data.ConsultaID, data.Localidad).Scan(&resultado.ConsultaID, 
		&resultado.Paciente, &resultado.Medico, &resultado.Especialidad,
		&resultado.Consultorio, &resultado.Ubicacion, &resultado.Localidad)
	if err != nil{
		return models.TurnoRequest{},err
	}

	log.Printf("resultado: %+v", resultado)

	return resultado, nil

}

func (repo *dbPantalla) ObtenerUltimosTurnos(limit int) ([]models.TurnoRequest, error) {

    query := `
        SELECT TOP (@p1)
            ExternalConsultaID as id_consulta,
            PacienteNombre as paciente,
            MedicoNombre as medico,
            Especialidad as especialidad,
            Consultorio as consultorio,
            Ubicacion as ubicacion,
            Localidad as localidad
        FROM TurnosPantalla
        ORDER BY UltimoLLamado DESC
    `
    rows, err := repo.database.Query(query, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var turnos []models.TurnoRequest
    for rows.Next() {
        var t models.TurnoRequest
        if err := rows.Scan(&t.ConsultaID, &t.Paciente, &t.Medico, &t.Especialidad, &t.Consultorio, &t.Ubicacion, &t.Localidad); err != nil {
            continue
        }
        turnos = append(turnos, t)
    }
    return turnos, nil
}