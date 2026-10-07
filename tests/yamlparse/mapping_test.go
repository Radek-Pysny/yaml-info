package yamlparse

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Test_mapping_differentWaysToRepresentTheSame ensures that three representations
// of different compactness of sequence of mapping are equal.
func Test_mapping_differentWaysToRepresentTheSame(t *testing.T) {
	testCases := []struct {
		title    string
		yamlText string
	}{
		{
			title: `expand`,
			yamlText: `
- x: 0
  y: 0
  z: 0
- x: 1
  y: 2
  z: 3`,
		},
		{
			title: `half-way`,
			yamlText: `
- {x: 0, y: 0, z: 0}
- {x: 1, y: 2, z: 3}`,
		},
		{
			title:    `compact`,
			yamlText: `[{x: 0, y: 0, z: 0}, {x: 1, y: 2, z: 3}]`,
		},
	}

	expected := []map[string]any{
		{"x": 0, "y": 0, "z": 0},
		{"x": 1, "y": 2, "z": 3},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			anySlice := AssertValidYamlSequence(t, tc.yamlText)

			require.Len(t, anySlice, len(expected))

			slice := make([]map[string]any, len(expected))
			for i, anyItem := range anySlice {
				item, ok := anyItem.(map[string]any)

				require.True(t, ok)

				slice[i] = item
			}

			require.Equal(t, expected, slice)
		})
	}
}

// Test_sequence_cannotStartSequenceOrMapOnTheSameLine brings in the topic of
// some invalid YAML inputs.
func Test_sequence_cannotStartSequenceOrMapOnTheSameLine(t *testing.T) {
	t.Run("nested-mapping-on-the-same-line", func(t *testing.T) {
		yamlText := `outer map: inner map: content`

		_, err := ParseYaml(yamlText)

		require.ErrorContains(t, err, `yaml: mapping values are not allowed in this context`)
	})

	t.Run("nested-sequence-on-the-same-line", func(t *testing.T) {
		yamlText := `outer map: - nested sequence`

		_, err := ParseYaml(yamlText)

		require.ErrorContains(t, err, `yaml: block sequence entries are not allowed in this context`)
	})
}
