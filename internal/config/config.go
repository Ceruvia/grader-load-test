package config

import (
	"log"

	"github.com/joho/godotenv"
)

type MachineryConfig struct {
	BrokerURL        string
	ResultBackendURL string
	QueueName        string
	ResultsExpireIn  int
}

type Config struct {
	MachineryCfg *MachineryConfig
}

func loadEnvFile() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func GetAppConfig() *Config {
	loadEnvFile()

	machineryConfig := &MachineryConfig{
		BrokerURL:        GetString("QUEUE_BROKER_URL", "amqp://guest:guest@localhost:5672/"),
		ResultBackendURL: GetString("QUEUE_RESULT_URL", "amqp://guest:guest@localhost:5672/"),
		QueueName:        GetString("QUEUE_NAME", "ceruvia_submissions"),
		ResultsExpireIn:  GetInt("QUEUE_RESULT_TTL", 36000),
	}

	return &Config{
		MachineryCfg: machineryConfig,
	}
}
