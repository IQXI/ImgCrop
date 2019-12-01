package structs

type LoggerConfig struct {
	Level    string `mapstructure:"level"`
	LogsPath string `mapstructure:"logs_path"`
	FileName string `mapstructure:"filename"`
	Name     string `mapstructure:"name"`
}

type CacheConfig struct {
	Path string `mapstructure:"path"`
	Size int64  `mapstructure:"size"`
}

//base config

type Config struct {
	Logger LoggerConfig `mapstructure:"logger"`
	Cache  CacheConfig  `mapstructure:"cache"`
}

type Image struct {
	Size       int64
	Path       string
	Headers    map[string]string
	Link       string
	FileName   string
	Exctension string
}

type Cache struct {
	Age   int64
	Image Image
}
