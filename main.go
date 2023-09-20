package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dop251/goja"
	"github.com/draganm/manifestor/interpolate"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) (err error) {

			fmt.Println("x")
			manifestorDir, err := findDotManifestorDir("")
			if err != nil {
				return err
			}

			manifestorJS := filepath.Join(manifestorDir, "manifestor.js")

			mfjsd, err := os.ReadFile(manifestorJS)
			if err != nil {
				return fmt.Errorf("could not read manifestor.js: %w", err)
			}

			encoder := yaml.NewEncoder(os.Stdout)

			vm := goja.New()
			vm.GlobalObject().Set("render", func(name string, values map[string]any) error {
				templateName := filepath.Join(manifestorDir, "templates", name)
				td, err := os.ReadFile(templateName)
				if err != nil {
					return fmt.Errorf("could not read template %s: %w", name, err)
				}
				err = interpolate.Interpolate(string(td), templateName, values, encoder)
				if err != nil {
					return fmt.Errorf("file %s: %w", name, err)
				}
				return nil
			})

			_, err = vm.RunScript("manifestor.js", string(mfjsd))

			if err != nil {
				return fmt.Errorf("could not generate manifests: %w", err)
			}
			return nil
		},
		Commands: []*cli.Command{
			&cli.Command{
				Name: "init",
				Action: func(c *cli.Context) error {
					return nil
				},
			},
		},
	}
	app.RunAndExitOnError()
}

func findDotManifestorDir(path string) (string, error) {
	fullPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("could not get full path of %s: %w", path, err)
	}

	for {

		dotManifestor := filepath.Join(fullPath, ".manifestor")
		st, err := os.Stat(dotManifestor)

		if os.IsNotExist(err) {
			parent := filepath.Dir(fullPath)

			fmt.Println(parent, fullPath)

			if parent == fullPath {
				return "", fmt.Errorf("could not find .manifestor directory in any parent of %s", fullPath)
			}

			fullPath = parent
			continue

		}

		if err != nil {
			return "", err
		}

		if st.IsDir() {
			return dotManifestor, nil
		}

	}
}
