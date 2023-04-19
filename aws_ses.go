package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type C_ses struct {
	s_region     string
	s_access_Key string
	s_secret_key string
	pc_client    *ses.Client
	s_sender     string
	s_recipient  string
	s_subject    string
	s_body       string
}

func (c *C_ses) Init(region, access_key, secret_key string) {
	c.s_region = region
	c.s_access_Key = access_key
	c.s_secret_key = secret_key
}

func (c *C_ses) Write_email(sender, recipient, subject, body string) {
	c.s_sender = sender
	c.s_recipient = recipient
	c.s_subject = subject
	c.s_body = body
}

func (c *C_ses) Set_config() error {
	cred := credentials.NewStaticCredentialsProvider(c.s_access_Key, c.s_secret_key, "")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithDefaultRegion(c.s_region), config.WithCredentialsProvider(cred))
	if err != nil {
		return err
	}

	c.pc_client = ses.NewFromConfig(cfg)

	return nil
}

func (c *C_ses) Send_email(client *ses.Client, sender, recipient, subject, body string) error {
	input := ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{recipient},
		},
		Message: &types.Message{
			Subject: &types.Content{
				Data: aws.String(subject),
			},
			Body: &types.Body{
				Text: &types.Content{
					Data: aws.String(body),
				},
			},
		},
		Source: aws.String(sender),
	}

	_, err := c.pc_client.SendEmail(context.Background(), &input)
	if err != nil {
		return err
	}

	return nil
}
