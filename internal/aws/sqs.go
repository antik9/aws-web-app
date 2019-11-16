package awsapp

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var sqsInstance *sqs.SQS
var queueUrl *string

func init() {
	sess := session.Must(
		session.NewSession(
			&aws.Config{
				Region: aws.String(endpoints.EuNorth1RegionID),
			},
		),
	)
	sqsInstance = sqs.New(sess)
	getQueueUrl(aws.String("IpQueue"))
}

func getQueueUrl(queueName *string) {
	queueUrlOutput, err := sqsInstance.GetQueueUrl(&sqs.GetQueueUrlInput{QueueName: queueName})
	if err != nil {
		log.Fatal(err)
	}
	queueUrl = queueUrlOutput.QueueUrl
}

func SendMessage(message *string) {
	_, err := sqsInstance.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    queueUrl,
		MessageBody: message,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func ReadMessage() *string {
	for {
		maxNumberOfMessages := int64(1)
		input := &sqs.ReceiveMessageInput{
			MaxNumberOfMessages: &maxNumberOfMessages,
			QueueUrl:            queueUrl,
		}
		receiveMessageOutput, err := sqsInstance.ReceiveMessage(input)
		if err == nil && len(receiveMessageOutput.Messages) == 1 {
			message := receiveMessageOutput.Messages[0]

			_, err = sqsInstance.DeleteMessage(&sqs.DeleteMessageInput{
				QueueUrl:      queueUrl,
				ReceiptHandle: message.ReceiptHandle,
			})
			if err == nil {
				return message.Body
			}
		}
		time.Sleep(time.Second)
	}
}
