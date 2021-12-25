/*
Copyright © 2021 Bilal Bhatti
*/

package curly

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type test struct {
	base, more string
	expected   string
}

func TestMerge(t *testing.T) {
	tests := []test{
		{
			// fail test
			// different types
			base: `{"x":"y"}`,
			more: `["one"]`,
		},
		{
			// fail test
			// value of key [a] is different type
			base: `{"x":"y","a":{"d":"e"}}`,
			more: `{"a":"b"}`,
		},
		{
			// fail test
			// value of key [a] is different type
			base: `{"a":"b"}`,
			more: `{"x":"y","a":{"d":"e"}}`,
		},
		{
			// simple types, merge test
			base:     `{}`,
			more:     `{"a":"b"}`,
			expected: `{"a":"b"}`,
		},
		{
			// simple types, merge test
			base:     `{"x":"y"}`,
			more:     `{"a":"b"}`,
			expected: `{"a":"b","x":"y"}`,
		},
		{
			// simple types, replace test
			base:     `{"x":"y","a":"c"}`,
			more:     `{"a":"b"}`,
			expected: `{"a":"b","x":"y"}`,
		},
		{
			// map merge test
			base:     `{"x":"y","a":{"d":"e"}}`,
			more:     `{"a":{"b":"c"}}`,
			expected: `{"a":{"b":"c","d":"e"},"x":"y"}`,
		},
		{
			// map merge test
			base:     `{"x":"y","a":{"d":"e","f":{"g":"h"}}}`,
			more:     `{"a":{"b":"c","f":{"i":"j"}}}`,
			expected: `{"a":{"b":"c","d":"e","f":{"g":"h","i":"j"}},"x":"y"}`,
		},
		{
			// replace array test
			base:     `{"x":"y","a":["d","e"]}`,
			more:     `{"a":["b","c"]}`,
			expected: `{"a":["b","c"],"x":"y"}`,
		},
	}

	for ti, test := range tests {
		var base, more interface{}

		mustMarshall(t, test.base, &base)

		mustMarshall(t, test.more, &more)

		if err := Merge(base, more); err != nil {
			if test.expected == "" {
				assert.NotNil(t, err)
			} else {
				t.Errorf("test[%d] failure, expected: `%s`, %v", ti, test.expected, err)
			}
			continue
		}

		var expected interface{}
		mustMarshall(t, test.expected, &expected)

		if !reflect.DeepEqual(base, expected) {
			bites, _ := json.Marshal(base)
			t.Error(
				fmt.Sprintf("test[%d] json should be equal\n", ti),
				fmt.Sprintf("expect: %v\n", test.expected),
				fmt.Sprintf("actual: %v\n", string(bites)),
			)
		}
	}
}

func mustMarshall(t *testing.T, js string, target interface{}) {
	if err := json.Unmarshal([]byte(js), &target); err != nil {
		t.Errorf("failed to unmarshal `%s` %v", js, err)
	}
}
