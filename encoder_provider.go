package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func manifestorProvider(outDir string) func(name string) (*yaml.Encoder, func() error, error) {

	if outDir == "" {
		return func(name string) (*yaml.Encoder, func() error, error) {
			return yaml.NewEncoder(os.Stdout), func() error { return nil }, nil
		}
	}

	return func(name string) (*yaml.Encoder, func() error, error) {
		outPath := filepath.Join(outDir, name)
		dir := filepath.Dir(outPath)
		if dir != "" {
			_, err := os.Stat(dir)
			if os.IsNotExist(err) {
				err = os.MkdirAll(dir, 0777)
				if err != nil {
					return nil, nil, fmt.Errorf("could not mkdir %s: %w", dir, err)
				}
				_, err = os.Stat(dir)
			}

			if err != nil {
				return nil, nil, fmt.Errorf("could not stat %s: %w", dir, err)
			}
		}

		f, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)

		if err != nil {
			return nil, nil, fmt.Errorf("could not open output file: %w", err)
		}
		return yaml.NewEncoder(f), f.Close, nil

	}
}
