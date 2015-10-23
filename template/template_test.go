package template

import (
	"reflect"
	"testing"
)

type testCase struct {
	map_key       interface{}
	map_value     interface{}
	expected_json []byte
}

func TestConvertToJSONDoesntBlowUp(t *testing.T) {
	// declare test cases
	testCases := []testCase{
		testCase{map_key: "foo", map_value: 6, expected_json: []byte("{\"foo\":6}")},
		testCase{map_key: "bar", map_value: []int{6}, expected_json: []byte("{\"bar\":[6]}")},
		testCase{map_key: "fu", map_value: "bum", expected_json: []byte("{\"fu\":\"bum\"}")},
		testCase{map_key: "baz", map_value: []string{"boom"}, expected_json: []byte("{\"baz\":[\"boom\"]}")},
	}

	var expected_json []byte
	var actual_json []byte

	// do the testing
	for idx := range testCases {
		var data_to_convert = make(map[interface{}]interface{})
		tcase := testCases[idx]
		expected_json = tcase.expected_json
		data_to_convert[tcase.map_key] = tcase.map_value

		actual_json = ConvertToJSON(data_to_convert)
		if !reflect.DeepEqual(actual_json, expected_json) {
			t.Errorf("Wanted %v but got %v", string(expected_json), string(actual_json))
		}
	}
}
