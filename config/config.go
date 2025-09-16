package config

import (
	"flag"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const (
	TrueStr = "true"
)

type Env string

var GlobalConfig *Config
var configMutex sync.RWMutex

type Config struct {
        Env           Env               `yaml:"env" mapstructure:"env"`
    	App           *AppConfig        `yaml:"app" mapstructure:"app"`
    	HTTPServer    *HttpServerConfig `yaml:"http_server" mapstructure:"http_server"`
    	MetricsServer *MetricsConfig    `yaml:"metrics_server" mapstructure:"metrics_server"`
    	Postgres      *PostgreSQLConfig `yaml:"postgres" mapstructure:"postgres"`
        Log           *LogConfig        `yaml:"log" mapstructure:"log"`
}
type AppConfig struct {
	Name    string `yaml:"name" mapstructure:"name"`
}

type HttpServerConfig struct {
	Address         string `yaml:"address" mapstructure:"address"`
	Timeout         string `yaml:"timeout" mapstructure:"read_timeout"`
	IdleTimeout     string `yaml:"idle_timeout" mapstructure:"write_timeout"`
}

type PostgreSQLConfig struct {
	User            string `yaml:"user" mapstructure:"user"`
	Password        string `yaml:"password" mapstructure:"password"`
	Host            string `yaml:"host" mapstructure:"host"`
	Port            int    `yaml:"port" mapstructure:"port"`
	Database        string `yaml:"database" mapstructure:"database"`
	IdleTimeout     int    `yaml:"idle_timeout" mapstructure:"idle_timeout"`
	ConnectTimeout  int    `yaml:"connect_timeout" mapstructure:"connect_timeout"`
}

type MetricsConfig struct {
	Address    string `yaml:"address" mapstructure:"address"`
	Enabled bool   `yaml:"enabled" mapstructure:"enabled"`
	Path    string `yaml:"path" mapstructure:"path"`
}

type LogConfig struct {
	SavePath         string `yaml:"save_path" mapstructure:"save_path"`
	FileName         string `yaml:"file_name" mapstructure:"file_name"`
	MaxSize          int    `yaml:"max_size" mapstructure:"max_size"`
	MaxAge           int    `yaml:"max_age" mapstructure:"max_age"`
	LocalTime        bool   `yaml:"local_time" mapstructure:"local_time"`
	Compress         bool   `yaml:"compress" mapstructure:"compress"`
	Level            string `yaml:"level" mapstructure:"level"`
	EnableConsole    bool   `yaml:"enable_console" mapstructure:"enable_console"`
	EnableColor      bool   `yaml:"enable_color" mapstructure:"enable_color"`
	EnableCaller     bool   `yaml:"enable_caller" mapstructure:"enable_caller"`
	EnableStacktrace bool   `yaml:"enable_stacktrace" mapstructure:"enable_stacktrace"`
}

func Load(configPath string, configFile string) (*Config, error) {
	var conf *Config
	vip := viper.New()
	vip.AddConfigPath(configPath)
	vip.SetConfigName(configFile)

	vip.SetConfigType("yaml")
	if err := vip.ReadInConfig(); err != nil {
		return nil, err
	}

	// Enable environment variables to override config
	vip.SetEnvPrefix("APP")
	vip.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	vip.AutomaticEnv()

	err := vip.Unmarshal(&conf)
	if err != nil {
		return nil, err
	}

	// Apply environment variable overrides
	applyEnvOverrides(conf)

	// Setup config file change monitoring
	vip.WatchConfig()
	vip.OnConfigChange(func(e fsnotify.Event) {
		// Reload configuration when file changes
		var newConf Config
		if err := vip.Unmarshal(&newConf); err == nil {
			// Apply environment variable overrides to the new config
			applyEnvOverrides(&newConf)

			// Update global config with new values - with mutex protection
			configMutex.Lock()
			*GlobalConfig = newConf
			configMutex.Unlock()
		}
	})

	return conf, nil
}

func applyEnvOverrides(conf *Config) {
	// Apply config overrides by category
	applyAppEnvOverrides(conf)
	applyHTTPServerEnvOverrides(conf)
	applyMetricsServerEnvOverrides(conf)
	applyPostgresEnvOverrides(conf)
    applyLogEnvOverrides(conf)

}

func applyAppEnvOverrides(conf *Config) {
	// Environment
	if env := os.Getenv("APP_ENV"); env != "" {
		conf.Env = Env(env)
	}
	// App config
	if name := os.Getenv("APP_APP_NAME"); name != "" {
		conf.App.Name = name
	}
}

func applyHTTPServerEnvOverrides(conf *Config) {
	if address := os.Getenv("APP_HTTP_SERVER_ADDRESS"); address != "" {
		conf.HTTPServer.Address = address
	}
	if timeout := os.Getenv("APP_HTTP_SERVER_TIMEOUT"); timeout != "" {
		conf.HTTPServer.Timeout = timeout
	}
	if idle_timeout := os.Getenv("APP_HTTP_SERVER_IDLE_TIMEOUT"); idle_timeout != "" {
		conf.HTTPServer.IdleTimeout = idle_timeout
	}
}

func applyMetricsServerEnvOverrides(conf *Config) {
	// Initialize MetricsServer if it doesn't exist
	if conf.MetricsServer == nil {
		conf.MetricsServer = &MetricsConfig{
			Address:    ":9090",
			Enabled: true,
			Path:    "/metrics",
		}
	}

	if address := os.Getenv("APP_METRICS_SERVER_ADDR"); address != "" {
		conf.MetricsServer.Address = address
	}
	if enabled := os.Getenv("APP_METRICS_SERVER_ENABLED"); enabled != "" {
		conf.MetricsServer.Enabled = enabled == TrueStr
	}
	if path := os.Getenv("APP_METRICS_SERVER_PATH"); path != "" {
		conf.MetricsServer.Path = path
	}
}

func applyPostgresEnvOverrides(conf *Config) {
	if host := os.Getenv("APP_POSTGRES_HOST"); host != "" {
		conf.Postgres.Host = host
	}
	if port := os.Getenv("APP_POSTGRES_PORT"); port != "" {
		if val, err := strconv.Atoi(port); err == nil {
			conf.Postgres.Port = val
		}
	}
	if username := os.Getenv("APP_POSTGRES_USERNAME"); username != "" {
		conf.Postgres.User = username
	}
	if password := os.Getenv("APP_POSTGRES_PASSWORD"); password != "" {
		conf.Postgres.Password = password
	}
	if database := os.Getenv("APP_POSTGRES_DB_NAME"); database != "" {
		conf.Postgres.Database = database
	}
	if idleTimeout := os.Getenv("APP_POSTGRES_IDLE_TIMEOUT"); idleTimeout != "" {
		if val, err := strconv.Atoi(idleTimeout); err == nil {
			conf.Postgres.IdleTimeout = val
		}
	}
	if connectTimeout := os.Getenv("APP_POSTGRES_CONNECT_TIMEOUT"); connectTimeout != "" {
		if val, err := strconv.Atoi(connectTimeout); err == nil {
			conf.Postgres.ConnectTimeout = val
		}
	}
}

func applyLogEnvOverrides(conf *Config) {
	if savePath := os.Getenv("APP_LOG_SAVE_PATH"); savePath != "" {
		conf.Log.SavePath = savePath
	}
	if fileName := os.Getenv("APP_LOG_FILE_NAME"); fileName != "" {
		conf.Log.FileName = fileName
	}
	if maxSize := os.Getenv("APP_LOG_MAX_SIZE"); maxSize != "" {
		if val, err := strconv.Atoi(maxSize); err == nil {
			conf.Log.MaxSize = val
		}
	}
	if maxAge := os.Getenv("APP_LOG_MAX_AGE"); maxAge != "" {
		if val, err := strconv.Atoi(maxAge); err == nil {
			conf.Log.MaxAge = val
		}
	}
	if localTime := os.Getenv("APP_LOG_LOCAL_TIME"); localTime != "" {
		conf.Log.LocalTime = localTime == TrueStr
	}
	if compress := os.Getenv("APP_LOG_COMPRESS"); compress != "" {
		conf.Log.Compress = compress == TrueStr
	}
	if level := os.Getenv("APP_LOG_LEVEL"); level != "" {
		conf.Log.Level = level
	}
	if enableConsole := os.Getenv("APP_LOG_ENABLE_CONSOLE"); enableConsole != "" {
		conf.Log.EnableConsole = enableConsole == TrueStr
	}
	if enableColor := os.Getenv("APP_LOG_ENABLE_COLOR"); enableColor != "" {
		conf.Log.EnableColor = enableColor == TrueStr
	}
	if enableCaller := os.Getenv("APP_LOG_ENABLE_CALLER"); enableCaller != "" {
		conf.Log.EnableCaller = enableCaller == TrueStr
	}
	if enableStacktrace := os.Getenv("APP_LOG_ENABLE_STACKTRACE"); enableStacktrace != "" {
		conf.Log.EnableStacktrace = enableStacktrace == TrueStr
	}
}

func Init(path, file string) {
	configPath := flag.String("config-path", path, "path to configuration path")
	configFile := flag.String("config-file", file, "name of configuration file (without extension)")
	flag.Parse()

	conf, err := Load(*configPath, *configFile)
	if err != nil {
		panic("Load config fail : " + err.Error())
	}
	GlobalConfig = conf
}