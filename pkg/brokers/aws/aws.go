package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type AWS struct {
	Config *AWSConfig
}

type AWSConfig struct {
	AwsRegion          string
	AwsAccessKeyId     string
	AwsSecretAccessKey string
}

func NewAWS(config *AWSConfig) *AWS {
	return &AWS{Config: config}
}

func (a *AWS) GetSession() *session.Session {
	awsConfig := &aws.Config{
		Region: &a.Config.AwsRegion,
	}

	if a.Config.AwsAccessKeyId != "" && a.Config.AwsSecretAccessKey != "" {
		// Getting credentials through environment variables
		awsConfig.Credentials = credentials.NewStaticCredentials(
			a.Config.AwsAccessKeyId, a.Config.AwsSecretAccessKey, "",
		)
	} else {
		// Getting credentials through host role
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		fmt.Println(err)
	}

	return sess
}
