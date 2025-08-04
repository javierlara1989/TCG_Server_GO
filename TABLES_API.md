# Tablas API Documentation

## Descripción General

El sistema de tablas permite a los usuarios crear mesas de juego para partidas de TCG. Cada tabla puede ser pública o privada, y puede tener diferentes categorías y premios.

## Modelos de Datos

### Table
- `id`: Identificador único de la tabla
- `category`: Categoría de la tabla (S, A, B, C, D)
- `privacy`: Privacidad de la tabla (private, public)
- `password`: Contraseña numérica opcional (máximo 10 dígitos)
- `prize`: Tipo de premio (money, card, aura)
- `amount`: Cantidad apostada (dinero o cartas) - entero opcional
- `winner`: Indica si hay un ganador (TRUE, FALSE, NULL)
- `created_at`: Fecha de creación
- `updated_at`: Fecha de última actualización
- `finished_at`: Fecha de finalización (opcional)

### UserTable
- `id`: Identificador único de la asociación
- `user_id`: ID del usuario propietario de la tabla
- `rival_id`: ID del rival (NULL si está esperando rival)
- `table_id`: ID de la tabla asociada

## Endpoints

### 1. Crear Tabla
**POST** `/api/tables`

Crea una nueva tabla y la asocia al usuario autenticado. La tabla se crea en estado de espera (sin rival).

#### Headers
```
Authorization: Bearer <token>
Content-Type: application/json
```

#### Request Body
```json
{
  "category": "A",
  "privacy": "public",
  "password": "1234",
  "prize": "money",
  "amount": 1000
}
```

#### Campos Requeridos
- `category`: Debe ser S, A, B, C, o D
- `privacy`: Debe ser "private" o "public"
- `prize`: Debe ser "money", "card", o "aura"

#### Campos Opcionales
- `password`: Contraseña numérica (máximo 10 dígitos)
- `amount`: Cantidad apostada (número entero positivo)

#### Response (201 Created)
```json
{
  "message": "Table created successfully",
  "table_id": 1
}
```

### 2. Obtener Tablas del Usuario
**GET** `/api/tables`

Obtiene todas las tablas asociadas al usuario autenticado (como propietario o rival).

#### Headers
```
Authorization: Bearer <token>
```

#### Response (200 OK)
```json
{
  "tables": [
    {
      "id": 1,
      "user_id": 1,
      "rival_id": null,
      "table_id": 1,
             "table": {
         "id": 1,
         "category": "A",
         "privacy": "public",
         "password": null,
         "prize": "money",
         "amount": 1000,
         "winner": null,
         "created_at": "2024-01-01T12:00:00Z",
         "updated_at": "2024-01-01T12:00:00Z",
         "finished_at": null
       },
      "user_name": "Usuario1",
      "user_email": "usuario1@example.com",
      "rival_name": null,
      "rival_email": null
    }
  ]
}
```

### 3. Actualizar Tabla
**PUT** `/api/tables/{id}`

Actualiza los parámetros de una tabla. Solo se puede actualizar si:
- El usuario es el propietario de la tabla
- La tabla está esperando rival (rival_id es NULL)

#### Headers
```
Authorization: Bearer <token>
Content-Type: application/json
```

#### Request Body
```json
{
  "category": "B",
  "privacy": "private",
  "password": "5678",
  "prize": "card",
  "amount": 500
}
```

#### Campos Modificables
- `category`: Nueva categoría (S, A, B, C, D)
- `privacy`: Nueva privacidad (private, public)
- `password`: Nueva contraseña (numérica, máximo 10 dígitos)
- `prize`: Nuevo premio (money, card, aura)
- `amount`: Nueva cantidad apostada (número entero positivo)

#### Response (200 OK)
```json
{
  "message": "Table updated successfully",
  "table_id": 1
}
```

## Validaciones

### Categorías Válidas
- S, A, B, C, D

### Privacidad Válida
- private, public

### Premios Válidos
- money, card, aura

### Contraseña
- Máximo 10 caracteres
- Solo dígitos numéricos (0-9)
- Opcional

### Cantidad (Amount)
- Número entero positivo
- Representa la cantidad de dinero o cartas apostadas
- Opcional

## Estados de la Tabla

1. **Esperando Rival**: `rival_id` es NULL
   - Se pueden modificar los parámetros
   - La tabla está disponible para que otros usuarios se unan

2. **Con Rival**: `rival_id` tiene un valor
   - No se pueden modificar los parámetros
   - La partida puede comenzar

3. **Finalizada**: `finished_at` tiene un valor
   - `winner` indica el resultado
   - La tabla está cerrada

## Códigos de Error

- `400 Bad Request`: Datos de entrada inválidos
- `401 Unauthorized`: Token de autenticación inválido o faltante
- `403 Forbidden`: No tienes permisos para realizar la acción
- `500 Internal Server Error`: Error interno del servidor

## Ejemplos de Uso

### Crear una tabla pública
```bash
curl -X POST http://localhost:8080/api/tables \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "category": "A",
    "privacy": "public",
    "prize": "money",
    "amount": 1000
  }'
```

### Crear una tabla privada con contraseña
```bash
curl -X POST http://localhost:8080/api/tables \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "category": "S",
    "privacy": "private",
    "password": "123456",
    "prize": "card",
    "amount": 500
  }'
```

### Actualizar parámetros de una tabla
```bash
curl -X PUT http://localhost:8080/api/tables/1 \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "category": "B",
    "privacy": "private",
    "password": "9999",
    "amount": 2000
  }'
``` 