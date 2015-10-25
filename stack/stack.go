package stack

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/s3"
)

type AWSConnection struct {
	Config *aws.Config
	S3     *s3.S3
	Cf     *cloudformation.CloudFormation
}

func ConnectAWS(r string) *AWSConnection {
	new_aws := new(AWSConnection)
	new_aws.Config = aws.NewConfig().WithRegion(r)
	new_aws.S3 = s3.New(new_aws.Config)
	new_aws.Cf = cloudformation.New(new_aws.Config)
	return new_aws
}

func (a *AWSConnection) ValidateTemplate(i string) (*cloudformation.ValidateTemplateOutput, error) {
	input := &cloudformation.ValidateTemplateInput{
		TemplateBody: aws.String(i),
	}
	t, err := a.Cf.ValidateTemplate(input)
	if err != nil {
		return nil, err
	}
	return t, nil
}
