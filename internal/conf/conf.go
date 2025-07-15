package conf

var (
	Conf = new(Config)
)

type Config struct{}

func (Config) Print() {}
