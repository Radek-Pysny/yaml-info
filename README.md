# YAML (Yet Another Markup Language)

This repo represents my notes on YAML format (mostly inspired by many online resources)
as one can see in this file.

Other than that, there is a collection of unit tests, that I used to learn and verify
my understanding of YAML format together with checking parsing of YAML in Golang.


## Syntax

YAML accepts the entire Unicode character set, except for some control characters. It may be encoded in UTF-8, UTF-16 or UTF-32.

- Whitespace indentation (no tabs) is used to denoting structure.
- Line comments begin with `#` outside of a string.
- List members are denoted by a leading hyphen `-` with one element per line.
    - Or it can be specified by enclosing comma separated elements in square brackets `[ ]`.
- Colon-space `: ` is used as separator between key and value of associative array with one entry per line.
    - Or it can be enclosed in curly braces `{ }` with comma-separated pairs.
    - Question mark in front of key can be used to allow key to contain leading dashes, square brackets, etc. without quotes.
- By `&` character can be denoted an anchor.
- Strings:
    - are ordinarily unquoted;
    - quoted by `"` may represent C-style escape sequences;
    - quoted by `'` supports only `''` escape sequence.
- Block scalars are delimited with indentation with optional modifiers processing line endings:
    - No modifier or `>-`: Folds each paragraph on single line + chomps last line ending.
    - `>` or `>+`: Folds each paragraph on single line + keeps last line ending.
    - `|-`: Keeps all line endings except the last one (is being chomped).
    - `|` or `|+`: Keeps all line endings.
- Multiple documents within a single stream are separated by three hyphens `---`.
    - Three periods `...` optionally end a document within a stream.
- Repeated nodes are initially denoted by an ampersand `&` and thereafter referenced with an asterisk `*`.
- Nodes may be labeled with a type or tag using a double exclamation mark `!!` followed by a string, which can be expanded into a URI.
- To combine two or more associative arrays by merging (overriding attributes) can be used `<<`.


## Data types

| data type | YAML representation |
|---|---|
| Null | `~`, `null`, `NULL`, `Null`, or leave it blank |
| Boolean | `true`, `TRUE`, `True`, `false`, `FALSE`, `False` (before YAML v1.2: `yes`, `no`, `on`, `off`) |
| Integer | `10`, `0b10`, `0x10`, `0o10` (before YAML v1.2: `010`) |
| Floating-point | `3.14`, `12.5e-9`, `.inf`, `.nan` |
| String | `lorem ipsum`, `'Don''t worry'`, `"Be\thappy"` |
| Date and Time | `2022-01-16, 23:59`, `2022-01-16 23:59:59` |

Using `YAML tags` can be used to cast values to a specific data type.

```
text: !!str 2022-01-16

numbers: !!set
  ? 5
  ? 8
  ? 13

image: !!binary
  R0lGODdhCAAIAPAAAAIGAfr4+SwAA
  AAACAAIAAACDIyPeWCsClxDMsZ3CgA7

pair: !!python/tuple
  - black
  - white

center_at: !!python/complex 3.14+2.72j

person: !!python/object:package_name.module_name.ClassName
  age: 42
  first_name: John
  last_name: Doe
```


## Examples

Lists:

```
--- # Block Sequence
- Casablanca
- North by Northwest
- The Man Who Wasn't There
--- # Inline Sequence
[Casablanca, North by Northwest, The Man Who Wasn't There]
```

Associative arrays:

```
--- # Indented Block
  name: John Smith
  age: 33
--- # Inline Block
{name: John Smith, age: 33}
```

String with preserved lines:

```
literal: | # Leading comment of literal block scalar
   There once was a tall man from Ealing
   Who got on a bus to Darjeeling
       It said on the door
       "Please don't sit on the floor"
   So he carefully sat on the ceiling
  # Trailing comment (with smaller indent) of literal block scalar 
```

and with folded lines:

```
data: > # Leading comment of folded block scalar
   Wrapped text
   will be folded
   into a single
   paragraph

   Blank lines denote
   paragraph breaks
      Each extra indented line has also leading newline added
  # Trailing comment (with smaller indent) of folded block scalar 
```

For both block scalars are working chomp operator:

```
no trailing newline: >-
  a
  b
  
  
one trailing newline: >
  a
  b
  
  
three trailing newlines: >+
  a
  b
  
  
```

Objects and lists are important components in yaml and can be mixed.

```
--- # The Smiths
- {name: John Smith, age: 33}
- name: Mary Smith
  age: 27
- [name, age]: [Rae Smith, 4]   # sequences as keys are supported
--- # People, by gender
men: [John Smith, Bill Jones]
women:
  - Mary Smith
  - Susan Williams
```

Easy repeating with anchor and references and merge operator `<<:`:

```
--- # Sequencer protocols for Laser eye surgery
- step:  &id001                  # defines anchor label &id001
    instrument:      KT 1885L
    title:           Function test
    type:            F1
    Unom:            230
    duration:        10s
- step: &id002
    instrument:      KT 1885L
    title:           HV test
    type:            H5
    Unom:            1500
    duration:        2s
    matrix:          +1,-4
- step: *id001                   # refers to the first step (with anchor &id001)
- step: *id002                   # refers to the second step
- step: *id002
- step:                          # second step with extra key-value pair in mapping
    <<: *id002
    enabled:         false
- step:                          # second step with overridden matrix value in mapping
    <<: *id002
    matrix:          +2,-4
```

Usage of explicit data type tag:

```
---
a: 123                     # an integer
b: "123"                   # a string, disambiguated by quotes
c: 123.0                   # a float
d: !!float 123             # also a float via explicit data type prefixed by (!!)
e: !!str 123               # a string, disambiguated by explicit type
f: !!str true              # a string via explicit type
g: Yes                     # a boolean True (yaml1.1), string "Yes" (yaml1.2 and higher)
h: Yes we have No bananas  # a string, "Yes" and "No" disambiguated by context.
i: NO                      # a boolean False (yaml1.1), string "NO" (yaml1.2 and higher), aka YAML Norway problem     
```

Example of base64-necoded binary data:

```
---
picture: !!binary |
  R0lGODdhDQAIAIAAAAAAANn
  Z2SwAAAAADQAIAAACF4SDGQ
  ar3xxbJ9p0qa7R0YxwzaFME
  1IAADs=
```

User-defined data types for object serialization:

```
---
myObject: !myClass { name: Joe, age: 15 }
```

