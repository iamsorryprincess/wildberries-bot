package config

func Load[TConfig any]() (TConfig, error) {
	var cfg TConfig
	return cfg, nil
}
