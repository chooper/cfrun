package main

import (
	"github.com/chooper/cfrun/stack"
	"github.com/chooper/cfrun/template"
	"log"
	"time"
)

// TODO
// - Create stack
// - --profile flag
// - --delete-before-update
// - diffs

func main() {
	// TODO(charles) accept these as arguments
	filename := "dummy.yaml"
	s3_bucket := "cch-test"
	s3_key := "cf.json"
	region := "us-west-2"
	stack_name := "my-cf-stack"

	cf_json := template.ConvertToJSON(template.LoadYAML(filename))

	aws := stack.ConnectAWS(region)
	err := aws.ValidateTemplate(cf_json)
	if err != nil {
		log.Fatal(err)
	}

	err = aws.UploadTemplate(s3_bucket, s3_key, cf_json)
	if err != nil {
		log.Fatal(err)
	}

	stack_id, err := aws.CreateStack(s3_bucket, s3_key, stack_name)
	if err != nil {
		log.Fatal(err)
	}

	for {
		status, err := aws.GetStackStatus(stack_id)
		if err != nil {
			log.Fatal(err)
		}
		// TODO(charles) there are other statuses we should stop for
		log.Printf("status: %v", *status)
		if *status == "CREATE_COMPLETE" {
			break
		}
		time.Sleep(10 * time.Second)
	}
}
