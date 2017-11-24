package jsonx

import (
	"fmt"
)

func ExampleEdit() {
	const input = `
/* comment */
{
  "a": 1 // oops! forgot a comma
  /* note the trailing comma */
  "b": 2,
}`

	// Insert value 3 at key path c/d.
	edits, _, _ := ComputePropertyEdit(input,
		PropertyPath("c", "d"),
		3,
		nil,
		FormatOptions{InsertSpaces: true, TabSize: 2},
	)
	output, _ := ApplyEdits(input, edits...)
	fmt.Println(output)
	// Output: /* comment */
	// {
	//   "a": 1 // oops! forgot a comma
	//   /* note the trailing comma */
	//   "b": 2,
	//   "c": {
	//     "d": 3
	//   },
	// }
}
