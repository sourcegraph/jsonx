// This file was ported from https://github.com/Microsoft/vscode/blob/c0bc1ace7ca3ce2d6b1aeb2bde9d1bb0f4b4bae6/src/vs/base/common/json.ts,
// which is licensed as follows:
//
// Copyright (c) Microsoft Corporation. All rights reserved. Licensed under the MIT License.

package jsonx

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	defaultOptions := ParseOptions{Comments: true, TrailingCommas: true}
	tests := map[string]struct {
		options *ParseOptions // if nil, use defaultOptions
		want    string
		errors  bool
	}{
		"": {want: ""},

		// literals
		"true":                      {want: "true"},
		"false":                     {want: "false"},
		"null":                      {want: "null"},
		`"foo"`:                     {want: `"foo"`},
		`"\"-\\-\/-\b-\f-\n-\r-\t"`: {want: `"\"-\\-/-\u0008-\u000c-\n-\r-\t"`},
		`"\u00DC"`:                  {want: `"Ü"`},
		"9":                         {want: "9"},
		"-9":                        {want: "-9"},
		"0.129":                     {want: "0.129"},
		"23e3":                      {want: "23e3"},
		"1.2E+3":                    {want: "1.2E+3"},
		"1.2E-3":                    {want: "1.2E-3"},
		"1.2E-3 // comment":         {want: "1.2E-3"},

		// objects
		"{}":                                                                                                        {want: "{}"},
		`{ "foo": true }`:                                                                                           {want: `{"foo":true}`},
		`{ "bar": 8, "xoo": "foo" }`:                                                                                {want: `{"bar":8,"xoo":"foo"}`},
		`{ "hello": [], "world": {} }`:                                                                              {want: `{"hello":[],"world":{}}`},
		`{ "a": false, "b": true, "c": [ 7.4 ] }`:                                                                   {want: `{"a":false,"b":true,"c":[7.4]}`},
		`{ "blockComment": ["/*", "*/"], "brackets": [ ["{", "}"], ["[", "]"], ["(", ")"] ], "lineComment": "//" }`: {want: `{"blockComment":["/*","*/"],"brackets":[["{","}"],["[","]"],["(",")"]],"lineComment":"//"}`},
		`{ "hello": { "again": { "inside": 5 }, "world": 1 }}`:                                                      {want: `{"hello":{"again":{"inside":5},"world":1}}`},
		`{ "foo": /*hello*/true }`:                                                                                  {want: `{"foo":true}`},

		// arrays
		"[]":                {want: "[]"},
		"[ [], [ [] ]]":     {want: "[[],[[]]]"},
		"[ 1, 2, 3 ]":       {want: "[1,2,3]"},
		`[ { "a": null } ]`: {want: `[{"a":null}]`},

		// objects with errors
		"{,}":                      {want: "{}", errors: true},
		`{ "foo": true, }`:         {options: &ParseOptions{TrailingCommas: false}, want: `{"foo":true}`, errors: true},
		`{ "bar": 8 "xoo": "foo"}`: {want: `{"bar":8,"xoo":"foo"}`, errors: true},
		`{ ,"bar": 8 }`:            {want: `{"bar":8}`, errors: true},
		`{ "bar": 8, "foo": }`:     {want: `{"bar":8}`, errors: true},
		`{ 8, "foo": 9 }`:          {want: `{"foo":9}`, errors: true},

		// array with errors
		"[,]":           {want: "[]", errors: true},
		"[ 1, 2, ]":     {options: &ParseOptions{TrailingCommas: false}, want: "[1,2]", errors: true},
		"[ 1 2, 3]":     {want: "[1,2,3]", errors: true},
		"[ ,1, 2, 3 ]":  {want: "[1,2,3]", errors: true},
		"[ ,1, 2, 3, ]": {options: &ParseOptions{TrailingCommas: false}, want: "[1,2,3]", errors: true},

		// disallow commments
		`[ 1, 2, null, "foo" ]`:         {options: &ParseOptions{Comments: false}, want: `[1,2,null,"foo"]`},
		`{ "hello1": [], "world": {} }`: {options: &ParseOptions{Comments: false}, want: `{"hello1":[],"world":{}}`},
		`{ "foo": /*comment*/ true }`:   {options: &ParseOptions{Comments: false}, want: `{"foo":true}`, errors: true},

		// trailing comma
		`{ "hello": [], }`:               {want: `{"hello":[]}`},
		`{ "hello": [] }`:                {want: `{"hello":[]}`},
		`{ "hello": [], "world": {}, }`:  {want: `{"hello":[],"world":{}}`},
		`{ "hello2": [], "world": {} }`:  {want: `{"hello2":[],"world":{}}`},
		"[ 1, 5, ]":                      {want: "[1,5]"},
		`{ "hello2": [], }`:              {options: &ParseOptions{TrailingCommas: false}, want: `{"hello2":[]}`, errors: true},
		`{ "hello2": [], "world": {}, }`: {options: &ParseOptions{TrailingCommas: false}, want: `{"hello2":[],"world":{}}`, errors: true},
		"[ 1, 6, ]":                      {options: &ParseOptions{TrailingCommas: false}, want: "[1,6]", errors: true},
	}
	for input, test := range tests {
		label := fmt.Sprintf("%q", input)

		options := test.options
		if options == nil {
			options = &defaultOptions
		} else {
			label += fmt.Sprintf(" with options %+v", options)
		}

		output, errors := Parse(input, *options)
		if test.errors && errors == nil {
			t.Errorf("%s: got no parse errors, want parse errors", label)
		}
		if !test.errors && errors != nil {
			t.Errorf("%s: got parse errors %v, want no parse errors", label, errors)
		}
		if string(output) != test.want {
			t.Errorf("%s: got output %s, want %s", label, output, test.want)
		}
	}
}
