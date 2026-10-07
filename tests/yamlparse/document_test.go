package yamlparse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_documentHeaderAndDocumentFooter ensures the very basic usage of document header
// `---` and document footer `...` in YAML. This unit test still relies on simple utilities
// that can handle just a single document in stream.
//
// There might be later introduced some additional unit test handling multiple documents
// within single YAML stream.
func Test_documentHeaderAndDocumentFooter(t *testing.T) {
	testCases := []struct {
		title    string
		yamlText string
		expected string
	}{
		{
			title: "just-document-header",
			yamlText: `---
document text`,
			expected: "document text",
		},
		{
			title: "document-header-with-comment",
			yamlText: `--- # sample document with just a top-level text
document text`,
			expected: "document text",
		},
		{
			title: "parsing-just-one-document-at-time-1",
			yamlText: `--- # this is parsed
document text
--- # the document is ignored as just one decode s requested
ignored document text`,
			expected: "document text",
		},
		{
			title: "parsing-just-one-document-at-time-2",
			yamlText: `--- # this is parsed
document text
--- # the document is ignored as just one decode s requested
ignored document text
---
other ignored document text
--- yet another ignored document text`,
			expected: "document text",
		},
		{
			title: "just-document-footer",
			yamlText: `document text
...`,
			expected: "document text",
		},
		{
			title: "document-footer-with-comment",
			yamlText: `document text
... # this is document footer`,
			expected: "document text",
		},
		{
			title: "ignore-stuff-after-document-footer",
			yamlText: `document text
... # everything below is ignored
ignored document text`,
			expected: "document text",
		},
		{
			title: "multiple-documents-in-yaml-stream-with-full-syntax",
			yamlText: `
--- # 1st document
processed document text
... end of 1st document

--- # 2nd document
ignored document content
... end of 2nd document

--- 3rd document
another ignored content
... end of 3rd document
`,
			expected: "processed document text",
		},
		{
			title: "scalar-with-no-document-header-inside",
			yamlText: `"valid
  ---
  scalar"
`,
			expected: "valid --- scalar",
		},
		{
			title: "scalar-with-no-document-footer-inside",
			yamlText: `"valid
  ...
  scalar"
`,
			expected: "valid ... scalar",
		},
		{
			title: "scalar-with-not-a-document-footer-inside",
			yamlText: `--- >
  block
  scalar
  ---
  continuation
`,
			expected: "block scalar --- continuation\n",
		},
		{
			title: "scalar-with-not-a-document-footer-inside",
			yamlText: `>
  block
  scalar
---
completely another document
`,
			expected: "block scalar\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			returned := AssertValidYamlString(t, tc.yamlText)

			assert.Equal(t, tc.expected, returned)
		})
	}
}
