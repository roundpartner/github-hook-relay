package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"os"
)

func GetSession() *session.Session {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secret := os.Getenv("AWS_SECRET_ACCESS_KEY")
	sessionToken := os.Getenv("AWS_SESSION_TOKEN")
	region := os.Getenv("AWS_REGION")

	config := GetConfig(region, accessKey, secret, sessionToken)
	session, err := session.NewSession(config)

	if err != nil {
		return nil
	}
	return session
}

func GetConfig(region, accessKey, secret, sessionToken string) *aws.Config {
	if accessKey == "" {
		return ConfigWithIAM(region)
	}
	return ConfigWithAccessKey(region, accessKey, secret, sessionToken)
}

func ConfigWithAccessKey(region, accessKey, secret, sessionToken string) *aws.Config {
	return &aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secret, sessionToken),
	}
}

func ConfigWithIAM(region string) *aws.Config {
	return &aws.Config{
		Region: aws.String(region),
	}
}
