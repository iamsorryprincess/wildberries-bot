package config

import (
	"errors"
	"fmt"
	"os"
)

type FileProvider[TConfig any] struct {
	filename string
	parser   Parser[TConfig]
}

func NewFileProvider[TConfig any](filename string, parser Parser[TConfig]) *FileProvider[TConfig] {
	return &FileProvider[TConfig]{
		filename: filename,
		parser:   parser,
	}
}

func (p *FileProvider[TConfig]) Get() (TConfig, error) {
	var cfg TConfig

	file, err := os.OpenFile(p.filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return cfg, fmt.Errorf("file provider open file: %w", err)
	}

	cfg, err = p.parser.Parse(file)
	if err != nil {
		err = fmt.Errorf("file provider parse: %w", err)
		if cErr := file.Close(); cErr != nil {
			return cfg, errors.Join(err, fmt.Errorf("file provider close file: %w", cErr))
		}

		return cfg, err
	}

	if err = file.Close(); err != nil {
		return cfg, fmt.Errorf("file provider close file: %w", err)
	}

	return cfg, nil
}
