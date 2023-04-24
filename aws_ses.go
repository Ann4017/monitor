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
	s_recipient  []string
	s_subject    string
	s_body       string
}

func (c *C_ses) Init(_s_region, _s_access_key, _s_secret_key string) {
	c.s_region = _s_region
	c.s_access_Key = _s_access_key
	c.s_secret_key = _s_secret_key
}

func (c *C_ses) Write_email(_s_sender string, _s_recipient []string, _s_subject, _s_body string) {
	c.s_sender = _s_sender
	c.s_recipient = _s_recipient
	c.s_subject = _s_subject
	c.s_body = _s_body
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

func (c *C_ses) Send_email(_pc_client *ses.Client, _s_sender string, _s_recipient []string, _s_subject string, _s_body string) error {
	input := ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: _s_recipient,
		},
		Message: &types.Message{
			Subject: &types.Content{
				Data: aws.String(_s_subject),
			},
			Body: &types.Body{
				Text: &types.Content{
					Data: aws.String(_s_body),
				},
			},
		},
		Source: aws.String(_s_sender),
	}

	_, err := c.pc_client.SendEmail(context.Background(), &input)
	if err != nil {
		return err
	}

	return nil
}
