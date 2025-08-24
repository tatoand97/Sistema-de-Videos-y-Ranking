import pika
import json
from datetime import datetime

# Conexión a RabbitMQ
connection = pika.BlockingConnection(
    pika.ConnectionParameters('localhost', 5672, '/', pika.PlainCredentials('admin', 'password'))
)
channel = connection.channel()

# Declarar la cola
channel.queue_declare(queue='audio_removal_queue', durable=True)

# Mensaje
message = {
    "filename": "SmoothCriminal_MichaelJackson.mp4"
}

# Enviar mensaje
channel.basic_publish(
    exchange='',
    routing_key='audio_removal_queue',
    body=json.dumps(message),
    properties=pika.BasicProperties(delivery_mode=2)  # Mensaje persistente
)

print(f"Mensaje enviado: {message}")
connection.close()