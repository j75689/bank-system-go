package config

import (
	"bank-system-go/pkg/logger"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Release  bool           `mapstructure:"release"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	HTTP     HTTPConfig     `mapstructure:"http"`
	DataBase DataBaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	MQ       MQConfig       `mapstructure:"mq"`
}

type LoggerConfig struct {
	Level  string           `mapstructure:"level"`
	Format logger.LogFormat `mapstructure:"format"`
}

type HTTPConfig struct {
	Port uint16 `mapstructure:"port"`
}

type DataBaseConfig struct {
	Driver         string        `mapstructure:"driver"`
	Host           string        `mapstructure:"host"`
	Port           uint          `mapstructure:"port"`
	Database       string        `mapstructure:"database"`
	InstanceName   string        `mapstructure:"instance_name"`
	User           string        `mapstructure:"user"`
	Password       string        `mapstructure:"password"`
	ConnectTimeout string        `mapstructure:"connect_timeout"`
	ReadTimeout    string        `mapstructure:"read_timeout"`
	WriteTimeout   string        `mapstructure:"write_timeout"`
	DialTimeout    time.Duration `mapstructure:"dial_timeout"`
	MaxLifetime    time.Duration `mapstructure:"max_lifetime"`
	MaxIdleTime    time.Duration `mapstructure:"max_idletime"`
	MaxIdleConn    int           `mapstructure:"max_idle_conn"`
	MaxOpenConn    int           `mapstructure:"max_open_conn"`
	SSLMode        bool          `mapstructure:"ssl_mode"`
}

type RedisConfig struct {
	Host         string        `mapstructure:"host"`
	Port         uint          `mapstructure:"port"`
	DB           int           `mapstructure:"db"`
	Password     string        `mapstructure:"password"`
	MinIdleConns int           `mapstructure:"min_idle_conns"`
	PoolSize     int           `mapstructure:"pool_size"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
}

type MQConfig struct {
	Driver      string            `mapstructure:"driver"`
	KafkaOption KafkaOptionConfig `mapstructure:"kafka_option"`
}

type KafkaOptionConfig struct {
	Brokers        []string `mapstructure:"brokers"`
	ConsumerGroup  string   `mapstructure:"consumer_group"`
	OffsetsInitial int64    `mapstructure:"offsets_initial"`
}

func NewConfig(configPath string) (Config, error) {
	var file *os.File
	file, _ = os.Open(configPath)

	v := viper.New()
	v.SetConfigType("yaml")
	v.AutomaticEnv()

	/* default */
	v.SetDefault("release", false)
	v.SetDefault("logger.level", "INFO")
	v.SetDefault("logger.format", logger.ConsoleFormat)
	v.SetDefault("http.port", "8080")

	/* database */
	v.SetDefault("database.driver", "mysql")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 3306)
	v.SetDefault("database.database", "")
	v.SetDefault("database.instance_name", "")
	v.SetDefault("database.user", "root")
	v.SetDefault("database.password", "")
	v.SetDefault("database.connect_timeout", "10s")
	v.SetDefault("database.read_timeout", "30s")
	v.SetDefault("database.write_timeout", "30s")
	v.SetDefault("database.dial_timeout", "10s")
	v.SetDefault("database.max_idletime", "1h")
	v.SetDefault("database.max_lifetime", "1h")
	v.SetDefault("database.max_idle_conn", 2)
	v.SetDefault("database.max_open_conn", 5)
	v.SetDefault("database.ssl_mode", false)

	/* redis */
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.pool_size", 10)
	v.SetDefault("redis.min_idle_conns", 5)
	v.SetDefault("redis.dial_timeout", "1s")

	/* mq */
	v.SetDefault("mq.driver", "")
	v.SetDefault("mq.kafka_option.brokers", []string{})
	v.SetDefault("mq.kafka_option.consumer_group", "")
	v.SetDefault("mq.kafka_option.offsets_initial", -2) // OffsetNewest = -1 ,OffsetOldest = -2

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.ReadConfig(file)

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}
