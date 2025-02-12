package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env          string       `yaml:"env" env-default:"prod"`
	Grpc         Grpc         `yaml:"grpc"`
	HTTP         HTTP         `yaml:"http"`
	Storage      Storage      `yaml:"storage"`
	URLGenerator URLGenerator `yaml:"url_generator"`
}

type Grpc struct {
	Address string        `yaml:"address" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-default:"5s"`
}

type HTTP struct {
	Address string `yaml:"address" env-required:"true"`
}

type Storage struct {
	Type     string         `yaml:"type" env-required:"true" env:"STORAGE_TYPE"`
	Postgres PostgresConfig `yaml:"postgres,omitempty"`
}

type PostgresConfig struct {
	Host     string     `yaml:"host"`
	Port     int        `yaml:"port"`
	User     string     `yaml:"user"`
	Password string     `yaml:"password"`
	DBName   string     `yaml:"dbname"`
	SSLMode  string     `yaml:"sslmode"`
	Pool     PoolConfig `yaml:"pool"`
}

type PoolConfig struct {
	MaxOpenConns    int           `yaml:"max_open_conns" env-default:"25"`
	MaxIdleConns    int           `yaml:"max_idle_conns" env-default:"5"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time" env-default:"10m"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env-default:"1h"`
}

type URLGenerator struct {
	MaxAttempt int `yaml:"max_attempt" env-default:"5"`
	Length     int `yaml:"length" env-default:"10"`
}

func (p *PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		p.Host, p.Port, p.User, p.Password, p.DBName, p.SSLMode,
	)
}

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Cannot read config: %s", err)
	}

	switch cfg.Storage.Type {
	case "postgres":
		if cfg.Storage.Postgres.Pool.MaxOpenConns < 1 {
			log.Fatal("postgres.pool.max_open_conns must be >= 1")
		}
		if cfg.Storage.Postgres.Pool.MaxIdleConns < 0 {
			log.Fatal("postgres.pool.max_idle_conns must be >= 0")
		}
		if cfg.Storage.Postgres.Host == "" ||
			cfg.Storage.Postgres.Port == 0 ||
			cfg.Storage.Postgres.User == "" ||
			cfg.Storage.Postgres.Password == "" ||
			cfg.Storage.Postgres.DBName == "" {
			log.Fatal("Postgres configuration is incomplete")
		}
	case "memory":
		// Нет требований
	default:
		log.Fatal("Unsupported storage type")
	}

	return &cfg
}
