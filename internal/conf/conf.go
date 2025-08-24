package conf

import "log/slog"

var (
	Conf = new(Config)
)

type Config struct {
	EnableCORS bool
	Proxy      []Proxy `yaml:"proxy"`
}

func (Config) Print() {
	logger := slog.Default()
	logger.Info("config", "proxy", Conf.Proxy)
}

type Proxy struct {
	Host     string `yaml:"host"`
	Upstream string `yaml:"upstream"`
}
