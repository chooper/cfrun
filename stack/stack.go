package stack

import (
	"bytes"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/s3"
)

type AWSConnection struct {
	Config *aws.Config
}

func ConnectAWS(r string) *AWSConnection {
	new_aws := new(AWSConnection)
	new_aws.Config = aws.NewConfig().WithRegion(r)
	return new_aws
}

func (a *AWSConnection) ValidateTemplate(i []byte) error {
	cf := cloudformation.New(a.Config)
	input := &cloudformation.ValidateTemplateInput{
		TemplateBody: aws.String(string(i)),
	}

	_, err := cf.ValidateTemplate(input)
	return err
}

func (a *AWSConnection) UploadTemplate(b string, k string, t []byte) error {
	s3_ := s3.New(a.Config)
	input := &s3.PutObjectInput{
		Bucket: aws.String(b),
		Key:    aws.String(k),
		Body:   bytes.NewReader([]byte(t)),
	}
	_, err := s3_.PutObject(input)
	return err
}

func (a *AWSConnection) CreateStack(b string, k string, s string) (*string, error) {
	cf := cloudformation.New(a.Config)
	input := &cloudformation.CreateStackInput{
		StackName:   aws.String(s),
		TemplateURL: aws.String("https://s3.amazonaws.com/" + b + "/" + k),
	}
	output, err := cf.CreateStack(input)
	return output.StackId, err
}

func (a *AWSConnection) GetStackStatus(sid *string) (*string, error) {
	cf := cloudformation.New(a.Config)
	input := &cloudformation.DescribeStacksInput{
		StackName: aws.String(*sid),
	}
	output, err := cf.DescribeStacks(input)
	if len(output.Stacks) != 1 {
		err = errors.New("DescribeStacks output does not have exactly 1 stacks")
	}
	if err != nil {
		return nil, err
	}
	stack := output.Stacks[0]
	return stack.StackStatus, nil
}
