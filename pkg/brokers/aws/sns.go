package aws

import (
	"context"

	brk "github.com/nanernunes/federation/pkg/brokers"

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

	attributes := make(map[string]*sns.MessageAttributeValue)

	for key, value := range message.Headers {
		dataType := "string"
		stringValue := value.(string)

		attributes[key] = &sns.MessageAttributeValue{
			DataType:    &dataType,
			StringValue: &stringValue,
		}
	}

	result, err := s.Client.Publish(&sns.PublishInput{
		TopicArn:          &target,
		Message:           &message.Body,
		MessageAttributes: attributes,
	})

	if err != nil {
		return "", err
	}

	return *result.MessageId, nil
}
