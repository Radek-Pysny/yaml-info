package yamlparse

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Test_sequence_atTopLevel presents very basic YAML with just a sequence on top-level.
func Test_sequence_atTopLevel(t *testing.T) {
	testCases := []struct {
		title    string
		yamlText string
		expected []string
	}{
		{
			title:    `seq-with-single-item`,
			yamlText: `- item`,
			expected: []string{`item`},
		},
		{
			title:    `seq-with-single-item-in-flow-style`,
			yamlText: `[item]`,
			expected: []string{`item`},
		},
		{
			title: `seq-with-multiple-items`,
			yamlText: `- 1st item
- 2nd item
- 3rd item`,
			expected: []string{`1st item`, `2nd item`, `3rd item`},
		},
		{
			title:    `seq-with-multiple-items-in-flow-style-on-one-line`,
			yamlText: `[1st item, 2nd item, 3rd item]`,
			expected: []string{`1st item`, `2nd item`, `3rd item`},
		},
		{
			title: `seq-with-multiple-items-in-flow-style-on-two-lines`,
			yamlText: `[1st item, 2nd item, 
  3rd item]`,
			expected: []string{`1st item`, `2nd item`, `3rd item`},
		},
		{
			title: `seq-with-multiple-items-in-flow-style-each-item-on-separate-line-1`,
			yamlText: `[
  1st item,
  2nd item, 
  3rd item # no extra trailing comma after the last item
  ]`,
			expected: []string{`1st item`, `2nd item`, `3rd item`},
		},
		{
			title: `seq-with-multiple-items-in-flow-style-each-item-on-separate-line-2`,
			yamlText: `[
  1st item,
  2nd item, 
  3rd item, # extra trailing comma after the last item
  ] # mind that closing bracket has to be indented as well`,
			expected: []string{`1st item`, `2nd item`, `3rd item`},
		},
		{
			title: `seq-with-multiple-items-in-flow-style-each-item-on-separate-line-3`,
			yamlText: `[
    1st item, # bigger indentation for items than for closing bracket
    2nd item, 
    3rd item,
  ] # mind that closing bracket has to be indented`,
			expected: []string{`1st item`, `2nd item`, `3rd item`},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			returned := AssertValidYamlSequenceOfText(t, tc.yamlText)

			require.Equal(t, tc.expected, returned)
		})
	}
}

// Test_sequence_underTopLevelMapping presents more common usage when some sequence under
// a top-level mapping. It is just a slight extension of Test_sequence_atTopLevel.
func Test_sequence_underTopLevelMapping(t *testing.T) {
	testCases := []struct {
		title    string
		key      string
		yamlText string
		expected []string
	}{
		{
			title: `seq-with-single-item`,
			key:   `go`,
			yamlText: `go:
  - v1.27`,
			expected: []string{`v1.27`},
		},
		{
			title:    `seq-with-single-item-in-flow-style`,
			key:      `go`,
			yamlText: `go: [v1.27]`,
			expected: []string{`v1.27`},
		},
		{
			title: `seq-with-multiple-items`,
			key:   `go`,
			yamlText: `go:
  - v1.27
  - v1.26
  - v1.25
  - v1.24
  - v1.23
  - v1.22
  - v1.21
  - v1.20
  - v1.19`,
			expected: []string{`v1.27`, `v1.26`, `v1.25`, `v1.24`, `v1.23`, `v1.22`, `v1.21`, `v1.20`, `v1.19`},
		},
		{
			title:    `seq-with-multiple-items-in-flow-style-on-one-line`,
			key:      `go`,
			yamlText: `go: [v1.27, v1.26, v1.25, v1.24, v1.23, v1.22, v1.21, v1.20, v1.19]`,
			expected: []string{`v1.27`, `v1.26`, `v1.25`, `v1.24`, `v1.23`, `v1.22`, `v1.21`, `v1.20`, `v1.19`},
		},
		{
			title: `seq-with-multiple-items-in-flow-style-on-two-lines`,
			key:   `go`,
			yamlText: `go: [v1.27, v1.26, v1.25, v1.24, v1.23, v1.22, v1.21,
  v1.20, v1.19]`,
			expected: []string{`v1.27`, `v1.26`, `v1.25`, `v1.24`, `v1.23`, `v1.22`, `v1.21`, `v1.20`, `v1.19`},
		},
		{
			title: `seq-with-multiple-items-in-flow-style-each-item-on-separate-line-1`,
			key:   `go`,
			yamlText: `go: [
  v1.27,
  v1.26,
  v1.25,
  v1.24,
  v1.23,
  v1.22,
  v1.21,
  v1.20,
  v1.19
  ]`,
			expected: []string{`v1.27`, `v1.26`, `v1.25`, `v1.24`, `v1.23`, `v1.22`, `v1.21`, `v1.20`, `v1.19`},
		},
		{
			title: `seq-with-multiple-items-in-flow-style-each-item-on-separate-line-2`,
			key:   `go`,
			yamlText: `go:
  [
  v1.27,
  v1.26,
  v1.25,
  v1.24,
  v1.23,
  v1.22,
  v1.21,
  v1.20,
  v1.19
  ]`,
			expected: []string{`v1.27`, `v1.26`, `v1.25`, `v1.24`, `v1.23`, `v1.22`, `v1.21`, `v1.20`, `v1.19`},
		},
		{
			title: `seq-with-multiple-items-in-flow-style-each-item-on-separate-line-3`,
			key:   `go`,
			yamlText: `go: [
  v1.27,
  v1.26,
  v1.25,
  v1.24,
  v1.23,
  v1.22,
  v1.21,
  v1.20,
  v1.19,
  ]`,
			expected: []string{`v1.27`, `v1.26`, `v1.25`, `v1.24`, `v1.23`, `v1.22`, `v1.21`, `v1.20`, `v1.19`},
		},
		{
			title: `seq-with-multiple-items-in-flow-style-each-item-on-separate-line-4`,
			key:   `go`,
			yamlText: `go:
  [
  v1.27,
  v1.26,
  v1.25,
  v1.24,
  v1.23,
  v1.22,
  v1.21,
  v1.20,
  v1.19,
  ]`,
			expected: []string{`v1.27`, `v1.26`, `v1.25`, `v1.24`, `v1.23`, `v1.22`, `v1.21`, `v1.20`, `v1.19`},
		},
		{
			title: `seq-with-multiple-items-in-flow-style-each-item-on-separate-line-5`,
			key:   `go`,
			yamlText: `go: [
    v1.27,
    v1.26,
    v1.25,
    v1.24,
    v1.23,
    v1.22,
    v1.21,
    v1.20,
    v1.19
  ]`,
			expected: []string{`v1.27`, `v1.26`, `v1.25`, `v1.24`, `v1.23`, `v1.22`, `v1.21`, `v1.20`, `v1.19`},
		},
		{
			title: `seq-with-multiple-items-in-flow-style-each-item-on-separate-line-6`,
			key:   `go`,
			yamlText: `go:
  [
    v1.27,
    v1.26,
    v1.25,
    v1.24,
    v1.23,
    v1.22,
    v1.21,
    v1.20,
    v1.19
  ]`,
			expected: []string{`v1.27`, `v1.26`, `v1.25`, `v1.24`, `v1.23`, `v1.22`, `v1.21`, `v1.20`, `v1.19`},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			mapping := AssertValidYamlMapping(t, tc.yamlText)

			seq, found := mapping[tc.key]

			require.True(t, found)

			returned := ConvertSequenceOfText(t, seq)

			require.Equal(t, tc.expected, returned)
		})
	}
}

// Test_sequence_nested presents small examples of nested sequences using multiple
// YAML representation but having the same meaning (being parsed into exactly the
// same structure).
func Test_sequence_nested(t *testing.T) {
	testCases := []struct {
		title     string
		yamlTexts []string
		expected  any
	}{
		{
			title: "item-in-two-levels",
			yamlTexts: []string{
				`sequence:
- level one
- - level two`,
				`sequence:
- level one
- [level two]`,
				`sequence: [level one, [level two]]`,
			},
			expected: map[string]any{
				"sequence": []any{"level one", []any{"level two"}},
			},
		},
		{
			title: "item-in-two-levels",
			yamlTexts: []string{
				`dice rolls:
- - 2
  - 5
- - 3
  - 4`,
				`dice rolls:
- [2, 5]
- [3, 4]`,
				`dice rolls: [ [2, 5], [3, 4] ]`,
			},
			expected: map[string]any{
				"dice rolls": []any{[]any{2, 5}, []any{3, 4}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			require.Greater(t, len(tc.yamlTexts), 0)

			for _, yamlText := range tc.yamlTexts {
				returned := AssertValidYamlMapping(t, yamlText)

				require.Equal(t, tc.expected, returned)
			}
		})
	}
}
