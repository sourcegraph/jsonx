package jsonx

import "testing"

func TestWalkerParseErrors(t *testing.T) {
	const noParseErrorCode ParseErrorCode = -1

	t.Run("comments and trailing commas enabled", func(t *testing.T) {
		for _, tc := range []struct {
			input string
			want  ParseErrorCode
		}{
			{
				input: `/* unexpected end of comment`,
				want:  ParseErrorUnexpectedEndOfComment,
			},
			{
				input: `"unexpected end of string`,
				want:  ParseErrorUnexpectedEndOfString,
			},
			{
				input: `2.`,
				want:  ParseErrorUnexpectedEndOfNumber,
			},
			{
				input: `"\u123"`,
				want:  ParseErrorInvalidUnicode,
			},
			{
				input: `"\."`,
				want:  ParseErrorInvalidEscapeCharacter,
			},
			{
				input: "\"\x01\"",
				want:  ParseErrorInvalidCharacter,
			},
			{
				input: `{"foo": "bar", /* this is a comment */}`,
				want:  noParseErrorCode,
			},
		} {
			t.Run(tc.want.String(), func(t *testing.T) {
				var have ParseErrorCode = noParseErrorCode
				v := Visitor{OnError: func(errorCode ParseErrorCode, offset, length int) {
					have = errorCode
				}}
				if !Walk(tc.input, ParseOptions{
					Comments:       true,
					TrailingCommas: true,
				}, v) {
					t.Error("Walk returned false unexpectedly")
				}

				if have != tc.want {
					t.Errorf("unexpected error code: have %v; want %v", have, tc.want)
				}
			})
		}
	})

	t.Run("comments disabled", func(t *testing.T) {
		for _, tc := range []struct {
			input string
			want  ParseErrorCode
		}{
			{
				input: `// line comment`,
				want:  InvalidCommentToken,
			},
			{
				input: `/* block comment */`,
				want:  InvalidCommentToken,
			},
			{
				input: `{"foo": "bar",}`,
				want:  noParseErrorCode,
			},
		} {
			t.Run(tc.want.String(), func(t *testing.T) {
				var have ParseErrorCode = noParseErrorCode
				v := Visitor{OnError: func(errorCode ParseErrorCode, offset, length int) {
					have = errorCode
				}}
				if !Walk(tc.input, ParseOptions{
					Comments:       false,
					TrailingCommas: true,
				}, v) {
					t.Error("Walk returned false unexpectedly")
				}

				if have != tc.want {
					t.Errorf("unexpected error code: have %v; want %v", have, tc.want)
				}
			})
		}
	})

	t.Run("trailing commas disabled", func(t *testing.T) {
		for _, tc := range []struct {
			input string
			want  ParseErrorCode
		}{
			{
				input: `{"foo": "bar"} // line comment`,
				want:  noParseErrorCode,
			},
			{
				input: `{"foo": "bar",}`,
				want:  ValueExpected,
			},
		} {
			t.Run(tc.want.String(), func(t *testing.T) {
				var have ParseErrorCode = noParseErrorCode
				v := Visitor{OnError: func(errorCode ParseErrorCode, offset, length int) {
					have = errorCode
				}}
				if !Walk(tc.input, ParseOptions{
					Comments:       true,
					TrailingCommas: false,
				}, v) {
					t.Error("Walk returned false unexpectedly")
				}

				if have != tc.want {
					t.Errorf("unexpected error code: have %v; want %v", have, tc.want)
				}
			})
		}
	})
}
