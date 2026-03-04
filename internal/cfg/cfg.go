package cfg

type (
	Config struct {
		Server Server
	}

	Server struct {
		Port string
	}
)

func NewConfig() (*Config, error) {
	var cfg Config

	cfg.Server.Port = ":8080"

	return &cfg, nil
}
