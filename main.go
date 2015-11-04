package main

import (
	"fmt"
	"github.com/chooper/cfrun/stack"
	"github.com/chooper/cfrun/template"
	"log"
)

// TODO
// - Create stack
// - --profile flag
// - --delete-before-update
// - diffs

func main() {
	// TODO(charles) accept filename as an argument
	filename := "advanced.yaml"
	cf_json := template.ConvertToJSON(template.LoadYAML(filename))
	fmt.Printf("--- cf_json:\n%v\n\n", string(cf_json))

	// TODO(charles) remove hardcoded region name
	aws := stack.ConnectAWS("us-west-2")
	s, err := aws.ValidateTemplate(string(cf_json))
	if err != nil {
		log.Fatal(err) // print error and exit
	}
	log.Println(*s.Description) // output the templates description if specified.
}
