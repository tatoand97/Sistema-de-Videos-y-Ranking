# Testing Architecture - Clean Architecture Pattern

## 📁 **Estructura Reorganizada**

```
Sistema-de-Videos-y-Ranking/
├── Api/
│   └── tests/
│       ├── mocks/                    # Mocks para interfaces
│       ├── unit/
│       │   ├── application/          # Tests de casos de uso
│       │   └── domain/               # Tests de entidades
│       └── integration/              # Tests de handlers
│
├── Workers/
│   ├── AudioRemoval/
│   │   └── tests/
│   │       ├── mocks/
│   │       └── unit/
│   │           ├── application/
│   │           └── domain/
│   ├── EditVideo/
│   │   └── tests/
│   │       ├── mocks/
│   │       └── unit/
│   │           ├── application/
│   │           └── domain/
│   ├── TrimVideo/
│   │   └── tests/
│   │       ├── mocks/
│   │       └── unit/
│   │           ├── application/
│   │           └── domain/
│   ├── Watermarking/
│   │   └── tests/
│   │       ├── mocks/
│   │       └── unit/
│   │           ├── application/
│   │           └── domain/
│   ├── gossipOpenClose/
│   │   └── tests/
│   │       ├── mocks/
│   │       └── unit/
│   │           ├── application/
│   │           └── domain/
│   └── StatesMachine/
│       └── tests/
│           ├── mocks/
│           └── unit/
│               ├── application/
│               └── domain/
│
└── shared/
    └── testing/                      # Helpers compartidos
```

## 🎯 **Principios de Arquitectura Limpia**

### **1. Separación por Capas**
- **Domain**: Tests de entidades y reglas de negocio
- **Application**: Tests de casos de uso y servicios
- **Infrastructure**: Mocks de repositorios y servicios externos
- **Presentation**: Tests de integración de handlers

### **2. Dependencias**
- Tests de dominio: Sin dependencias externas
- Tests de aplicación: Usan mocks de infraestructura
- Tests de integración: Usan servicios reales o containers

### **3. Ubicación**
- Cada worker tiene sus tests junto al código
- API tiene estructura similar para consistencia
- Shared contiene helpers comunes

## 🚀 **Comandos de Testing**

### **Todos los tests**
```bash
make test-workers
```

### **Por componente**
```bash
cd Api && go test -v ./tests/unit/...
cd Workers/AudioRemoval && go test -v ./tests/unit/...
cd Workers/EditVideo && go test -v ./tests/unit/...
```

### **Con cobertura**
```bash
make coverage-workers
```

## ✅ **Beneficios de la Reorganización**

1. **Consistencia**: Todos los workers siguen el mismo patrón
2. **Mantenibilidad**: Tests cerca del código que prueban
3. **Escalabilidad**: Fácil agregar nuevos workers con tests
4. **Arquitectura Limpia**: Respeta las capas y dependencias
5. **CI/CD**: Estructura uniforme para pipelines

## 🔄 **Migración Completada**

- ✅ AudioRemoval: Movido de raíz a Workers/
- ✅ EditVideo: Movido de raíz a Workers/
- ✅ TrimVideo: Movido de raíz a Workers/
- ✅ Watermarking: Estructura creada
- ✅ gossipOpenClose: Estructura creada
- ✅ StatesMachine: Estructura creada
- ✅ Makefile: Actualizado con nuevas rutas
- ✅ Duplicados: Eliminados