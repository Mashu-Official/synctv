package conf

//nolint:tagliatelle
type LogConfig struct {
	Enable     bool   `env:"LOG_ENABLE"      yaml:"enable"`
	LogFormat  string `env:"LOG_FORMAT"      hc:"can be set: text | json"                                        yaml:"log_format"`
	FilePath   string `env:"LOG_FILE_PATH"   hc:"if it is a relative path, the data-dir directory will be used." yaml:"file_path"`
	MaxSize    int    `cm:"mb"               env:"LOG_MAX_SIZE"                                                  hc:"max size per log file" yaml:"max_size"`
	MaxBackups int    `env:"LOG_MAX_BACKUPS" yaml:"max_backups"`
	MaxAge     int    `env:"LOG_MAX_AGE"     yaml:"max_age"`
	Compress   bool   `env:"LOG_COMPRESS"    yaml:"compress"`
}

func DefaultLogConfig() LogConfig {
	return LogConfig{
		Enable:     true,
		LogFormat:  "text",
		FilePath:   "log/log.log",
		MaxSize:    10,
		MaxBackups: 10,
		MaxAge:     28,
		Compress:   false,
	}
}
