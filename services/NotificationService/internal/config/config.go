package config

type (
	RabbitMQURL         string
	SubscriberQueueName string
)

type Config struct {
	RabbitURL string `env:"RABBIT_URL" envDefault:"amqp://guest:guest@localhost:5672/"`
}
