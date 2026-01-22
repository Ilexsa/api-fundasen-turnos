package models

type TurnoRequest struct {
    ConsultaID  string `json:"consulta_id"`
    Medico      string `json:"medico"`
    Paciente    string `json:"paciente"`
    Consultorio string `json:"consultorio"`
    Especialidad string `json:"especialidad"`
	Localidad string `json:"localidad"`
	Ubicacion string `json:"ubicacion"`
}

type WebSocketMessage struct {
    Tipo string      `json:"tipo"` 
    Data TurnoRequest `json:"data"`
}

type Respuesta struct {
	ConsultaID  string `json:"consulta_id"`
	Localidad string `json:"localidad"`
}