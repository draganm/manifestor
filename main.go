package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/drone/envsubst"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

func main() {
	app := &cli.App{
		Action: func(c *cli.Context) (err error) {
			defer func() {
				if err != nil {
					err = cli.NewExitError(err.Error(), 1)
				}
			}()

			enc := yaml.NewEncoder(os.Stdout)

			for _, f := range c.Args().Slice() {
				var b []byte
				if f == "-" {
					b, err = io.ReadAll(os.Stdin)
					if err != nil {
						return fmt.Errorf("while reading from stdin: %w", err)
					}
				} else {
					b, err = os.ReadFile(f)
					if err != nil {
						return fmt.Errorf("while reading file %s: %w", f, err)
					}

				}

				dec := yaml.NewDecoder(bytes.NewReader(b))

				for {
					var obj interface{}
					err = dec.Decode(&obj)
					if err == io.EOF {
						break
					}

					if err != nil {
						return fmt.Errorf("while decoding yaml file %s: %w", f, err)
					}

					iobj, err := interpolate(obj)
					if err != nil {
						return fmt.Errorf("while interpolating values into %s: %w", f, err)
					}

					err = enc.Encode(iobj)
					if err != nil {
						return fmt.Errorf("while encoding interpolated manifests: %w", err)
					}

				}

			}
			return nil
		},
	}
	app.Run(os.Args)
}

func interpolate(o interface{}) (interface{}, error) {

	switch v := o.(type) {
	case map[string]interface{}:
		nm := make(map[string]interface{}, len(v))
		for mk, mv := range v {
			nv, err := interpolate(mv)
			if err != nil {
				return nil, fmt.Errorf("while interpolating map value for key %q: %w", mk, err)
			}
			nm[mk] = nv
		}
		return nm, nil
	case []interface{}:
		ns := make([]interface{}, len(v))
		for i, sv := range v {
			nv, err := interpolate(sv)
			if err != nil {
				return nil, fmt.Errorf("while interpolating slice value with index %d: %w", i, err)
			}
			ns[i] = nv
		}
		return ns, nil
	case string:
		return envsubst.EvalEnv(v)
	default:
		return o, nil
	}
}
