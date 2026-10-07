package yamlparse

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_keysAreCaseSensitive(t *testing.T) {
	yamlText := `
  Port: titlecase
  port: lowercase
  PORT: uppercase
  PoRt: togglecase
  PorT: mixedcase
`

	testCases := map[string]string{
		"Port": "titlecase",
		"port": "lowercase",
		"PORT": "uppercase",
		"PoRt": "togglecase",
		"PorT": "mixedcase",
	}

	returned := AssertValidYamlMapping(t, yamlText)

	for key, expected := range testCases {
		t.Run(key, func(t *testing.T) {
			anyValue, exists := returned[key]

			require.True(t, exists)

			value, ok := anyValue.(string)

			require.True(t, ok)
			require.Equal(t, expected, value)
		})
	}
}

func Test_complexKeyMappings(t *testing.T) {
	t.Run("just-typecast-example", func(t *testing.T) {
		yamlText := `
? !!str abc
: alphabetic
? !!str 123
: numeric
`
		mapping := map[string]string{
			"abc": "alphabetic",
			"123": "numeric",
		}

		returned := AssertValidYamlMapping(t, yamlText)

		for key, expected := range mapping {
			t.Run(key, func(t *testing.T) {
				anyValue, exists := returned[key]

				require.True(t, exists)

				value, ok := anyValue.(string)

				require.True(t, ok)
				require.Equal(t, expected, value)
			})
		}
	})

	t.Run("list-as-a-key", func(t *testing.T) {
		// It is valid YAML, but Go is not able to process it. Unfortunatelly.

		yamlText := `
? - abc
  - xyz
: value of a list key
`

		decoded, err := ParseYaml(yamlText)

		require.ErrorContains(t, err, `yaml: invalid map key: []interface {}{"abc", "xyz"}`)

		yamlRoot, ok := decoded.(map[any]any)

		require.True(t, ok)
		require.Len(t, yamlRoot, 0)
	})
}
