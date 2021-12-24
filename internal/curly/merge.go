package curly

import (
	"reflect"

	"github.com/pkg/errors"
)

func Merge(base, more interface{}) error {
	return merge(reflect.ValueOf(base), reflect.ValueOf(more))
}

func merge(base, more reflect.Value) error {
	for base.Kind() == reflect.Ptr || base.Kind() == reflect.Interface {
		base = base.Elem()
	}

	for more.Kind() == reflect.Ptr || more.Kind() == reflect.Interface {
		more = more.Elem()
	}

	if base.Kind() != more.Kind() {
		return errors.New("values not mergeable")
	}

	switch base.Kind() {
	case reflect.Map:
		err := mergeMap(base, more)
		if err != nil {
			return err
		}
	case reflect.Array, reflect.Slice:

	default:
		return errors.New("unexpected error")
	}

	return nil
}

func mergeMap(base, more reflect.Value) error {
	for _, k := range more.MapKeys() {
		left := base.MapIndex(k)

		for left.Kind() == reflect.Ptr || left.Kind() == reflect.Interface {
			left = left.Elem()
		}

		// left side does not have key
		if !left.IsValid() {
			base.SetMapIndex(k, more.MapIndex(k))
			continue
		}

		// left side has key
		right := more.MapIndex(k)

		// if left side a map merge map
		if _, ok := left.Interface().(map[string]interface{}); ok {
			err := merge(left, right)
			if err != nil {
				return err
			}

			continue
		}

		if _, ok := right.Interface().(map[string]interface{}); ok {
			// if left side is not a map, but right side is
			return errors.New("values not mergeable")
		}

		// if both sides are not maps
		base.SetMapIndex(k, right)
	}

	return nil
}
