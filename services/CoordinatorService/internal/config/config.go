package config

type (
	StateClientAddr    string
	SearchClientAddr   string
	PublisherQueueName string
	RabbitURL          string
)

type Config struct {
	// --- Infrastructure: Message Broker & Cache ---
	RabbitURL     string `env:"RABBIT_URL" envDefault:"amqp://guest:guest@localhost:5672/"`
	RedisHost     string `env:"REDIS_HOST" envDefault:"localhost"`
	RedisPort     string `env:"REDIS_PORT" envDefault:"6379"`
	RedisPassword string `env:"REDIS_PASSWORD" envDefault:""`

	// --- Internal Microservices (gRPC/HTTP Addresses) ---
	SearchServiceAddress string `env:"SEARCH_SERVICE_ADDRESS" envDefault:"http://localhost:8082"`
	StateServiceAddress  string `env:"STATE_SERVICE_ADDRESS" envDefault:"http://localhost:8081"`

	// --- Server Settings ---
	RestApiPort string `env:"HTTP_PORT" envDefault:"8080"`
}
