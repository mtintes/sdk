package flatmap_test

import (
	"encoding/json"
	"fmt"

	"github.com/nextmv-io/sdk/flatmap"
)

func ExampleDo() {
	nested := map[string]any{
		"a": "foo",
		"b": []any{
			map[string]any{
				"c": "bar",
				"d": []any{
					map[string]any{
						"e": 2,
					},
					true,
				},
			},
			map[string]any{
				"c": "baz",
				"d": []any{
					map[string]any{
						"e": 3,
					},
					false,
				},
			},
		},
	}

	flattened := flatmap.Do(nested)

	b, err := json.MarshalIndent(flattened, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))

	// Output:
	// {
	//   "$.a": "foo",
	//   "$.b[0].c": "bar",
	//   "$.b[0].d[0].e": 2,
	//   "$.b[0].d[1]": true,
	//   "$.b[1].c": "baz",
	//   "$.b[1].d[0].e": 3,
	//   "$.b[1].d[1]": false
	// }
}

func ExampleUndo() {
	flattened := map[string]any{
		"$.a":           "foo",
		"$.b[0].c":      "bar",
		"$.b[0].d[0].e": 2,
		"$.b[0].d[1]":   true,
		"$.b[1].c":      "baz",
		"$.b[1].d[0].e": 3,
		"$.b[1].d[1]":   false,
	}

	nested, err := flatmap.Undo(flattened)
	if err != nil {
		panic(err)
	}

	b, err := json.MarshalIndent(nested, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
	// Output:
	// {
	//   "a": "foo",
	//   "b": [
	//     {
	//       "c": "bar",
	//       "d": [
	//         {
	//           "e": 2
	//         },
	//         true
	//       ]
	//     },
	//     {
	//       "c": "baz",
	//       "d": [
	//         {
	//           "e": 3
	//         },
	//         false
	//       ]
	//     }
	//   ]
	// }
}
