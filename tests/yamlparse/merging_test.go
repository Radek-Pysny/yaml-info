package yamlparse

// This file shows usage of anchor `&REF`, alias `*REF`, and merge operator `<<: *REF`
// which enabled sharing of parts of YAML mappings.

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Test_anchor ensures that anchor (`&REF` where REF is name of anchor) does not
// break anything (even when it is not used anyhow).
func Test_anchor(t *testing.T) {
	t.Run("no-data", func(t *testing.T) {

		yamlText := `object: &anchor`

		returned := AssertValidYamlMapping(t, yamlText)

		anyValue, exists := returned["object"]

		require.True(t, exists)
		require.Nil(t, anyValue) // no data
	})

	t.Run("mapping-data", func(t *testing.T) {

		yamlText := `
object: &anchor
  name: the Object
  age: rather unknown
`

		returned := AssertValidYamlMapping(t, yamlText)

		anyValue, exists := returned["object"]

		require.True(t, exists)

		value, ok := anyValue.(map[string]any)

		require.True(t, ok)

		require.Equal(t, value["name"], "the Object")
		require.Equal(t, value["age"], "rather unknown")
	})
}

// Test_alias shows sharing of data from anchor (`&REF`) via simple aliases (`*REF`).
func Test_alias(t *testing.T) {
	yamlText := `
shipping-address: &address
  street: Dulce
  city: UnL
  zip: 400 01
billing-address: *address
contact-address: *address
`
	address := map[string]any{
		"street": "Dulce",
		"city":   "UnL",
		"zip":    "400 01",
	}

	returned := AssertValidYamlMapping(t, yamlText)

	require.Equal(t, address, returned["shipping-address"])
	require.Equal(t, address, returned["billing-address"])
	require.Equal(t, address, returned["contact-address"])
}

// Test_merging shows usage of merge operator (`<<: *REF`), but also
// extension (adding a brand-new fields) and override (change of field
// from referenced data).
func Test_merging(t *testing.T) {
	t.Run("just-merge", func(t *testing.T) {
		yamlText := `
default: &default
  timeout: 5s
  retries: 3
  endpoint: example.com

service-connector:
  name: connector
  config:
    <<: *default
`
		connector := map[string]any{
			"name": "connector",
			"config": map[string]any{
				"timeout":  "5s",
				"retries":  3,
				"endpoint": "example.com",
			},
		}

		returned := AssertValidYamlMapping(t, yamlText)

		require.Equal(t, connector, returned["service-connector"])
	})

	t.Run("merge-and-extend", func(t *testing.T) {
		yamlText := `
default: &default
  timeout: 5s
  retries: 3
  endpoint: example.com

service-connector:
  name: connector
  config:
    <<: *default
    build-flags: production, no-db
`
		connector := map[string]any{
			"name": "connector",
			"config": map[string]any{
				"timeout":     "5s",
				"retries":     3,
				"endpoint":    "example.com",
				"build-flags": "production, no-db",
			},
		}

		returned := AssertValidYamlMapping(t, yamlText)

		require.Equal(t, connector, returned["service-connector"])
	})

	t.Run("merge-extend-and-override", func(t *testing.T) {
		yamlText := `
default: &default
  timeout: 5s
  retries: 3
  endpoint: example.com

service-connector:
  name: connector
  config:
    <<: *default
    build-flags: production, no-db
    timeout: 42s
`
		connector := map[string]any{
			"name": "connector",
			"config": map[string]any{
				"timeout":     "42s",
				"retries":     3,
				"endpoint":    "example.com",
				"build-flags": "production, no-db",
			},
		}

		returned := AssertValidYamlMapping(t, yamlText)

		require.Equal(t, connector, returned["service-connector"])
	})
}

// Test_multiMerging_notPossible checks that merge operator `<<: *REF` can be used
// just once within a single mapping level.
func Test_multiMerging_notPossible(t *testing.T) {
	yamlText := `
component A: &A
  name:     a
  type:     simple

component B: &B
  date:     2086-04-26
  repeated: false

composite A+B:
  <<: *A
  <<: *B
`

	_, err := ParseYaml(yamlText)

	require.ErrorContains(t, err, `yaml: unmarshal errors:`)
	require.ErrorContains(t, err, `mapping key "<<" already defined at line `)
}
