package yamlparse

// This file contains unit test for every scalar supported by YAML.
// Namely: null, boolean, integer, floating-point, string (with extra
// care on flow scalars and block scalars), and timestamps.

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Test_validScalar_null shows all the possibilities to denote a null in YAML document.
// There is a single YAML with top-level mapping that contains each possible way to
// denote a `null`-containing key.
//
// Null value can be either empty (completely missing) or one of literals: null, Null, NULL,
// !!null (tag-only), or ~ (simply tilde character).
//
// The used YAML library represent `null` key as true `nil` within `any`.
func Test_validScalar_null(t *testing.T) {
	yamlText := `
lowercase: null
titlecase: Null
uppercase: NULL
tilde:     ~
missing:   
tagged:    !!null
`

	returned := AssertValidYamlMapping(t, yamlText)

	for key := range returned {
		t.Run(key, func(t *testing.T) {
			anyValue, exists := returned[key]

			require.True(t, exists)
			require.Nil(t, anyValue)
		})
	}
}

// Test_validScalar_boolean shows all possibilities how to write a boolean literals
// in YAML document. For each possibility is generated a simple YAML with top-level
// mapping with just one key — enabled.
//
// Integers can be in decimal form, octal form (`0o`-prefixed), hexadecimal form
// (`0x`-prefixed), or binary form (`0b`-prefixed).
//
// In the testCases map is included also commented-out deprecated boolean literals
// that were valid before YAML v1.2. Check "YAML Norway problem" for more details.
func Test_validScalar_boolean(t *testing.T) {
	testCases := map[string]bool{
		// Valid variants of keywords are: lowercase, titlecase, and uppercase.
		`true`:  true,
		`True`:  true,
		`TRUE`:  true,
		`false`: false,
		`False`: false,
		`FALSE`: false,
		// tagged variants
		`!!bool 'true'`:  true,
		`!!bool 'True'`:  true,
		`!!bool 'TRUE'`:  true,
		`!!bool 'false'`: false,
		`!!bool 'False'`: false,
		`!!bool 'FALSE'`: false,
		`!!bool "true"`:  true,
		`!!bool "True"`:  true,
		`!!bool "TRUE"`:  true,
		`!!bool "false"`: false,
		`!!bool "False"`: false,
		`!!bool "FALSE"`: false,
		// legacy variants that were valid before YAML v1.2
		//`on`:  true,
		//`off`: false,
		//`yes`: true,
		//`no`:  false,
		//`On`:  true,
		//`Off`: false,
		//`Yes`: true,
		//`No`:  false,
		//`ON`:  true,
		//`OFF`: false,
		//`YES`: true,
		//`NO`:  false,
	}

	for input, expected := range testCases {
		t.Run(input, func(t *testing.T) {
			yamlText := "enabled: " + input

			anyValue := AssertValidYamlMappingKey(t, yamlText, "enabled")
			value, ok := anyValue.(bool)

			require.True(t, ok)
			require.Equal(t, expected, value)
		})
	}
}

// Test_validScalar_integer shows basic possibilities how to represent integer in YAML.
// Each test case has own small YAML text with mapping of just one key.
//
// Integers can be in decimal form, octal form (`0o`-prefixed), hexadecimal form
// (`0x`-prefixed), or binary form (`0b`-prefixed).
func Test_validScalar_integer(t *testing.T) {
	testCases := []struct {
		title    string
		yamlText string
		key      string
		expected int
	}{
		{
			title:    `positive-decimal`,
			yamlText: `dec: 42`,
			key:      `dec`,
			expected: 42,
		},
		{
			title:    `explicit-positive-decimal`,
			yamlText: `dec: +42`,
			key:      `dec`,
			expected: 42,
		},
		{
			title:    `negative-decimal`,
			yamlText: `dec: -11`,
			key:      `dec`,
			expected: -11,
		},
		{
			title:    `positive-binary`,
			yamlText: `bin: 0b101`,
			key:      `bin`,
			expected: 5,
		},
		{
			title:    `explicit-positive-binary`,
			yamlText: `bin: +0b101`,
			key:      `bin`,
			expected: 5,
		},
		{
			title:    `negative-binary`,
			yamlText: `bin: -0b10101`,
			key:      `bin`,
			expected: -21,
		},
		{
			title:    `positive-octal`,
			yamlText: `oct: 0o377`,
			key:      `oct`,
			expected: 255,
		},
		{
			title:    `explicit-positive-octal`,
			yamlText: `oct: +0o377`,
			key:      `oct`,
			expected: 255,
		},
		{
			title:    `negative-octal`,
			yamlText: `oct: -0o77`,
			key:      `oct`,
			expected: -63,
		},
		{
			title:    `positive-hexadecimal-uppercase-deadbeef`,
			yamlText: `hex: 0xDEADBEEF`,
			key:      `hex`,
			expected: 3_735_928_559,
		},
		{
			title:    `explciit-positive-hexadecimal-lowercase-deadbeef`,
			yamlText: `hex: +0xdeadbeef`,
			key:      `hex`,
			expected: 3_735_928_559,
		},
		{
			title:    `positive-hexadecimal-mixedcase-deadbeef`,
			yamlText: `hex: 0xDeadBeef`,
			key:      `hex`,
			expected: 3_735_928_559,
		},
		{
			title:    `negative-hexadecimal-mixedcase`,
			yamlText: `hex: -0x1aF9`,
			key:      `hex`,
			expected: -6_905,
		},
		{
			title:    `tagged-no-quote`,
			yamlText: `tagged: !!int 12`,
			key:      `tagged`,
			expected: 12,
		},
		{
			title:    `tagged-no-quote-negative`,
			yamlText: `tagged: !!int -69`,
			key:      `tagged`,
			expected: -69,
		},
		{
			title:    `tagged-double-quotes`,
			yamlText: `tagged: !!int "12"`,
			key:      `tagged`,
			expected: 12,
		},
		{
			title:    `tagged-single-quotes`,
			yamlText: `tagged: !!int '12'`,
			key:      `tagged`,
			expected: 12,
		},
		{
			title:    `tagged-single-quotes-hexadecimal`,
			yamlText: `tagged hex: !!int '0x12'`,
			key:      `tagged hex`,
			expected: 18,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			anyValue := AssertValidYamlMappingKey(t, tc.yamlText, tc.key)
			value, ok := anyValue.(int)

			require.True(t, ok)
			require.Equal(t, tc.expected, value)
		})
	}
}

// Test_validScalar_floatingPoint shows basic possibilities how to represent
// floating-point values in YAML. Each test case has own small YAML text with
// mapping of just one key.
//
// Floating-point values have special values .inf (.Inf, .INF; mind also its
// -/+ prefixed variants) and .nan (.NaN, .NAN; mind mixed case usage instead
// of usual titlecase). Next to this general form (3.14) and scientific form
// (1e-12) are supported, as well as missing zero before decimal separator
// (.5 ≡ 0.5) or missing zero after decimal separator (1. ≡ 1.0).
func Test_validScalar_floatingPoint(t *testing.T) {
	const eps = 1e-10

	testCases := []struct {
		title    string
		yaml     string
		key      string
		expected float64
	}{
		{
			title:    `positive-infinity-lowercase`,
			yaml:     `positive infinity:    .inf`,
			key:      `positive infinity`,
			expected: math.Inf(1),
		},
		{
			title:    `positive-infinity-titlecase`,
			yaml:     `positive infinity:    .Inf`,
			key:      `positive infinity`,
			expected: math.Inf(1),
		},
		{
			title:    `positive-infinity-uppercase`,
			yaml:     `positive infinity:    .INF`,
			key:      `positive infinity`,
			expected: math.Inf(1),
		},
		{
			title:    `positive-infinity-lowercase`,
			yaml:     `positive infinity:    +.inf`,
			key:      `positive infinity`,
			expected: math.Inf(1),
		},
		{
			title:    `positive-infinity-titlecase`,
			yaml:     `positive infinity:    +.Inf`,
			key:      `positive infinity`,
			expected: math.Inf(1),
		},
		{
			title:    `positive-infinity-uppercase`,
			yaml:     `positive infinity:    +.INF`,
			key:      `positive infinity`,
			expected: math.Inf(1),
		},
		{
			title:    `negative-infinity-lowercase`,
			yaml:     `neg_inf:   -.inf`,
			key:      `neg_inf`,
			expected: math.Inf(-1),
		},
		{
			title:    `negative-infinity-titlecase`,
			yaml:     `neg_inf:   -.Inf`,
			key:      `neg_inf`,
			expected: math.Inf(-1),
		},
		{
			title:    `negative-infinity-uppercase`,
			yaml:     `neg_inf:   -.INF`,
			key:      `neg_inf`,
			expected: math.Inf(-1),
		},
		{
			title:    `not-a-number-lowercase`,
			yaml:     `not-a-number:   .nan`,
			key:      `not-a-number`,
			expected: math.NaN(),
		},
		{
			title:    `not-a-number-mixedcase`,
			yaml:     `not-a-number:   .NaN`,
			key:      `not-a-number`,
			expected: math.NaN(),
		},
		{
			title:    `not-a-number-uppercase`,
			yaml:     `not-a-number:   .NAN`,
			key:      `not-a-number`,
			expected: math.NaN(),
		},
		{
			title:    `human-readable-form`,
			yaml:     `x: 3.14`,
			key:      `x`,
			expected: 3.14,
		},
		{
			title:    `scientific-form`,
			yaml:     `x: 12.5e-3`,
			key:      `x`,
			expected: 0.0125,
		},
		{
			title:    `no-leading-zero`,
			yaml:     `x: .5`,
			key:      `x`,
			expected: 0.5,
		},
		{
			title:    `no-treailing-zero`,
			yaml:     `x: 2.`,
			key:      `x`,
			expected: 2.0,
		},
		{
			title:    `tagged`,
			yaml:     `tagged int: !!float 1`,
			key:      `tagged int`,
			expected: 1.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			anyValue := AssertValidYamlMappingKey(t, tc.yaml, tc.key)
			value, ok := anyValue.(float64)

			require.True(t, ok)

			if math.IsNaN(tc.expected) { // special handling as NaN ≠ NaN (not a single value in general)
				require.True(t, math.IsNaN(value))

				return
			}

			require.InDelta(t, tc.expected, value, eps)
		})
	}
}

// Test_validScalar_string shows basic possibilities of string representations in YAML.
// Usage of string without quotes, with single- and double-quotes, and vary basics of
// literal/folring block scalars.
//
// For more examples on block scalars (both literal and folding) see Test_validScalar_blockScalars.
func Test_validScalar_string(t *testing.T) {
	testCases := []struct {
		title    string
		yaml     string
		key      string
		expected string
	}{
		{
			title:    `empty-within-double-quotes`,
			yaml:     `empty-string: ""`,
			key:      `empty-string`,
			expected: "",
		},
		{
			title:    `empty-within-single-quotes`,
			yaml:     `empty-string: ''`,
			key:      `empty-string`,
			expected: "",
		},
		{
			title:    `no-quotes-needed`,
			yaml:     `text: simple text`,
			key:      `text`,
			expected: "simple text",
		},
		{
			title:    `no-quotes-needed-and-leading-whitespace-ignored`,
			yaml:     `text:    simple text`,
			key:      `text`,
			expected: "simple text",
		},
		{
			title:    `no-quotes-needed-and-trailing-whitespace-ignored`,
			yaml:     `text: simple text     `,
			key:      `text`,
			expected: "simple text",
		},
		{
			title:    `single-quote-escape`,
			yaml:     `what: 'Don''t worry...'`,
			key:      `what`,
			expected: "Don't worry...",
		},
		{
			title:    `double-quote-escape`,
			yaml:     `what: "She said: \"Hey!\"\tand he run..."`,
			key:      `what`,
			expected: `She said: "Hey!"` + "\t" + `and he run...`,
		},
		{
			title: `folding-block-scalar-with-last-line-ending`,
			yaml: `folded-block: >` + "\n" +
				`  that long line` + "\n" +
				`  is still on one line` + "\n",
			key:      `folded-block`,
			expected: `that long line is still on one line` + "\n",
		},
		{
			title: `folding-block-scalar-without-last-line-ending`,
			yaml: `folded-block: >-` + "\n" +
				`  that long line` + "\n" +
				`  is still on one line` + "\n",
			key:      `folded-block`,
			expected: `that long line is still on one line`,
		},
		{
			title: `block-scalar-with-last-line-ending`,
			yaml: `block: |` + "\n" +
				`  1st line` + "\n" +
				`  2nd line` + "\n" +
				`  3rd line` + "\n",
			key:      `block`,
			expected: `1st line` + "\n" + `2nd line` + "\n" + `3rd line` + "\n",
		},
		{
			title: `block-scalar-without-last-line-ending`,
			yaml: `block: |-` + "\n" +
				`  1st line` + "\n" +
				`  2nd line` + "\n" +
				`  3rd line` + "\n",
			key:      `block`,
			expected: `1st line` + "\n" + `2nd line` + "\n" + `3rd line`,
		},
		{
			title: `block-scalar-removing-just-1-space-from-indentation-of-2`,
			yaml: `block: |1` + "\n" +
				`  1st line` + "\n" +
				`  2nd line` + "\n" +
				`  3rd line` + "\n",
			key:      `block`,
			expected: ` 1st line` + "\n" + ` 2nd line` + "\n" + ` 3rd line` + "\n",
		},
		{
			title: `block-scalar-removing-just-1-space-from-indentation-of-5`,
			yaml: `block: |1` + "\n" +
				`     1st line` + "\n" +
				`     2nd line` + "\n" +
				`     3rd line` + "\n",
			key:      `block`,
			expected: `    1st line` + "\n" + `    2nd line` + "\n" + `    3rd line` + "\n",
		},
		{
			title:    `parsing-true-as-strign`,
			yaml:     `boolean-as-string: "true"`,
			key:      `boolean-as-string`,
			expected: "true",
		},
		{
			title:    `parsing-false-as-strign`,
			yaml:     `bool-as-str: "false"`,
			key:      `bool-as-str`,
			expected: "false",
		},
		{
			title:    `parsing-integer-as-strign`,
			yaml:     `int-as-str: "42"`,
			key:      `int-as-str`,
			expected: "42",
		},
		{
			title:    `parsing-floating-point-as-strign`,
			yaml:     `fp-as-str: "3.14"`,
			key:      `fp-as-str`,
			expected: "3.14",
		},
		{
			title:    `tagged-floating-point`,
			yaml:     `float tagged as str: !!str 3.14`,
			key:      `float tagged as str`,
			expected: "3.14",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			anyValue := AssertValidYamlMappingKey(t, tc.yaml, tc.key)
			value, ok := anyValue.(string)

			require.True(t, ok)
			require.Equal(t, tc.expected, value)
		})
	}
}

// Test_validScalar_timestamp shows possibilities of timestamps in YAML.
// Timestamps are all ISO 9601 based. There is supported
//   - date-only (YYYY-MM-DD);
//   - with less precise base (YYYY-MM-DDThh:mm:ss) or
//     more precise base (YYYY-MM-DDThh:mm:ss.z, where z is 1-9 digits);
//   - with suffix Z suffix (meaning Zulu time, aka UTC timezone), or
//     +hh:mm, or -hh:mm timezone offset suffix.
func Test_validScalar_timestamp(t *testing.T) {
	var zonePlusHour = time.FixedZone("", 3600)

	testCases := []struct {
		title    string
		yaml     string
		expected time.Time
	}{
		{
			title:    `date-only`,
			yaml:     `ts: 2026-10-07`,
			expected: time.Date(2026, 10, 07, 0, 0, 0, 0, time.UTC),
		},
		{
			title:    `basic-instant-in-zulu`,
			yaml:     `ts: 2026-10-07T16:48:01Z`,
			expected: time.Date(2026, 10, 07, 16, 48, 1, 0, time.UTC),
		},
		{
			title:    `basic-instant-with-custom-timezone`,
			yaml:     `ts: 2026-10-07T16:48:01+01:00`,
			expected: time.Date(2026, 10, 07, 16, 48, 1, 0, zonePlusHour),
		},
		{
			title:    `precise-instant-in-zulu`,
			yaml:     `ts: 2026-10-07T16:48:01.123456789Z`,
			expected: time.Date(2026, 10, 07, 16, 48, 1, 123456789, time.UTC),
		},
		{
			title:    `precise-instant-with-custom-timezone`,
			yaml:     `ts: 2026-10-07T16:48:01.123456789+01:00`,
			expected: time.Date(2026, 10, 07, 16, 48, 1, 123456789, zonePlusHour),
		},
		{
			title:    `tagged-date-only`,
			yaml:     `ts: !!timestamp 2026-10-07`,
			expected: time.Date(2026, 10, 07, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			anyValue := AssertValidYamlMappingKey(t, tc.yaml, "ts")
			value, ok := anyValue.(time.Time)

			require.True(t, ok)
			require.Equal(t, tc.expected, value)
		})
	}
}

// Test_validScalar_flowScalars shows more complex usages of so-called flow scalars
// (scalars that are not block scalars nor complex structures). It is mainly about
// how is handled whitespace/indentation, where can appear special characters in
// normal text corpus (without any special meaning or constructing invalid YAML)
// and special character combinations (e.g. colon followed by space is trickier).
//
// There might be covered almost all examples from:
// https://www.yaml.info/learn/quote.html#flow.
func Test_validScalar_flowScalars(t *testing.T) {
	testCases := []struct {
		title    string
		yaml     string
		expected string
	}{
		{
			title:    `backslash-without-escape`,
			yaml:     `text: a string with \ backslash does not need escape`,
			expected: `a string with \ backslash does not need escape`,
		},
		{
			title:    `other-special-characters-without-escape`,
			yaml:     `text: when not first a lot can be included @ $ % ^ & * ( ) [ ] { } `,
			expected: `when not first a lot can be included @ $ % ^ & * ( ) [ ] { }`,
		},
		{
			title: `single-line-text-on-multiple-lines`,
			yaml: `text: this is
  all still
  just one line`,
			expected: `this is all still just one line`,
		},
		{
			title: `single-line-text-on-multiple-lines-nothing-after-key`,
			yaml: `text:
  this is
  all still
  just one line`,
			expected: `this is all still just one line`,
		},
		{
			title: `ignoring-leading-and-trailing-shitespaces`,
			yaml: `text: a
  b              
                     c
  d 
     e`,
			expected: `a b c d e`,
		},
		{
			title: `whitespace-inside-is-kept`,
			yaml: `text: a   b
  c   d     e
`,
			expected: `a   b c   d     e`,
		},
		{
			title: `tabs-are-tricky`,
			yaml: `text: a	b
  c	d	e
`,
			expected: "a\tb c\td\te",
		},
		{
			title: `blank-line-is-interpreted-as-newline`,
			yaml: `text: a

  b c

  d
  e
`,
			expected: "a\nb c\nd e",
		},
		{
			title: `even-repeated--blank-lines-are-interpreted-as-newlines`,
			yaml: `text: a


  b c


  d
  e
`,
			expected: "a\n\nb c\n\nd e",
		},
		{
			title: `comment-can-be-on-last-line-only`,
			yaml: `text: a
  b # this is comment and it is skipped by the parser
`,
			expected: "a b",
		},
		{
			title: `caveat-with-dash-followed-by-space`,
			yaml: `text: a
  - b
     - c
`,
			expected: "a - b - c",
		},
		{
			title:    `started-with-ampersand-needs-quote`,
			yaml:     `text: '& this is not an anchor'`,
			expected: `& this is not an anchor`,
		},
		{
			title:    `started-with-ampersand-needs-quote`,
			yaml:     `text: "! this is not a tag"`,
			expected: `! this is not a tag`,
		},
		{
			title:    `started-with-asterisk-needs-quote`,
			yaml:     `text: '* this is not an alias'`,
			expected: `* this is not an alias`,
		},
		{
			title:    `started-with-dash-needs-quote`,
			yaml:     `text: "- this is not a block sequence entry"`,
			expected: `- this is not a block sequence entry`,
		},
		{
			title:    `started-with-colon-needs-quote`,
			yaml:     `text: ': this is not a block mapping entry'`,
			expected: `: this is not a block mapping entry`,
		},
		{
			title:    `started-with-question-mark-needs-quote`,
			yaml:     `text: "? this is not an explicit mapping key"`,
			expected: `? this is not an explicit mapping key`,
		},
		{
			title:    `started-with-bracket-needs-quote`,
			yaml:     `text: '[ this is not a flow sequence'`,
			expected: `[ this is not a flow sequence`,
		},
		{
			title:    `started-with-curly-bracket-needs-quote`,
			yaml:     `text: "{ this is not a flow mapping"`,
			expected: `{ this is not a flow mapping`,
		},
		{
			title:    `started-with-comma-needs-quote`,
			yaml:     `text: ', this is not a flow collection separator'`,
			expected: `, this is not a flow collection separator`,
		},
		{
			title:    `started-with-hash-needs-quote`,
			yaml:     `text: "# this is not a comment"`,
			expected: `# this is not a comment`,
		},
		{
			title:    `started-with-vertical-bar-needs-quote`,
			yaml:     `text: '| this is not a block scalar'`,
			expected: `| this is not a block scalar`,
		},
		{
			title:    `started-with-greater-than-sign-needs-quote`,
			yaml:     `text: "> this is not a block scalar"`,
			expected: `> this is not a block scalar`,
		},
		{
			title:    `started-with-ampersand-needs-quote`,
			yaml:     `text: '@ this is not a reserved character'`,
			expected: `@ this is not a reserved character`,
		},
		{
			title:    `started-with-backtick-needs-quote`,
			yaml:     "text: \"` this is not a reserved character\"",
			expected: "` this is not a reserved character",
		},
		{
			title:    `started-with-percent-needs-quote`,
			yaml:     `text: '% this is not a directive'`,
			expected: `% this is not a directive`,
		},
		{
			title:    `colon-followed-by-space-needs-quote`,
			yaml:     `text: 'any text with : colon followed by space has to be quoted'`,
			expected: `any text with : colon followed by space has to be quoted`,
		},
		{
			title:    `colon-not-followed-by-space-does-not-need-quote`,
			yaml:     `text: This is ok :because of no space after colon`,
			expected: `This is ok :because of no space after colon`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			anyValue := AssertValidYamlMappingKey(t, tc.yaml, "text")
			value, ok := anyValue.(string)

			require.True(t, ok)
			require.Equal(t, tc.expected, value)
		})
	}
}

// Test_validScalar_blockScalars shows perhaps all possibilities of block scalars
// (aka literal block scalar denoted by `|` and folding block scalar denoted by `>`).
// There are covered topics like:
// - which whitespaces are skipped and which are kept;
// - how to include newline in folded block scalar (blank line or extra indent);
// - trailing newlines with chomp indicators (`|-`, `|+`, `>-`, and `>+`);
// - chomping the desired count of leading spaces (e.g., `|1`);
// - only leading/trailing comments on block scalars.
//
// The vast majority of examples from https://www.yaml.info/learn/quote.html#block
// should be covered by this unit test.
func Test_validScalar_blockScalars(t *testing.T) {
	testCases := []struct {
		title    string
		key      string
		yaml     string
		expected string
	}{
		{
			title: `literal-scalar-block`,
			key:   "literal",
			yaml: `literal: |
  line 1
  line 2
`,
			expected: "line 1\nline 2\n",
		},
		{
			title: `literal-scalar-block-extra-indent`,
			key:   "literal",
			yaml: `literal: |
  line 1
    line 2
  line 3
`,
			expected: "line 1\n  line 2\nline 3\n",
		},
		{
			title: `literal-scalar-block-with-bash-script`,
			key:   "bash",
			yaml: `bash: |
  #!/usr/bin/env bash
  echo "Help, I'm trapped in a YAML document!"
  exit 1
`,
			expected: "#!/usr/bin/env bash\necho \"Help, I'm trapped in a YAML document!\"\nexit 1\n",
		},
		{
			title: `literal-scalar-block-with-embedded-yaml`,
			key:   "yaml",
			yaml: `yaml: |
  literal: |
    line 1
    line 2
`,
			expected: "literal: |\n  line 1\n  line 2\n",
		},
		{
			title: `folded-scalar-block-with-long-bash-command`,
			key:   "folded",
			yaml: `folded: >
  apt-get update
  && apt-get install -y
  git tig vim jq tmux tmate git-subrepo
`,
			expected: "apt-get update && apt-get install -y git tig vim jq tmux tmate git-subrepo\n",
		},
		{
			title: `folded-scalar-block-with-two-long-lines`,
			key:   "a text with two long lines",
			yaml: `a text with two long lines: >
  this is the first
  long line

  and this is the
  second
`,
			expected: "this is the first long line\nand this is the second\n",
		},
		{
			title: `folded-scalar-block-with-four-long-lines`,
			key:   "a long text with enforced newlines",
			yaml: `a long text with enforced newlines: >
  line
  one
    line two
    line three
  line
  four
`,
			expected: "line one\n  line two\n  line three\nline four\n",
		},
		{
			title: `folded-scalar-block-with-traling-spaces`,
			key:   "a long text with trailing spaces",
			yaml: `a long text with trailing spaces: >
  trailing spaces   
  continued
`,
			expected: "trailing spaces    continued\n",
		},
		{
			title: `literal-scalar-block-with-no-comments`,
			key:   "literal",
			yaml: `literal: |
  # no comment
  a
  # no comment
  b
  # no comment
`,
			expected: "# no comment\na\n# no comment\nb\n# no comment\n",
		},
		{
			title: `folded-scalar-block-with-no-comments`,
			key:   "folded",
			yaml: `folded: >
  # no comment
  a
  # no comment
  b
  # no comment
`,
			expected: "# no comment a # no comment b # no comment\n",
		},
		{
			title: `literal-scalar-block-with-trailing-comment`,
			key:   "literal",
			yaml: `literal: |
  a
  b
 # comment for scalar block (less indented comment is a comment and it stops block)
`,
			expected: "a\nb\n",
		},
		{
			title: `folded-scalar-block-with-trailing-comment`,
			key:   "folded",
			yaml: `folded: >
  a
  b
 # comment for scalar block (less indented comment is a comment and it stops block)
`,
			expected: "a b\n",
		},
		{
			title: `literal-scalar-block-with-leading-comment`,
			key:   "literal",
			yaml: `literal: | # a block scalar leading comment
  a
  b
`,
			expected: "a\nb\n",
		},
		{
			title: `folded-scalar-block-with-leading-comment`,
			key:   "folded",
			yaml: `folded: > # a block scalar leading comment
  a
  b
`,
			expected: "a b\n",
		},
		{
			title: `literal-scalar-block-no-chomp-indicator`,
			key:   "literal",
			yaml: `literal: |
  a
  b



`,
			expected: "a\nb\n", // no chomp indicator → just one trailing newline
		},
		{
			title: `folded-scalar-block-no-chomp-indicator`,
			key:   "folded",
			yaml: `folded: > # a block scalar leading comment
  a
  b



`,
			expected: "a b\n", // no chomp indicator → just one trailing newline
		},
		{
			title: `literal-scalar-block-negative-chomp-indicator`,
			key:   "literal",
			yaml: `literal: |-
  a
  b



`,
			expected: "a\nb", // - chomp indicator → no trailing newline
		},
		{
			title: `folded-scalar-block-negative-chomp-indicator`,
			key:   "folded",
			yaml: `folded: >-
  a
  b



`,
			expected: "a b", // - chomp indicator → no trailing newline
		},
		{
			title: `literal-scalar-block-positive-chomp-indicator`,
			key:   "literal",
			yaml: `literal: |+
  a
  b



`,
			expected: "a\nb\n\n\n\n", // + chomp indicator → keeping all trailing newline + 1
		},
		{
			title: `folded-scalar-block-positive-chomp-indicator`,
			key:   "folded",
			yaml: `folded: >+
  a
  b



`,
			expected: "a b\n\n\n\n", // + chomp indicator → keeping all trailing newline + 1
		},
		{
			title: `literal-scalar-block-chomping-1-column`,
			key:   "literal",
			yaml: `literal: |1
  a
    b
   c
`,
			expected: " a\n   b\n  c\n", // from indent is chomped just 1 space
		},
		{
			title: `folded-scalar-block-chomping-1-column`,
			key:   "folded",
			yaml: `folded: >1 # it is exactly the same as '|1'
  a
    b
   c
`,
			expected: " a\n   b\n  c\n", // from indent is chomped just 1 space
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			anyValue := AssertValidYamlMappingKey(t, tc.yaml, tc.key)
			value, ok := anyValue.(string)

			require.True(t, ok)
			require.Equal(t, tc.expected, value)
		})
	}
}

// Test_validScalar_blockScalar_withinItemize shows that each bullet of itemize
// can be handled as block scalar, as it might not be visible at the first sight.
func Test_validScalar_blockScalar_withinItemize(t *testing.T) {
	yamlText := `- |
  trailing spaces   
  trailing tab	
`

	expected := []string{"trailing spaces   \ntrailing tab\t\n"}

	returned := AssertValidYamlSequenceOfText(t, yamlText)

	require.Equal(t, len(expected), len(returned))
	require.Equal(t, expected, returned)
}
