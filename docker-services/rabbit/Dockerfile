# syntax=docker/dockerfile:1
FROM rabbitmq:3.13-management-alpine

ENV RABBITMQ_DEFAULT_USER=admin \
    RABBITMQ_DEFAULT_PASS=admin

COPY --chown=rabbitmq:rabbitmq enabled_plugins /etc/rabbitmq/enabled_plugins
COPY --chown=rabbitmq:rabbitmq rabbitmq.conf /etc/rabbitmq/rabbitmq.conf
COPY --chown=rabbitmq:rabbitmq definitions.json /etc/rabbitmq/definitions.json

HEALTHCHECK --interval=10s --timeout=5s --retries=5 CMD rabbitmq-diagnostics -q ping

EXPOSE 5672 15672
