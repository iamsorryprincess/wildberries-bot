package config

import (
	"errors"
	"fmt"
	"os"
)

type FileProvider struct {
	filename string
}

func NewFileProvider(filename string) *FileProvider {
	return &FileProvider{
		filename: filename,
	}
}

func (p *FileProvider) Read() ([]byte, error) {
	file, err := os.OpenFile(p.filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("config file provider open file: %w", err)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		err = fmt.Errorf("config file provider stat file: %w", err)

		var cErr error
		if cErr = file.Close(); cErr != nil {
			return nil, errors.Join(err, fmt.Errorf("config file provider close file: %w", cErr))
		}

		return nil, err
	}

	data := make([]byte, fileInfo.Size())
	if _, err = file.Read(data); err != nil {
		err = fmt.Errorf("config file provider read file: %w", err)

		var cErr error
		if cErr = file.Close(); cErr != nil {
			return nil, errors.Join(err, fmt.Errorf("config file provider close file: %w", cErr))
		}

		return nil, err
	}

	if err = file.Close(); err != nil {
		return nil, fmt.Errorf("config file provider close file: %w", err)
	}

	return data, nil
}
