package config

type Config interface {
	SetDefaults()
}

func Load[TConfig Config]() (TConfig, error) {
	return LoadYAML[TConfig]("config.yaml")
}
