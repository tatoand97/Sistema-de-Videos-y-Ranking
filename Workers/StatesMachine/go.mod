module statesmachine

go 1.21

require (
	github.com/streadway/amqp v1.1.0
	github.com/joho/godotenv v1.4.0
	github.com/sirupsen/logrus v1.9.3
	shared v0.0.0
)

replace shared => ./shared