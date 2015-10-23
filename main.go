package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/chooper/cfrun/template"
	"log"
)

// TODO
// - Create stack
// - --region flag
// - --profile flag
// - Variable file names
// - --delete-before-update
// - diffs

func main() {
	filename := "advanced.yaml"
	cf_json := template.Load(filename)
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
