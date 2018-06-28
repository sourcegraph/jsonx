// This file was ported from https://github.com/Microsoft/vscode/blob/c0bc1ace7ca3ce2d6b1aeb2bde9d1bb0f4b4bae6/src/vs/base/common/json.ts,
// which is licensed as follows:
//
// Copyright (c) Microsoft Corporation. All rights reserved. Licensed under the MIT License.

package jsonx

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestParseTree(t *testing.T) {
	tests := map[string]struct {
		want   *Node
		errors []ParseErrorCode
	}{
		// literals
		`true`:      {want: &Node{Type: Boolean, Offset: 0, Length: 4, Value: true}},
		`false`:     {want: &Node{Type: Boolean, Offset: 0, Length: 5, Value: false}},
		`null`:      {want: &Node{Type: Null, Offset: 0, Length: 4, Value: nil}},
		`23`:        {want: &Node{Type: Number, Offset: 0, Length: 2, Value: json.Number("23")}},
		`-1.93e-19`: {want: &Node{Type: Number, Offset: 0, Length: 9, Value: json.Number("-1.93e-19")}},
		`"hello"`:   {want: &Node{Type: String, Offset: 0, Length: 7, Value: "hello"}},

		// arrays
		`[]`: {want: &Node{Type: Array, Offset: 0, Length: 2}},
		`[ 1 ]`: {
			want: &Node{
				Type: Array, Offset: 0, Length: 5, Children: []*Node{
					{Type: Number, Offset: 2, Length: 1, Value: json.Number("1")},
				},
			},
		},
		`[ 1,"x"]`: {
			want: &Node{
				Type: Array, Offset: 0, Length: 8, Children: []*Node{
					{Type: Number, Offset: 2, Length: 1, Value: json.Number("1")},
					{Type: String, Offset: 4, Length: 3, Value: "x"},
				},
			},
		},
		`[[]]`: {
			want: &Node{Type: Array, Offset: 0, Length: 4, Children: []*Node{
				{Type: Array, Offset: 1, Length: 2},
			}},
		},

		// objects
		`{ }`: {want: &Node{Type: Object, Offset: 0, Length: 3}},
		`{ "val": 1 }`: {
			want: &Node{
				Type: Object, Offset: 0, Length: 12, Children: []*Node{
					{
						Type: Property, Offset: 2, Length: 8, ColumnOffset: 7, Children: []*Node{
							{Type: String, Offset: 2, Length: 5, Value: "val"},
							{Type: Number, Offset: 9, Length: 1, Value: json.Number("1")},
						},
					},
				},
			},
		},
		`{"id": "$", "v": [ null, null] }`: {
			want: &Node{
				Type: Object, Offset: 0, Length: 32, Children: []*Node{
					{
						Type: Property, Offset: 1, Length: 9, ColumnOffset: 5, Children: []*Node{
							{Type: String, Offset: 1, Length: 4, Value: "id"},
							{Type: String, Offset: 7, Length: 3, Value: "$"},
						},
					},
					{
						Type: Property, Offset: 12, Length: 18, ColumnOffset: 15, Children: []*Node{
							{Type: String, Offset: 12, Length: 3, Value: "v"},
							{
								Type: Array, Offset: 17, Length: 13, Children: []*Node{
									{Type: Null, Offset: 19, Length: 4, Value: nil},
									{Type: Null, Offset: 25, Length: 4, Value: nil},
								},
							},
						},
					},
				},
			},
		},
		`{  "id": { "foo": { } } , }`: {
			want: &Node{
				Type: Object, Offset: 0, Length: 27, Children: []*Node{
					{
						Type: Property, Offset: 3, Length: 20, ColumnOffset: 7, Children: []*Node{
							{Type: String, Offset: 3, Length: 4, Value: "id"},
							{
								Type: Object, Offset: 9, Length: 14, Children: []*Node{
									{
										Type: Property, Offset: 11, Length: 10, ColumnOffset: 16, Children: []*Node{
											{Type: String, Offset: 11, Length: 5, Value: "foo"},
											{Type: Object, Offset: 18, Length: 3},
										},
									},
								},
							},
						},
					},
				},
			},
			errors: []ParseErrorCode{PropertyNameExpected, ValueExpected},
		},
	}
	for input, test := range tests {
		tree, errors := ParseTree(input, ParseOptions{Comments: false, TrailingCommas: false})
		if !reflect.DeepEqual(errors, test.errors) {
			t.Errorf("%q: got errors %v, want %v", input, errors, test.errors)
		}
		clearParentPointers(tree)
		if !reflect.DeepEqual(tree, test.want) {
			t.Errorf("%q: tree\ngot  %s\nwant %s", input, printAndDestroyTree(tree), printAndDestroyTree(test.want))
		}
	}
}

func clearParentPointers(node *Node) {
	// Clear parent pointers to avoid cycles when (e.g) stringifying.
	if node == nil {
		return
	}
	node.Parent = nil
	for _, child := range node.Children {
		clearParentPointers(child)
	}
}

func printAndDestroyTree(root *Node) string {
	clearParentPointers(root)
	data, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(data)
}

func TestPath_JSON(t *testing.T) {
	p1 := Path{
		{IsProperty: true, Property: "a"},
		{IsProperty: true, Property: ""},
		{Index: 0},
		{Index: 1},
	}

	data, err := json.Marshal(p1)
	if err != nil {
		t.Fatal(err)
	}
	if want := `["a","",0,1]`; string(data) != want {
		t.Errorf("got %s, want %s", data, want)
	}

	var p2 Path
	if err := json.Unmarshal(data, &p2); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(p1, p2) {
		t.Errorf("got %+v, want %+v", p1, p2)
	}

}
