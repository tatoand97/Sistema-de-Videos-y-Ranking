#!/bin/sh

# Esperar a que RabbitMQ esté disponible
echo "Waiting for RabbitMQ..."
until nc -z rabbitmq 5672; do
  echo "RabbitMQ is unavailable - sleeping"
  sleep 2
done
echo "RabbitMQ is up - executing command"

# Esperar a que PostgreSQL esté disponible
echo "Waiting for PostgreSQL..."
until nc -z postgres 5432; do
  echo "PostgreSQL is unavailable - sleeping"
  sleep 2
done
echo "PostgreSQL is up - executing command"

# Ejecutar la aplicación
exec "$@"