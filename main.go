package main

import (
	"bytes"
	"io"
	"os"

	"github.com/drone/envsubst"
	"github.com/pkg/errors"
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
						return errors.Wrap(err, "while reading from stdin")
					}
				} else {
					b, err = os.ReadFile(f)
					if err != nil {
						return errors.Wrapf(err, "while reading file %s", f)
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
						return errors.Wrapf(err, "while decoding yaml file %s", f)
					}

					iobj, err := interpolate(obj)
					if err != nil {
						return errors.Wrapf(err, "while interpolating values into %s", f)
					}

					err = enc.Encode(iobj)
					if err != nil {
						return errors.Wrap(err, "while encoding interpolated manifests")
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
				return nil, errors.Wrapf(err, "while interpolating map value for key %q", mk)
			}
			nm[mk] = nv
		}
		return nm, nil
	case []interface{}:
		ns := make([]interface{}, len(v))
		for i, sv := range v {
			nv, err := interpolate(sv)
			if err != nil {
				return nil, errors.Wrapf(err, "while interpolating slice value with index %d", i)
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
