package yamlparse

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v3"
)

// AssertValidYamlSequenceOfText is more concrete variant of AssertValidYamlSequence,
// as it returns slice of strings. Most of unit tests are OK with that simple usage.
func AssertValidYamlSequenceOfText(t *testing.T, yamlText string) []string {
	t.Helper()

	seq := AssertValidYamlSequence(t, yamlText)
	return ConvertSequenceOfText(t, seq)
}

// ConvertSequenceOfText converts output of AssertValidYamlSequence into slice of strings.
// It might fail if the original YAML text was not parsed into just strings.
func ConvertSequenceOfText(t *testing.T, anySeq any) []string {
	t.Helper()

	seq, ok := anySeq.([]any)

	require.True(t, ok)

	result := make([]string, len(seq))
	for i, anyItem := range seq {
		item, ok := anyItem.(string)

		require.True(t, ok)

		result[i] = item
	}

	return result
}

// AssertValidYamlSequence parses the single top-level sequence out of the given YAML text.
// Expect fail if YAML does not consist of just a sequence. The parsed slice is returned
// for the purpose of unit tests.
func AssertValidYamlSequence(t *testing.T, yamlText string) []any {
	t.Helper()

	decoded, err := ParseYaml(yamlText)

	require.NoError(t, err)

	var result []any
	result, ok := decoded.([]any)

	require.True(t, ok)

	return result
}

// AssertValidYamlMappingKey parses one field of the top-level mapping out of the given YAML text.
// Expect fail if YAML does not consist of just a mapping or the mapping does not contain that key.
// The parsed item (as any) is returned for the purpose of unit tests.
func AssertValidYamlMappingKey(t *testing.T, yamlText string, key string) any {
	t.Helper()

	mapping := AssertValidYamlMapping(t, yamlText)

	result, exists := mapping[key]

	require.True(t, exists)

	return result
}

// AssertValidYamlMapping parses the single top-level mapping out of the given YAML text.
// Expect fail if YAML does not consist of just a mapping. The parsed map is returned
// for the purpose of unit tests.
func AssertValidYamlMapping(t *testing.T, yamlText string) map[string]any {
	t.Helper()

	decoded, err := ParseYaml(yamlText)

	require.NoError(t, err)

	var result map[string]any
	result, ok := decoded.(map[string]any)

	require.True(t, ok)

	return result
}

// AssertValidYamlString parses the single top-level string out of the given YAML text.
// Expect fail if YAML does not consist of just a string. The parsed string is returned
// for the purpose of unit tests.
func AssertValidYamlString(t *testing.T, yamlText string) string {
	t.Helper()

	decoded, err := ParseYaml(yamlText)

	require.NoError(t, err)

	var result string
	result, ok := decoded.(string)

	require.True(t, ok)

	return result
}

// ParseYaml parses a single document of the given YAML text. It returns error from decoding.
// It is used as the basic building block for most of the unit tests.
func ParseYaml(yamlText string) (any, error) {
	var decoded any
	buffer := bytes.NewBufferString(yamlText)
	err := yaml.NewDecoder(buffer).Decode(&decoded)

	return decoded, err
}
