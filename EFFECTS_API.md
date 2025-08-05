# Effects API - Uso Interno

Este documento describe el sistema de efectos internos para el TCG Server. Los efectos son entidades de uso interno que no requieren rutas públicas de API.

## Estructura de Base de Datos

### Tabla `effects`
- `id` (INT, AUTO_INCREMENT, PRIMARY KEY)
- `description` (TEXT, NOT NULL) - Descripción del efecto
- `created_at` (TIMESTAMP, DEFAULT CURRENT_TIMESTAMP)
- `updated_at` (TIMESTAMP, DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP)
- `deleted_at` (TIMESTAMP, NULL) - Para soft delete

### Tabla `card_effects`
- `card_id` (INT, NOT NULL) - Referencia a la tabla cards
- `effect_id` (INT, NOT NULL) - Referencia a la tabla effects
- PRIMARY KEY (card_id, effect_id)
- FOREIGN KEY constraints con CASCADE DELETE

## Modelos

### Effect
```go
type Effect struct {
    ID          int        `json:"id" db:"id"`
    Description string     `json:"description" db:"description"`
    CreatedAt   time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
    DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}
```

### CardEffect
```go
type CardEffect struct {
    CardID   int `json:"card_id" db:"card_id"`
    EffectID int `json:"effect_id" db:"effect_id"`
}
```

## Funciones de Base de Datos Disponibles

### Gestión de Efectos
- `GetEffectByID(id int)` - Obtener efecto por ID
- `GetAllEffects()` - Obtener todos los efectos no eliminados
- `SoftDeleteEffect(id int)` - Soft delete del efecto
- `HardDeleteEffect(id int)` - Eliminación permanente del efecto y sus relaciones

### Gestión de Relaciones Card-Effect
- `CreateCardEffect(cardID, effectID int)` - Asociar carta con efecto
- `GetEffectsByCardID(cardID int)` - Obtener todos los efectos de una carta
- `GetCardsByEffectID(effectID int)` - Obtener todas las cartas con un efecto específico
- `DeleteCardEffect(cardID, effectID int)` - Eliminar relación específica
- `DeleteAllCardEffects(cardID int)` - Eliminar todas las relaciones de una carta

## Uso Interno

Los efectos están diseñados para uso interno del sistema. No se exponen rutas públicas para su gestión, ya que serán programados y gestionados internamente por la lógica del juego.

### Ejemplo de Uso Interno
```go
// Obtener todos los efectos disponibles
rows, err := database.GetAllEffects()
if err != nil {
    // Manejar error
}

// Asociar efecto a una carta
err = database.CreateCardEffect(cardID, effectID)
if err != nil {
    // Manejar error
}

// Obtener efectos de una carta
rows, err := database.GetEffectsByCardID(cardID)
if err != nil {
    // Manejar error
}
```

## Notas Importantes

1. **Uso Interno**: Los efectos no tienen rutas públicas de API
2. **Seeder**: Los efectos serán llenados por un seeder más adelante, no se pueden crear ni actualizar programáticamente
3. **Soft Delete**: Los efectos usan soft delete por defecto
4. **Relaciones**: Una carta puede tener múltiples efectos y un efecto puede estar en múltiples cartas
5. **Integridad**: Las relaciones se eliminan automáticamente cuando se elimina una carta o efecto
6. **Índices**: Se han creado índices para optimizar las consultas frecuentes

## Implementación Futura

Los efectos serán programados más tarde para implementar la lógica del juego, permitiendo que las cartas tengan comportamientos especiales durante las partidas. 