# Configuración de Pruebas Unitarias para Workers

## Estructura Creada

```
AudioRemoval/
├── tests/
│   ├── unit/
│   │   ├── domain/
│   │   │   └── audio_removal_test.go
│   │   └── application/
│   │       ├── audio_removal_usecase_test.go
│   │       └── message_handler_test.go
│   └── mocks/
EditVideo/
├── tests/
│   └── unit/
│       └── domain/
│           └── edit_video_test.go
TrimVideo/
├── tests/
│   └── unit/
│       └── domain/
│           └── trim_video_test.go
shared/
└── testing/
    └── test_helpers.go
```

## Comandos para Ejecutar

```bash
# Instalar dependencias de testing
make install-test-deps

# Ejecutar tests de todos los workers
make test-workers

# Ejecutar tests con cobertura
make coverage-workers

# Limpiar artefactos de test
make clean-test
```

## Dependencias Agregadas

- `github.com/stretchr/testify` - Assertions y testing utilities
- `github.com/golang/mock` - Mock generation
- Shared testing helpers en `/shared/testing/`

## Patrones Implementados

### 1. **Domain Layer Tests**
- Validación de entidades
- Reglas de negocio
- Sin dependencias externas

### 2. **Application Layer Tests**
- Use cases con mocks
- Message handlers
- Dependency injection

### 3. **Mocks**
- Interfaces para repositories
- Services externos (Storage, Logger)
- Message brokers

## Cobertura Objetivo

- **Domain**: >90%
- **Application**: >85%
- **Infrastructure**: >70%

## CI/CD Integration

- GitHub Actions configurado
- Tests automáticos en PR/push
- Reporte de cobertura por worker

La estructura mantiene la arquitectura limpia separando las pruebas por capas y usando dependency injection para facilitar el testing.