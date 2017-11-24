// This file was ported from https://github.com/Microsoft/vscode/blob/c0bc1ace7ca3ce2d6b1aeb2bde9d1bb0f4b4bae6/src/vs/base/common/jsonFormatter.ts,
// which is licensed as follows:
//
// Copyright (c) Microsoft Corporation. All rights reserved. Licensed under the MIT License.

package jsonx

import (
	"strings"
	"testing"
)

func TestFormat(t *testing.T) {
	defaultFormatOptions := FormatOptions{TabSize: 2, InsertSpaces: true, EOL: "\n"}
	tests := map[string]struct {
		input   string
		options *FormatOptions
		want    string
	}{
		"object - single property": {
			input: `{"x" : 1}`,
			want: `
{
  "x": 1
}`},
		"object - unicode": {
			input: `{"你好" : 1}`,
			want: `
{
  "你好": 1
}`},
		"object - multi-line": {
			input: `{
  "x": "y"
}`,
			want: `
{
  "x": "y"
}`},
		"object - multiple properties": {
			input: `{"x" : 1,  "y" : "foo", "z"  : true}`,
			want: `{
  "x": 1,
  "y": "foo",
  "z": true
}`},
		"object - no properties ": {
			input: `{"x" : {    },  "y" : {}}`,
			want: `{
  "x": {},
  "y": {}
}`},
		"object - nesting": {
			input: `{"x" : {  "y" : { "z"  : { }}, "a": true}}`,
			want: `{
  "x": {
    "y": {
      "z": {}
    },
    "a": true
  }
}`},
		"array - single items": {
			input: `["[]"]`,
			want: `[
  "[]"
]`},
		"array - multiple items": {
			input: `[true,null,1.2]`,
			want: `[
  true,
  null,
  1.2
]`},
		"array - no items": {
			input: `[      ]`,
			want:  `[]`},
		"array - nesting": {
			input: `[ [], [ [ {} ], "a" ]  ]`,
			want: `[
  [],
  [
    [
      {}
    ],
    "a"
  ]
]`},
		"syntax errors": {
			input: `[ null 1.2 ]`,
			want: `[
  null 1.2
]`},
		"empty lines": {
			input: `{
"a": true,

"b": true
}`,
			options: &FormatOptions{TabSize: 2, InsertSpaces: false, EOL: "\n"},
			want: `{
	"a": true,
	"b": true
}`},
		"single line comment": {
			input: `[ 
//comment 你好
"foo", "bar"
] `,
			want: `[
  //comment 你好
  "foo",
  "bar"
]`},
		"block line comment": {
			input: `[{
        /*comment 你好*/     
"foo" : true
}] `,
			want: `[
  {
    /*comment 你好*/
    "foo": true
  }
]`},
		"single line comment on same line": {
			input: ` {  
        "a": {}// comment 你好
 } `,
			want: `{
  "a": {} // comment 你好
}`},
		"single line comment on same line 2": {
			input: `{ //comment 你好
}`,
			want: `{ //comment 你好
}`},
		"block comment on same line": {
			input: `{      "a": {}, /*comment 你好*/    
        /*comment 你好*/ "b": {},    
		"c": {/*comment 你好*/}    } `,
			want: `{
  "a": {}, /*comment 你好*/
  /*comment 你好*/ "b": {},
  "c": { /*comment 你好*/}
}`},

		"block comment on same line advanced": {
			input: ` {       "d": [
             null
        ] /*comment 你好*/
		,"e": /*comment 你好*/ [null] }`,
			want: `{
  "d": [
    null
  ] /*comment 你好*/,
  "e": /*comment 你好*/ [
    null
  ]
}`},
		"multiple block comments on same line": {
			input: `{      "a": {} /*comment 你好*/, /*comment 你好*/   
        /*comment 你好*/ "b": {}  /*comment 你好*/  } `,
			want: `{
  "a": {} /*comment 你好*/, /*comment 你好*/
  /*comment 你好*/ "b": {} /*comment 你好*/
}`},
		"multiple mixed comments on same line": {
			input: `[ /*comment 你好*/  /*comment 你好*/   // comment 
]`,
			want: `[ /*comment 你好*/ /*comment 你好*/ // comment 
]`},
		"range": {
			input: `{ "a": {},
|"b": [null, null]|
} `,
			want: `{ "a": {},
"b": [
  null,
  null
]
} `},
		"range with existing indent": {
			input: `{ "a": {},
   |"b": [null],
"c": {}
} |`,
			want: `{ "a": {},
  "b": [
    null
  ],
  "c": {}
}`},
		"range with existing indent - tabs": {
			input: `{ "a": {},
|  "b": [null],   
"c": {}
} |    `,
			options: &FormatOptions{TabSize: 2, InsertSpaces: false, EOL: "\n"},
			want: `{ "a": {},
	"b": [
		null
	],
	"c": {}
}`},
		"block comment none-line breaking symbols": {
			input: `{ "a": [ 1
/* comment 你好 */
, 2
/* comment 你好 */
]
/* comment 你好 */
,
 "b": true
/* comment 你好 */
}`,
			want: `{
  "a": [
    1
    /* comment 你好 */
    ,
    2
    /* comment 你好 */
  ]
  /* comment 你好 */
  ,
  "b": true
  /* comment 你好 */
}`},
		"line comment after none-line breaking symbols": {
			input: `{ "a":
// comment 你好
null,
 "b"
// comment 你好
: null
// comment 你好
}`,
			want: `{
  "a":
  // comment 你好
  null,
  "b"
  // comment 你好
  : null
  // comment 你好
}`},
	}
	for label, test := range tests {
		t.Run(label, func(t *testing.T) {
			options := test.options
			if options == nil {
				options = &defaultFormatOptions
			}

			var input string
			var offset, length int
			if strings.Count(test.input, "|") == 2 {
				rangeStart := strings.Index(test.input, "|")
				rangeEnd := strings.LastIndex(test.input, "|")
				input = test.input[:rangeStart] + test.input[rangeStart+1:rangeEnd] + test.input[rangeEnd+1:]
				offset = rangeStart
				length = rangeEnd - rangeStart
			} else {
				input = test.input
				length = len([]rune(test.input))
			}

			edits := FormatRange(input, offset, length, *options)
			output, err := ApplyEdits(input, edits...)
			if err != nil {
				t.Fatal(err)
			}
			if want := strings.TrimPrefix(test.want, "\n"); output != want {
				t.Errorf("formatted\ngot  %s\nwant %s", output, want)
			}
		})
	}
}
