package config

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	RabbitMQ RabbitMQConfig `yaml:"rabbitmq"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type Secrets struct {
	JWT JWTConfig `yaml:"jwt"`
}

type JWTConfig struct {
	SecretKey string `yaml:"secret_key"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
}

type RabbitMQConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	VHost    string `yaml:"vhost"`
}
