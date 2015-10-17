package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"strconv"
)

// TODO
// - Create stack
// - --region flag
// - --profile flag
// - Variable file names
// - --delete-before-update
// - diffs

const BUFFER_SIZE = 65536

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
	return in, nil
}

func loadTemplate(filename string) []byte {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	file_contents := make([]byte, BUFFER_SIZE)
	_, err = file.Read(file_contents)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	file_contents = bytes.Trim(file_contents, "\x00")

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

func main() {
	filename := "advanced.yaml"
	cf_json := loadTemplate(filename)
	fmt.Printf("--- cf_json:\n%v\n\n", string(cf_json))

	config := aws.NewConfig().WithRegion("us-west-2")
	svc := cloudformation.New(config)
	input := &cloudformation.ValidateTemplateInput{
		TemplateBody: aws.String(string(cf_json)),
	}
	template, err := svc.ValidateTemplate(input)
	if err != nil {
		log.Fatal(err) // print error and exit
	}
	log.Println(*template.Description) // output the templates description if specified.

}
