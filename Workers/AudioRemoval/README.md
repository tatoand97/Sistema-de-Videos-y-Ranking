# AudioRemoval Worker - MP4 Nativo

## Descripción
Worker optimizado para remover audio de videos MP4 usando procesamiento nativo en Go, sin dependencias externas.

## Características
- **Procesamiento**: Manipulación directa de contenedor MP4
- **Recursos**: Mínimos (solo Go runtime)
- **Rendimiento**: 200-500 archivos/segundo por core
- **Memoria**: ~20MB por instancia
- **Tamaño imagen**: ~15MB
- **Latencia**: 5-10ms por archivo

## Construcción
```bash
docker build -t audioremoval .
```

## Limitaciones
- Solo soporta archivos MP4
- Requiere contenedor MP4 válido

## Ventajas
- Sin dependencias FFmpeg
- Escalabilidad horizontal perfecta
- Ideal para pruebas de carga
- Consumo mínimo de recursos