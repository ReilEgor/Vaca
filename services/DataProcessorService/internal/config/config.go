package config

type (
	StateClientAddr     string
	SearchClientAddr    string
	RabbitURL           string
	SubscriberQueueName string
	PublisherQueueName  string
)

type Config struct {
	// REDIS
	RedisHost            string `env:"REDIS_HOST" envDefault:"localhost"`
	RedisPort            string `env:"REDIS_PORT" envDefault:"6379"`
	RedisPassword        string `env:"REDIS_PASSWORD" envDefault:""`
	DSN                  string `env:"DB_SOURCE" envDefault:"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"`
	RabbitURL            string `env:"RABBIT_URL" envDefault:"amqp://guest:guest@localhost:5672/"`
	StateServiceAddress  string `env:"STATE_SERVICE_ADDRESS" envDefault:"http://localhost"`
	SearchServiceAddress string `env:"SEARCH_SERVICE_ADDRESS" envDefault:"http://localhost"`
	SubscriberQueueName  string `env:"SUBSCRIBER_QUEUE" envDefault:"vacancy_result_queue"`
	PublisherQueueName   string `env:"NOTIFICATION_QUEUE" envDefault:"notification_queue"`
}
