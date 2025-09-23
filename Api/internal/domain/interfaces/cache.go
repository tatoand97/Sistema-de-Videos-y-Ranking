package interfaces

import "context"

// Cache abstrae un almacenamiento clave-valor de solo lectura usado por la capa de API.
type Cache interface {
	// GetBytes returns the value stored at key or an error if not found.
	GetBytes(ctx context.Context, key string) ([]byte, error)
}
