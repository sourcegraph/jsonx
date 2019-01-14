// This file was ported from https://github.com/Microsoft/vscode/blob/c0bc1ace7ca3ce2d6b1aeb2bde9d1bb0f4b4bae6/src/vs/base/common/json.ts,
// which is licensed as follows:
//
// Copyright (c) Microsoft Corporation. All rights reserved. Licensed under the MIT License.

package jsonx

import (
	"reflect"
	"testing"
)

func TestScanner(t *testing.T) {
	tests := map[string][]SyntaxKind{
		"{": {OpenBraceToken},
		"}": {CloseBraceToken},
		"[": {OpenBracketToken},
		"]": {CloseBracketToken},
		":": {ColonToken},
		",": {CommaToken},

		// comments
		"// this is a comment 你好":       {LineCommentTrivia},
		"// this is a comment 你好\n":     {LineCommentTrivia, LineBreakTrivia},
		"/* this is a comment 你好*/":     {BlockCommentTrivia},
		"/* this is a \r\ncomment 你好*/": {BlockCommentTrivia},
		"/* this is a \ncomment 你好*/":   {BlockCommentTrivia},
		"/* this is a":                  {BlockCommentTrivia}, // unexpected end,
		"/* this is a \ncomment 你好":     {BlockCommentTrivia},
		"/ ttt":                         {Unknown, Trivia, Unknown}, // broken comment,

		// strings
		`"test"`:              {StringLiteral},
		`"\""`:                {StringLiteral},
		`"\/"`:                {StringLiteral},
		`"\b"`:                {StringLiteral},
		`"\f"`:                {StringLiteral},
		`"\n"`:                {StringLiteral},
		`"\r"`:                {StringLiteral},
		`"\t"`:                {StringLiteral},
		`"\v"`:                {StringLiteral},
		`"` + "\u88ff" + `"`:  {StringLiteral},
		`"` + "​\u2028" + `"`: {StringLiteral},
		`"你好"`:                {StringLiteral},

		// unexpected end
		`"test`:              {StringLiteral},
		`"test` + "\n" + `"`: {StringLiteral, LineBreakTrivia, StringLiteral},

		// numbers
		"0":         {NumericLiteral},
		"0.1":       {NumericLiteral},
		"-0.1":      {NumericLiteral},
		"-1":        {NumericLiteral},
		"1":         {NumericLiteral},
		"123456789": {NumericLiteral},
		"10":        {NumericLiteral},
		"90":        {NumericLiteral},
		"90E+123":   {NumericLiteral},
		"90e+123":   {NumericLiteral},
		"90e-123":   {NumericLiteral},
		"90E-123":   {NumericLiteral},
		"90E123":    {NumericLiteral},
		"90e123":    {NumericLiteral},

		// zero handling
		"01":  {NumericLiteral, NumericLiteral},
		"-01": {NumericLiteral, NumericLiteral},

		// unexpected end
		"-":  {Unknown},
		".0": {Unknown},

		// malformed input
		"/": {Unknown},

		// keywords: true, false, null
		"true":  {TrueKeyword},
		"false": {FalseKeyword},
		"null":  {NullKeyword},

		"true false null": {TrueKeyword, Trivia, FalseKeyword, Trivia, NullKeyword},

		// invalid words
		"nulllll": {Unknown},
		"True":    {Unknown},
		"foo-bar": {Unknown},
		"foo bar": {Unknown, Trivia, Unknown},

		// trivia
		" ":              {Trivia},
		"  \t  ":         {Trivia},
		"  \t  \n  \t  ": {Trivia, LineBreakTrivia, Trivia},
		"\r\n":           {LineBreakTrivia},
		"\r":             {LineBreakTrivia},
		"\n":             {LineBreakTrivia},
		"\n\r":           {LineBreakTrivia, LineBreakTrivia},
		"\n   \n":        {LineBreakTrivia, Trivia, LineBreakTrivia},
	}
	for input, want := range tests {
		scanner := NewScanner(input, ScanOptions{Trivia: true})
		var kinds []SyntaxKind
		for {
			kind := scanner.Scan()
			if kind == EOF {
				break
			}
			kinds = append(kinds, kind)
		}
		if !reflect.DeepEqual(kinds, want) {
			t.Errorf("%q: got kinds %s, want %s", input, kinds, want)
		}
	}
}

func TestScannerErrors(t *testing.T) {
	tests := map[string]struct {
		kind      SyntaxKind
		errorCode ScanErrorCode
	}{
		// invalid characters
		`"` + "\t" + `"`:  {StringLiteral, InvalidCharacter},
		`"` + "\t " + `"`: {StringLiteral, InvalidCharacter},
	}
	for input, test := range tests {
		scanner := NewScanner(input, ScanOptions{Trivia: true})
		kind := scanner.Scan()
		if kind != test.kind {
			t.Errorf("%q: got kind %s, want %s", input, kind, test.kind)
		}
		if scanner.Err() != test.errorCode {
			t.Errorf("%q: got scan error code %s, want %s", input, scanner.Err(), test.errorCode)
		}
	}
}
