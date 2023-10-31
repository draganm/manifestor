package interpolate_test

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/draganm/manifestor/interpolate"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestInterpolate(t *testing.T) {
	testCases := []struct {
		name           string
		template       string
		values         map[string]any
		fileName       string
		expectedOutput string
		expectedError  string
	}{
		{
			name:           "no values",
			template:       `foo: bar`,
			expectedOutput: `foo: bar`,
		},
		{
			name:           "simple string",
			template:       `foo: ${foo}`,
			values:         map[string]any{"foo": "bar"},
			expectedOutput: `foo: bar`,
		},
		{
			name:           "simple number",
			template:       `foo: ${foo}`,
			values:         map[string]any{"foo": 42},
			expectedOutput: `foo: 42`,
		},
		{
			name:           "null",
			template:       `foo: ${foo}`,
			values:         map[string]any{"foo": nil},
			expectedOutput: `foo: null`,
		},
		{
			name:          "not existing",
			template:      `foo: ${foo}`,
			values:        map[string]any{},
			expectedError: "could not interpolate: line 1, col 6: interpolation error: ReferenceError: foo is not defined at <eval>:1:1(0)",
		},
		{
			name:           "array value",
			template:       `foo: ${[1,2,'bar']}`,
			expectedOutput: "foo:\n    - 1\n    - 2\n    - bar",
		},
		{
			name:           "map value",
			template:       `foo: '${({"bar": "baz"})}'`,
			expectedOutput: "foo:\n    bar: baz",
		},
		{
			name:           "for each",
			template:       `foo: {_forEach: "${[1,2,3]}", _template: {foo: "${eachValue}"}}`,
			expectedOutput: "foo: [{foo: 1}, {foo: 2}, {foo: 3}]",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			bb := &bytes.Buffer{}
			enc := yaml.NewEncoder(bb)

			err := interpolate.Interpolate(tc.template, tc.fileName, tc.values, enc)

			if tc.expectedError != "" {
				require.Error(err, tc.expectedError)
			} else {
				require.NoError(err)
			}

			require.Equal(tc.expectedOutput, strings.TrimSuffix(bb.String(), "\n"))

		})
	}
}

func yamlEqual(expected, actual string) (bool, error) {

	var expectedYAMLAsInterface, actualYAMLAsInterface interface{}

	err := yaml.Unmarshal([]byte(expected), &expectedYAMLAsInterface)
	if err != nil {
		return false, fmt.Errorf("could not parse expected: %w", err)
	}

	err = yaml.Unmarshal([]byte(actual), &actualYAMLAsInterface)
	if err != nil {
		return false, fmt.Errorf("could not parse actual: %w", err)
	}

	return reflect.DeepEqual(expected, actual), nil

}
