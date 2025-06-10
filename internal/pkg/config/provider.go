package config

import "flag"

type Reader interface {
	Read() ([]byte, error)
}

type Parser[TConfig any] interface {
	Parse(data []byte) (TConfig, error)
}

type SetDefaults func()

func ReadConfig[TConfig any](reader Reader, parser Parser[TConfig], defaults ...SetDefaults) (TConfig, error) {
	var cfg TConfig

	data, err := reader.Read()
	if err != nil {
		return cfg, err
	}

	for _, setDefaultFunc := range defaults {
		setDefaultFunc()
	}

	cfg, err = parser.Parse(data)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func Load[TConfig any](defaults ...SetDefaults) (TConfig, error) {
	filename := flag.String("config", "config.yaml", "config file path")
	flag.Parse()
	return ReadConfig(NewFileProvider(*filename), NewYamlParser[TConfig](), defaults...)
}
