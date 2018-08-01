// This file was ported from https://github.com/Microsoft/vscode/blob/c0bc1ace7ca3ce2d6b1aeb2bde9d1bb0f4b4bae6/src/vs/base/common/jsonEdit.ts,
// which is licensed as follows:
//
// Copyright (c) Microsoft Corporation. All rights reserved. Licensed under the MIT License.

package jsonx

import (
	"encoding/json"
	"testing"
)

func TestComputePropertyEdit(t *testing.T) {
	defaultFormatOptions := FormatOptions{
		TabSize:      2,
		InsertSpaces: true,
		EOL:          "\n",
	}

	type testCase struct {
		input       string
		path        Path
		value       interface{}
		remove      bool
		insertIndex func(properties []string) int
		options     *FormatOptions
		want        string
	}
	assertEdits := func(t *testing.T, tests []testCase) {
		t.Helper()
		for _, test := range tests {
			options := test.options
			if options == nil {
				options = &defaultFormatOptions
			}
			var edits []Edit
			var err error
			if test.remove {
				edits, _, err = ComputePropertyRemoval(test.input, test.path, *options)
			} else {
				edits, _, err = ComputePropertyEdit(test.input, test.path, test.value, test.insertIndex, *options)
			}
			if err != nil {
				t.Errorf("%q: ComputePropertyEdit: %s", test.input, err)
				continue
			}
			output, err := ApplyEdits(test.input, edits...)
			if err != nil {
				t.Errorf("%q: ApplyEdits: %s", test.input, err)
				continue
			}
			if output != test.want {
				t.Errorf("%q: output\ngot  %s\nwant %s", test.input, output, test.want)
			}
		}
	}

	t.Run("set property", func(t *testing.T) {
		assertEdits(t, []testCase{
			{
				input: "{\n  \"x\": \"y\"\n}",
				path:  PropertyPath("x"),
				value: "bar",
				want:  "{\n  \"x\": \"bar\"\n}",
			},
			{
				input: "true",
				path:  nil,
				value: "bar",
				want:  `"bar"`,
			},
			{
				input: "{\n  \"x\": \"y\"\n}",
				path:  PropertyPath("x"),
				value: map[string]bool{"key": true},
				want:  "{\n  \"x\": {\n    \"key\": true\n  }\n}",
			},
			{
				input: "{\n  \"a\": \"b\",  \"x\": \"y\"\n}",
				path:  PropertyPath("a"),
				value: nil,
				want:  "{\n  \"a\": null,  \"x\": \"y\"\n}",
			},
		})
	})
	t.Run("insert property", func(t *testing.T) {
		assertEdits(t, []testCase{
			{
				input: "{}",
				path:  PropertyPath("foo"),
				value: "bar",
				want:  "{\n  \"foo\": \"bar\"\n}",
			},
			{
				input: "{}",
				path:  PropertyPath("foo", "foo2"),
				value: "bar",
				want:  "{\n  \"foo\": {\n    \"foo2\": \"bar\"\n  }\n}",
			},
			{
				input: "{\n}",
				path:  PropertyPath("foo"),
				value: "bar",
				want:  "{\n  \"foo\": \"bar\"\n}",
			},
			{
				input: "  {\n  }",
				path:  PropertyPath("foo"),
				value: "bar",
				want:  "  {\n    \"foo\": \"bar\"\n  }",
			},
			{
				input: "{\n  \"x\": \"y\"\n}",
				path:  PropertyPath("foo"),
				value: "bar",
				want:  "{\n  \"x\": \"y\",\n  \"foo\": \"bar\"\n}",
			},
			{
				input: "{\n  \"x\": \"y\"\n}",
				path:  PropertyPath("e"),
				value: "null",
				want:  "{\n  \"x\": \"y\",\n  \"e\": \"null\"\n}",
			},
			{
				input: "{\n  \"x\": \"y\"\n}",
				path:  PropertyPath("x"),
				value: "bar",
				want:  "{\n  \"x\": \"bar\"\n}",
			},
			{
				input: "{\n  \"x\": {\n    \"a\": 1,\n    \"b\": true\n  }\n}\n",
				path:  PropertyPath("x"),
				value: "bar",
				want:  "{\n  \"x\": \"bar\"\n}\n",
			},
			{
				input: "{\n  \"x\": {\n    \"a\": 1,\n    \"b\": true\n  }\n}\n",
				path:  PropertyPath("x", "b"),
				value: "bar",
				want:  "{\n  \"x\": {\n    \"a\": 1,\n    \"b\": \"bar\"\n  }\n}\n",
			},
			{
				input:       "{\n  \"x\": {\n    \"a\": 1,\n    \"b\": true\n  }\n}\n",
				path:        PropertyPath("x", "c"),
				value:       "bar",
				insertIndex: func([]string) int { return 0 },
				want:        "{\n  \"x\": {\n    \"c\": \"bar\",\n    \"a\": 1,\n    \"b\": true\n  }\n}\n",
			},
			{
				input:       "{\n  \"x\": {\n    \"a\": 1,\n    \"b\": true\n  }\n}\n",
				path:        PropertyPath("x", "c"),
				value:       "bar",
				insertIndex: func([]string) int { return 1 },
				want:        "{\n  \"x\": {\n    \"a\": 1,\n    \"c\": \"bar\",\n    \"b\": true\n  }\n}\n",
			},
			{
				input:       "{\n  \"x\": {\n    \"a\": 1,\n    \"b\": true\n  }\n}\n",
				path:        PropertyPath("x", "c"),
				value:       "bar",
				insertIndex: func([]string) int { return 2 },
				want:        "{\n  \"x\": {\n    \"a\": 1,\n    \"b\": true,\n    \"c\": \"bar\"\n  }\n}\n",
			},
			{
				input: "{\n  \"x\": {\n    \"a\": 1,\n    \"b\": true\n  }\n}\n",
				path:  PropertyPath("c"),
				value: "bar",
				want:  "{\n  \"x\": {\n    \"a\": 1,\n    \"b\": true\n  },\n  \"c\": \"bar\"\n}\n",
			},
			{
				input: "{\n  \"a\": [\n    {\n    } \n  ]  \n}",
				path:  PropertyPath("foo"),
				value: "bar",
				want:  "{\n  \"a\": [\n    {\n    } \n  ],\n  \"foo\": \"bar\"\n}",
			},
			{
				input: "",
				path:  MakePath("foo", 0),
				value: "bar",
				want:  "{\n  \"foo\": [\n    \"bar\"\n  ]\n}",
			},
			{
				input: "//comment",
				path:  MakePath("foo", 0),
				value: "bar",
				want:  "{\n  \"foo\": [\n    \"bar\"\n  ]\n} //comment\n",
			},
			{
				input: "{\n  \"你\": [\n    \"好\"\n  ]  \n}",
				path:  PropertyPath("foo"),
				value: "bar",
				want:  "{\n  \"你\": [\n    \"好\"\n  ],\n  \"foo\": \"bar\"\n}",
			},
		})
	})
	t.Run("remove property", func(t *testing.T) {
		assertEdits(t, []testCase{
			{
				input:  "{\n  \"x\": \"y\"\n}",
				path:   PropertyPath("x"),
				remove: true,
				want:   "{}",
			},
			{
				input:  "{\n  \"x\": \"y\", \"a\": []\n}",
				path:   PropertyPath("x"),
				remove: true,
				want:   "{\n  \"a\": []\n}",
			},
			{
				input:  "{\n  \"x\": \"y\", \"a\": []\n}",
				path:   PropertyPath("a"),
				remove: true,
				want:   "{\n  \"x\": \"y\"\n}",
			},
		})
	})
	t.Run("remove property if ends with comma", func(t *testing.T) {
		assertEdits(t, []testCase{
			{
				input:  "{\n  \"x\": \"y\",\n}",
				path:   PropertyPath("x"),
				remove: true,
				want:   "{}",
			},
			{
				input:  "{\n  \"x\": \"y\" ,\n}",
				path:   PropertyPath("x"),
				remove: true,
				want:   "{}",
			},
			{
				input:  "{\n  \"x\": \"y\", \"a\": [],\n}",
				path:   PropertyPath("a"),
				remove: true,
				want:   "{\n  \"x\": \"y\",\n}",
			},
		})
	})
	t.Run("insert item to empty array", func(t *testing.T) {
		assertEdits(t, []testCase{
			{
				input: "[\n]",
				path:  MakePath(-1),
				value: "bar",
				want:  "[\n  \"bar\"\n]",
			},
		})
	})
	t.Run("insert item to undefined array", func(t *testing.T) {
		assertEdits(t, []testCase{
			{
				input: "{\n}",
				path:  MakePath("foo", -1),
				value: "bar",
				want:  "{\n  \"foo\": [\n    \"bar\"\n  ]\n}",
			},
		})
	})
	t.Run("insert item", func(t *testing.T) {
		assertEdits(t, []testCase{
			{
				input: "[\n  1,\n  2\n]",
				path:  MakePath(-1),
				value: "bar",
				want:  "[\n  1,\n  2,\n  \"bar\"\n]",
			},
		})
	})
	t.Run("remove item in array with one item", func(t *testing.T) {
		assertEdits(t, []testCase{
			{
				input:  "[\n  1\n]",
				path:   MakePath(0),
				remove: true,
				want:   "[]",
			},
		})
	})
	t.Run("remove item in the middle of the array", func(t *testing.T) {
		assertEdits(t, []testCase{
			{
				input:  "[\n  1,\n  2,\n  3\n]",
				path:   MakePath(1),
				remove: true,
				want:   "[\n  1,\n  3\n]",
			},
		})
	})
	t.Run("remove last item in the array", func(t *testing.T) {
		assertEdits(t, []testCase{
			{
				input:  "[\n  1,\n  2,\n  \"bar\"\n]",
				path:   MakePath(2),
				remove: true,
				want:   "[\n  1,\n  2\n]",
			},
		})
	})
	t.Run("remove last item in the array if ends with comma", func(t *testing.T) {
		assertEdits(t, []testCase{
			{
				input:  "[\n  1,\n  \"foo\",\n  \"bar\",\n]",
				path:   MakePath(2),
				remove: true,
				want:   "[\n  1,\n  \"foo\"\n]",
			},
		})
	})
	t.Run("remove last item in the array if there is a comment in the beginning", func(t *testing.T) {
		assertEdits(t, []testCase{
			{
				input:  "// This is a comment\n[\n  1,\n  \"foo\",\n  \"bar\"\n]",
				path:   MakePath(2),
				remove: true,
				want:   "// This is a comment\n[\n  1,\n  \"foo\"\n]",
			},
		})
	})
	t.Run("edit item in array with one item", func(t *testing.T) {
		assertEdits(t, []testCase{
			{
				input: "[\n  1\n]",
				path:  MakePath(0),
				value: 2,
				want:  "[\n  2\n]",
			},
		})
	})
	t.Run("edit item in the middle of the array", func(t *testing.T) {
		assertEdits(t, []testCase{
			{
				input: "[\n  1,\n  2,\n  3\n]",
				path:  MakePath(1),
				value: 4,
				want:  "[\n  1,\n  4,\n  3\n]",
			},
		})
	})
	t.Run("edit last item in the array", func(t *testing.T) {
		assertEdits(t, []testCase{
			{
				input: "[\n  1,\n  2,\n  \"foo\"\n]",
				path:  MakePath(2),
				value: "bar",
				want:  "[\n  1,\n  2,\n  \"bar\"\n]",
			},
		})
	})
	t.Run("edit last item in the array if ends with comma", func(t *testing.T) {
		assertEdits(t, []testCase{
			{
				input: "[\n  1,\n  \"foo\",\n  \"bar\",\n]",
				path:  MakePath(2),
				value: "qux",
				want:  "[\n  1,\n  \"foo\",\n  \"qux\",\n]",
			},
		})
	})
	t.Run("edit last item in the array if there is a comment in the beginning", func(t *testing.T) {
		assertEdits(t, []testCase{
			{
				input: "// This is a comment\n[\n  1,\n  \"foo\",\n  \"bar\"\n]",
				path:  MakePath(2),
				value: "qux",
				want:  "// This is a comment\n[\n  1,\n  \"foo\",\n  \"qux\"\n]",
			},
		})
	})
	t.Run("edit item in a nested array", func(t *testing.T) {
		assertEdits(t, []testCase{
			{
				input: "[\n  1,\n  {\n    \"foo\": [\n      2 // This is a comment\n    ]\n  },\n  3\n]",
				path:  MakePath(1, "foo", 0),
				value: 4,
				want:  "[\n  1,\n  {\n    \"foo\": [\n      4 // This is a comment\n    ]\n  },\n  3\n]",
			},
		})
	})
	t.Run("set raw JSON", func(t *testing.T) {
		assertEdits(t, []testCase{
			{
				input: "{\n  \"x\": \"y\"\n}",
				path:  PropertyPath("x"),
				value: json.RawMessage(`/*c*/"z"`),
				want:  "{\n  \"x\": /*c*/ \"z\"\n}",
			},
		})
	})
}
