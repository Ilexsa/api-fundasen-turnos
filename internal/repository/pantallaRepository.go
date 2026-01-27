package repositories

import (
	"api-turnos/internal/models"
	"database/sql"
	"log"
)

type PantallaReposotory interface{
	AgregarTurno(data models.Respuesta) (models.TurnoRequest, error)
	ObtenerUltimosTurnos(limit int, ubicacion string) ([]models.TurnoRequest, error)
}

type dbPantalla struct {
	database *sql.DB
}

func NewPantallaRepository (database *sql.DB) PantallaReposotory{
	return &dbPantalla{database: database}
}

func (repo *dbPantalla) AgregarTurno(data models.Respuesta) (models.TurnoRequest, error){
	query :=`EXEC sp_GestionarTurnoPantalla @p1, @p2, @p3`
	var resultado models.TurnoRequest
	err := repo.database.QueryRow(query, data.ConsultaID, data.Localidad, data.ConsultorioID).Scan(&resultado.ConsultaID, 
		&resultado.Paciente, &resultado.Medico, &resultado.Especialidad,
		&resultado.Consultorio, &resultado.Ubicacion, &resultado.Localidad)
	if err != nil{
		return models.TurnoRequest{},err
	}

	log.Printf("resultado: %+v", resultado)

	return resultado, nil

}

func (repo *dbPantalla) ObtenerUltimosTurnos(limit int, ubicacion string) ([]models.TurnoRequest, error) {

    query := `
-- Declaración de variables
DECLARE @limite INT = @p1; -- Puse 10 por defecto, cámbialo a lo que necesites
DECLARE @ubicacionEscogida VARCHAR(100) = @p2; -- Si lo dejas NULL o '', trae todo. Prueba poner 'Piso 1'

WITH TurnosOrdenados AS (
    SELECT
        ExternalConsultaID as id_consulta,
        PacienteNombre as paciente,
        MedicoNombre as medico,
        Especialidad as especialidad,
        Consultorio as consultorio,
        Ubicacion as ubicacion,
        Localidad as localidad,
        UltimoLLamado,
        -- Numeramos por consultorio para sacar el más nuevo
        ROW_NUMBER() OVER (
            PARTITION BY Consultorio 
            ORDER BY UltimoLLamado DESC
        ) as Fila
    FROM TurnosPantalla
    WHERE 
        -- Filtro de fecha (Solo hoy)
        CAST(UltimoLLamado AS DATE) = CAST(GETDATE() AS DATE)
)
SELECT TOP (@limite)
    id_consulta,
    paciente,
    medico,
    especialidad,
    consultorio,
    ubicacion,
    localidad
FROM TurnosOrdenados
WHERE Fila = 1 
  AND (
      -- Lógica del filtro opcional:
      @ubicacionEscogida IS NULL 
      OR @ubicacionEscogida = '' 
      OR ubicacion = @ubicacionEscogida
  )
ORDER BY UltimoLLamado DESC;
    `
    rows, err := repo.database.Query(query, limit, ubicacion)
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