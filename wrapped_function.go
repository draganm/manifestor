package main

import (
	"fmt"

	"github.com/dop251/goja"
)

func wrap(name string, f func(interface{}) error) func(interface{}) error {
	return func(i interface{}) (err error) {
		defer func() {
			p := recover()
			if p == nil {
				return
			}

			ex, isException := p.(*goja.Exception)

			if isException {
				err = fmt.Errorf("while executing %s():\n%v", name, ex.String())
			} else {
				err = fmt.Errorf("while executing %s: %v", name, p)
			}

		}()

		return f(i)
	}
}
