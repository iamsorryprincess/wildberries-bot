package config

import "io"

type Parser[TConfig any] interface {
	Parse(reader io.Reader) (TConfig, error)
}

type Provider[TConfig any] interface {
	Get() (TConfig, error)
}

func Load[TConfig any]() (TConfig, error) {
	return NewFileProvider[TConfig]("config.json", NewJSONParser[TConfig]()).Get()
}
