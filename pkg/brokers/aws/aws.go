package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

type AWS struct {
	AccountId string
	Config    *AWSConfig
}

type AWSConfig struct {
	AwsRegion          string
	AwsAccessKeyId     string
	AwsSecretAccessKey string
	AwsSessionToken    string
}

func NewAWS(config *AWSConfig) (aws *AWS) {
	aws = &AWS{Config: config}
	aws.AccountId = aws.GetAccountId()
	return
}

func (a *AWS) GetAccountId() string {
	svc := sts.New(a.GetSession())
	input := &sts.GetCallerIdentityInput{}

	result, err := svc.GetCallerIdentity(input)
	if err != nil {
		fmt.Println(err)
	}
	return *result.Account
}

func (a *AWS) GetSession() *session.Session {
	awsConfig := &aws.Config{
		Region: &a.Config.AwsRegion,
	}

	if a.Config.AwsAccessKeyId != "" && a.Config.AwsSecretAccessKey != "" {
		// Getting credentials through environment variables
		awsConfig.Credentials = credentials.NewStaticCredentials(
			a.Config.AwsAccessKeyId,
			a.Config.AwsSecretAccessKey,
			a.Config.AwsSessionToken,
		)
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		log.Println(err)
	}

	return sess
}
