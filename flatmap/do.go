package flatmap

import (
	"fmt"
	"reflect"
)

/*
Do takes a nested map and flattens it into a single level map. The flattening
roughly follows the [JSONPath] standard. Please see the example to understand
how the flattened output looks like.

[JSONPath]: https://goessner.net/articles/JsonPath/
*/
func Do(nested map[string]any) map[string]any {
	flattened := map[string]any{}
	for childKey, childValue := range nested {
		setChildren(flattened, childKey, childValue)
	}

	return flattened
}

// setChildren is a helper function for flatten. It is invoked recursively on a
// child value. If the child is not a map or a slice, then the value is simply
// set on the flattened map. If the child is a map or a slice, then the
// function is invoked recursively on the child's values, until a
// non-map-non-slice value is hit.
func setChildren(flattened map[string]any, parentKey string, parentValue any) {
	newKey := fmt.Sprintf(".%s", parentKey)
	if reflect.TypeOf(parentValue) == nil {
		flattened[newKey] = parentValue
		return
	}

	if reflect.TypeOf(parentValue).Kind() == reflect.Map {
		children := parentValue.(map[string]any)
		for childKey, childValue := range children {
			newKey = fmt.Sprintf("%s.%s", parentKey, childKey)
			setChildren(flattened, newKey, childValue)
		}
		return
	}

	if reflect.TypeOf(parentValue).Kind() == reflect.Slice {
		children := parentValue.([]any)
		if len(children) == 0 {
			flattened[newKey] = children
			return
		}

		for childIndex, childValue := range children {
			newKey = fmt.Sprintf("%s[%v]", parentKey, childIndex)
			setChildren(flattened, newKey, childValue)
		}
		return
	}

	flattened[newKey] = parentValue
}
