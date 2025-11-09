# Base idéntica a la que usas en compose
FROM redis:7-alpine

# Expone el puerto estándar
EXPOSE 6379

# Declara volumen para persistencia (equivale a /data del compose)
VOLUME ["/data"]

# Healthcheck como en tu compose
HEALTHCHECK --interval=10s --timeout=3s --retries=5 \
  CMD redis-cli -h 127.0.0.1 -p 6379 ping || exit 1

# Comando por defecto (AOF activado)
CMD ["redis-server","--appendonly","yes"]
