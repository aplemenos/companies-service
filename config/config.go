package config

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"
)

// App config struct
type Config struct {
	Server      ServerConfig
	Postgres    PostgresConfig
	Redis       RedisConfig
	Cookie      Cookie
	Metrics     Metrics
	Logger      Logger
	Jaeger      Jaeger
	KafkaTopics KafkaTopics
	Kafka       *Kafka
}

// Server config struct
type ServerConfig struct {
	AppVersion           string
	Port                 string
	PprofPort            string
	Mode                 string
	JwtSecretKey         string
	CookieName           string
	ReadTimeout          time.Duration
	WriteTimeout         time.Duration
	SSL                  bool
	CtxDefaultTimeout    time.Duration
	Debug                bool
	CheckIntervalSeconds int
}

// Logger config
type Logger struct {
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
}

// Postgresql config
type PostgresConfig struct {
	PostgresqlHost     string
	PostgresqlPort     string
	PostgresqlUser     string
	PostgresqlPassword string
	PostgresqlDbname   string
	PostgresqlSSLMode  bool
	PgDriver           string
}

// Redis config
type RedisConfig struct {
	RedisAddr      string
	RedisPassword  string
	RedisDB        string
	RedisDefaultdb string
	MinIdleConns   int
	PoolSize       int
	PoolTimeout    int
	Password       string
	DB             int
}

// Cookie config
type Cookie struct {
	Name     string
	MaxAge   int
	Secure   bool
	HTTPOnly bool
}

// Metrics config
type Metrics struct {
	URL         string
	ServiceName string
}

// Jaeger
type Jaeger struct {
	Host        string
	ServiceName string
	LogSpans    bool
}

// KafkaTopics
type KafkaTopics struct {
	CompanyCreated TopicConfig
	CompanyUpdated TopicConfig
	CompanyDeleted TopicConfig
}

// TopicConfig kafka topic config
type TopicConfig struct {
	TopicName         string
	Partitions        int
	ReplicationFactor int
}

// Kafka
type Kafka struct {
	Brokers    []string
	InitTopics bool
}

// Load config file from given path
func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

// Parse config file
func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}

// Get config path for local or docker
func GetConfigPath(configPath string) string {
	if configPath == "docker" {
		return "./config/config-docker"
	}
	return "./config/config-local"
}
