package aws

import (
	"context"
	"fmt"

	brk "github.com/nanernunes/federation/pkg/brokers"
	"github.com/nanernunes/federation/pkg/util"

	"github.com/aws/aws-sdk-go/service/sns"
)

type SNS struct {
	Name   string
	AWS    *AWS
	Client *sns.SNS
	Errors chan error
}

func NewSNS(name string, aws *AWS) *SNS {
	svc := sns.New(aws.GetSession())
	return &SNS{Name: name, AWS: aws, Client: svc}
}

func (s *SNS) GetName() string {
	return s.Name
}

func (s *SNS) Connect(chan error) bool {
	return false
}

func (s *SNS) Ack(message *brk.Message) error {
	return nil
}

func (s *SNS) Subscribe(ctx context.Context, source string, chErr chan error) <-chan brk.Message {
	return nil
}

func (s *SNS) Publish(
	target string, message *brk.Message, options map[string]interface{},
) (string, error) {

	dataType := "String"
	attributes := make(map[string]*sns.MessageAttributeValue)

	for key, value := range message.Headers {
		if parsed, err := util.ToString(value); err == nil {
			attributes[key] = &sns.MessageAttributeValue{
				DataType:    &dataType,
				StringValue: &parsed,
			}
		}
	}

	arn := s.ToTopicArn(target)

	result, err := s.Client.Publish(&sns.PublishInput{
		TopicArn:          &arn,
		Message:           &message.Body,
		MessageAttributes: attributes,
	})

	if err != nil {
		return "", err
	}

	return *result.MessageId, nil
}

func (s *SNS) ToTopicArn(topic string) string {
	return fmt.Sprintf(
		"arn:aws:sns:%s:%s:%s",
		s.AWS.Config.AwsRegion,
		s.AWS.AccountId,
		topic,
	)
}
