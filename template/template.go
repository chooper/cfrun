package template

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strconv"
)

// stolen: https://github.com/bronze1man/yaml2json/blob/master/main.go
func transformData(in interface{}) (out interface{}, err error) {
	switch in.(type) {
	case map[interface{}]interface{}:
		o := make(map[string]interface{})
		for k, v := range in.(map[interface{}]interface{}) {
			sk := ""
			switch k.(type) {
			case string:
				sk = k.(string)
			case int:
				sk = strconv.Itoa(k.(int))
			default:
				return nil, errors.New(
					fmt.Sprintf("type not match: expect map key string or int get: %T", k))
			}
			v, err = transformData(v)
			if err != nil {
				return nil, err
			}
			o[sk] = v
		}
		return o, nil
	case []interface{}:
		in1 := in.([]interface{})
		len1 := len(in1)
		o := make([]interface{}, len1)
		for i := 0; i < len1; i++ {
			o[i], err = transformData(in1[i])
			if err != nil {
				return nil, err
			}
		}
		return o, nil
	default:
		return in, nil
	}
}

func Load(filename string) []byte {
	file_contents, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var file_structure interface{}
	err = yaml.Unmarshal([]byte(file_contents), &file_structure)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	file_structure, err = transformData(file_structure)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	cf_json, err := json.Marshal(file_structure)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return cf_json
}
