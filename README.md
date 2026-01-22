# API Turnos (pantalla) - Go + SQL Server

## Requisitos
- Go 1.23+
- SQL Server con las tablas + stored procedures creadas (los SPs que definimos)

## Configuración
Edita `.env`:

```ini
API_ADDR=":8080"
API_TOKEN="..."     # opcional, si está vacío no se exige auth
ALLOWED_ORIGINS="*"

DB_HOST="localhost"
DB_PORT="1433"
DB_USER="sa"
DB_PASS="..."
DB_ENCRYPT="disable"
DB_NAME="TURNOS"
```

## Ejecutar
```bash
go run ./cmd/server
```

## Endpoints
Base: `/api/v1`

### Lectura (pantalla)
- `GET /api/v1/screen/waiting?queueId=&consultingRoomId=&top=`
- `GET /api/v1/screen/calling?queueId=&consultingRoomId=&top=`

### Realtime (SSE)
- `GET /api/v1/stream?branchCode=&queueId=&consultingRoomId=`

Ejemplo JS:
```js
const es = new EventSource('/api/v1/stream?branchCode=ALB&consultingRoomId=18')
es.addEventListener('ticket', (e) => {
  const evt = JSON.parse(e.data)
  // evt.type: CALLED, RECALLED, REASSIGNED_ROOM, IN_SERVICE, DONE, CANCELLED, NO_SHOW
  console.log(evt)
})
```

### Mutaciones (requieren Bearer token si `API_TOKEN` no está vacío)
- `POST /api/v1/tickets/call`   (OPCIÓN B: crea si no existe y deja en CALLING)
- `POST /api/v1/tickets`
- `POST /api/v1/tickets/:id/status`
- `POST /api/v1/tickets/:id/recall`
- `POST /api/v1/tickets/:id/reassign-room`

#### Call upsert (opción B)
El sistema externo solo hace el “llamado”; si el turno no existe, la API lo crea y lo deja en CALLING.

Requiere el SP `dbo.usp_Ticket_CallUpsert` (script en `sql/usp_Ticket_CallUpsert.sql`).

```bash
curl -X POST http://localhost:8080/api/v1/tickets/call \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TU_TOKEN" \
  -d '{
    "sourceSystem":"SistemaClinicoX",
    "branchCode":"ALB",
    "externalId":"28349",
    "patientDisplayName":"Maria F.",
    "specialtyId": 1,
    "consultingRoomId": 18,
    "queueId": 1,
    "payload": {"by":"dr", "reason":"llamado"}
  }'
```

#### Crear turno (ejemplo)
```bash
curl -X POST http://localhost:8080/api/v1/tickets \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TU_TOKEN" \
  -d '{
    "sourceSystem":"SistemaClinicoX",
    "branchCode":"ALB",
    "externalId":"28349",
    "patientDisplayName":"Maria F.",
    "specialtyId": 1,
    "consultingRoomId": 18,
    "priority": 0
  }'
```

#### Cambiar estado
```bash
curl -X POST http://localhost:8080/api/v1/tickets/<TICKET_UUID>/status \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TU_TOKEN" \
  -d '{"newStatus":1, "payload":"{\"by\":\"kiosk\"}"}'
```

Estados:
- 0 WAITING
- 1 CALLING
- 2 IN_SERVICE
- 3 DONE
- 4 CANCELLED
- 5 NO_SHOW

### Admin/Bootstrap (Upserts)
- `POST /api/v1/admin/specialties/upsert`
- `POST /api/v1/admin/consulting-rooms/upsert`
- `POST /api/v1/admin/doctors/upsert`
- `POST /api/v1/admin/queues/upsert`
- `POST /api/v1/admin/doctor-assignments/upsert`

> Nota: En esta versión, el `isActive` es opcional; si no lo envías, se asume `true`.
