/*
Copyright © 2021 Bilal Bhatti
*/

package curly

func MapI2MapS(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = MapI2MapS(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = MapI2MapS(v)
		}
	}
	return i
}
