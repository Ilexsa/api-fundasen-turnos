package models

type Paciente struct {
    IdPaciente  string `json:"id_paciente"`
    Nombre      string `json:"nombre"`
    Dni 		string `json:"dni"`
    Extranjero bool `json:"extranjero"`
    TipoIdentificacion string `json:"tipo_identificacion"`
    FechaNacimiento string `json:"fecha_nacimiento"`
    Sexo string `json:"sexo"`
	Estado string `json:"estado"`
	Direccion string `json:"direccion"`
	Telefono string `json:"telefono"`
	Email string `json:"email"`
}