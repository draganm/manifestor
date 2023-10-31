package interpolate

import (
	"fmt"
	"io"
	"strings"

	"github.com/dop251/goja"
	"github.com/drone/envsubst"
	"gopkg.in/yaml.v3"
)

func Interpolate(template string, fileName string, values map[string]any, encoder *yaml.Encoder) error {

	dec := yaml.NewDecoder(strings.NewReader(template))
	for {
		doc := &yaml.Node{}
		err := dec.Decode(doc)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("could not parse template: %w", err)
		}

		err = interpolate(doc, values)
		if err != nil {
			return fmt.Errorf("could not interpolate: %w", err)
		}

		if fileName != "" {
			doc.HeadComment = fmt.Sprintf("file %s", fileName)
		}

		encoder.Encode(doc)
	}

	return nil
}

func getMapValue(n *yaml.Node, key string) *yaml.Node {
	for i, c := range n.Content {
		if c.Kind == yaml.ScalarNode && c.Tag == "!!str" && c.Value == key {
			if (i + 1) < len(n.Content) {
				return n.Content[i+1]
			}
		}
	}
	return nil
}

func getValue(src string, values map[string]any) (any, error) {
	rt := goja.New()
	for k, v := range values {
		err := rt.GlobalObject().Set(k, v)
		if err != nil {
			return nil, fmt.Errorf("could not set global %s: %w", k, err)
		}
	}

	pr, err := goja.Compile("", src, true)
	if err != nil {
		return nil, fmt.Errorf("could not compile %q: %w", src, err)
	}

	v, err := rt.RunProgram(pr)
	if err != nil {
		return nil, fmt.Errorf("interpolation error: %w", err)
	}

	return v.Export(), nil

}

func withValue(m map[string]any, k string, v any) map[string]any {
	ww := map[string]any{}
	for k, v := range m {
		ww[k] = v
	}

	ww[k] = v

	return ww
}

func interpolate(n *yaml.Node, values map[string]any) (err error) {
	defer func() {
		if err != nil {
			if !(strings.Contains(err.Error(), "line ") && strings.Contains(err.Error(), "col ")) {
				err = fmt.Errorf("line %d, col %d: %w", n.Line, n.Column, err)
			}
		}
	}()
	switch n.Kind {
	case yaml.MappingNode:
		forEach := getMapValue(n, "_forEach")
		if forEach != nil {

			template := getMapValue(n, "_template")
			if template == nil {
				return fmt.Errorf("could not find _template for _forEach")
			}

			v, err := getValue(forEach.Value[2:len(forEach.Value)-1], values)
			if err != nil {
				return fmt.Errorf("could not get forEach values: %w", err)
			}

			vals, isSlice := v.([]any)
			if !isSlice {
				return fmt.Errorf("forEach expression did not return an array")
			}

			children := []*yaml.Node{}
			for i, v := range vals {
				tc := cloneNode(template)
				ww := withValue(values, "eachValue", v)
				err = interpolate(tc, ww)
				if err != nil {
					return fmt.Errorf("while interpolating each %d: %w", i, err)
				}

				children = append(children, tc)

			}

			n.Kind = yaml.SequenceNode
			n.Tag = "!!seq"
			n.Content = children

			return nil

		}

		for _, c := range n.Content {
			err = interpolate(c, values)
			if err != nil {
				return err
			}
		}

		return nil
		// spew.Dump(n)
		// for _, c := range n.Content {
		// 	panic(c.Value)
		// }
		// if n.Tag != "!!map" {
		// 	panic(n.Tag)
		// }
	case yaml.ScalarNode:

		if n.Tag != "!!str" {
			return nil
		}

		if strings.HasPrefix(n.Value, "${") && strings.HasSuffix(n.Value, "}") {
			src := n.Value[2 : len(n.Value)-1]
			val, err := getValue(src, values)
			if err != nil {
				return err
			}
			d, err := yaml.Marshal(val)
			if err != nil {
				return fmt.Errorf("could not marshal value %v: %w", val, err)
			}

			vn := &yaml.Node{}
			err = yaml.Unmarshal(d, vn)
			if err != nil {
				return fmt.Errorf("could not unmarshal value: %w", err)
			}

			valueNode := vn.Content[0]

			*n = *valueNode
			return nil
		}

		n.Value, err = envsubst.Eval(n.Value, func(s string) string {
			return fmt.Sprintf("%v", values[s])
		})

		if err != nil {
			return fmt.Errorf("could not interpolate %d:%d: %w", n.Line, n.Column, err)
		}
	default:
		for _, c := range n.Content {
			err = interpolate(c, values)
			if err != nil {
				return err
			}
		}

	}
	return nil

}

func cloneNode(o *yaml.Node) *yaml.Node {

	contentClone := []*yaml.Node{}
	for _, c := range o.Content {
		contentClone = append(contentClone, cloneNode(c))
	}
	c := &yaml.Node{
		Kind:        o.Kind,
		Style:       o.Style,
		Tag:         o.Tag,
		Value:       o.Value,
		Anchor:      o.Anchor,
		Alias:       o.Alias,
		Content:     contentClone,
		HeadComment: o.HeadComment,
		LineComment: o.LineComment,
		FootComment: o.FootComment,
	}

	return c
}
