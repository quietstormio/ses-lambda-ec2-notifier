package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type EC2StateChangeDetail struct {
	InstanceID string `json:"instance-id"`
	State      string `json:"state"`
}

type CloudWatchEvent struct {
	Source     string               `json:"source"`
	DetailType string               `json:"detail-type"`
	Detail     EC2StateChangeDetail `json:"detail"`
}

func Handler(ctx context.Context, event CloudWatchEvent) (string, error) {
	// Configure AWS session
	session, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		log.Println("Failed creating AWS session:", err)
		return "", err
	}

	// Initiate new SES session
	svc := ses.New(session)

	fromAddress := "fromexample@gmail.com"

	toAddress := "toexample@gmail.com"

	textBody := "AN EC2 instanace is ready: " + event.Detail.InstanceID + " is now in state: " + event.Detail.State

	htmlBody := "<h1>An EC2 instance has entered the running state</h1><p>Instance ID: " + event.Detail.InstanceID + "</p>"

	charSet := "UTF-8"

	emailInput := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(toAddress),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(htmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(textBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(charSet),
				Data:    aws.String("EC2 instance ready"),
			},
		},
		Source: aws.String(fromAddress),
	}

	result, err := svc.SendEmail(emailInput)
	if err != nil {
		log.Println("Failed sending email:", err)
		return "", err
	}

	log.Println("Sent email to: ", toAddress)
	return result.String(), nil
}

func main() {
	lambda.Start(Handler)
}
