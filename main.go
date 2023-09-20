package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dop251/goja"
	"github.com/draganm/manifestor/interpolate"
	"github.com/go-git/go-git/v5"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) (err error) {

			repo, err := git.PlainOpenWithOptions("", &git.PlainOpenOptions{DetectDotGit: true})
			if err != nil {
				return fmt.Errorf("could not open git repo: %w", err)
			}

			head, err := repo.Head()
			if err != nil {
				return fmt.Errorf("could not get git head: %w", err)
			}

			fmt.Println("head commit", head.Hash().String())

			gitValues := map[string]string{
				"headSha":      head.Hash().String(),
				"headShaShort": head.Hash().String()[:7],
			}

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

			env := map[string]string{}

			for _, ev := range os.Environ() {
				name, value, found := strings.Cut(ev, "=")
				if found {
					env[name] = value
				}
			}

			vm.GlobalObject().Set("env", env)
			vm.GlobalObject().Set("git", gitValues)

			vm.GlobalObject().Set("render", func(name string, values map[string]any, fileName string) error {
				templateName := filepath.Join(manifestorDir, "templates", name)
				td, err := os.ReadFile(templateName)
				if err != nil {
					return fmt.Errorf("could not read template %s: %w", name, err)
				}
				err = interpolate.Interpolate(string(td), fileName, values, encoder)
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
