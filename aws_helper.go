package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
	"os"
)

func GetSession() *session.Session {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secret := os.Getenv("AWS_SECRET_ACCESS_KEY")
	sessionToken := os.Getenv("AWS_SESSION_TOKEN")
	region := os.Getenv("AWS_REGION")

	session, err := session.NewSession(
		&aws.Config{
			Region:      aws.String(region),
			Credentials: credentials.NewStaticCredentials(accessKey, secret, sessionToken),
		})
	if err != nil {
		log.Printf("AWS Error: %s\n", err.Error())
		return nil
	}
	return session
}
