package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Cat struct {
	Fact   string `json:"fact"`
	Length int    `json:"length"`
}

type Response struct {
	StatusCode int    `json:"statuscode"`
	Message    string `json:"message"`
	Cat
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("Hello Handler with Name", req.Body)

	// Make request to third-party API
	resp, err := http.Get("https://catfact.ninja/fact")
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Failed to call third-party API"}, err
	}
	defer resp.Body.Close()

	// Decode response body
	var cat Cat
	if err := json.NewDecoder(resp.Body).Decode(&cat); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Failed to decode API response"}, err
	}

	// Construct Lambda response
	response := Response{
		StatusCode: resp.StatusCode,
		Cat:        cat,
		Message:    "Successful",
	}

	// Serialize Lambda response
	body, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Failed to serialize response"}, err
	}

	// Upload response to S3
	err = uploadToS3(body)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Failed to upload to S3"}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: resp.StatusCode,
		Body:       string(body),
	}, nil
}

func uploadToS3(data []byte) error {
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create S3 service client
	svc := s3.New(sess)

	// Generate a unique object key based on the current timestamp
	objectKey := time.Now().Format("2006-01-02T15:04:05Z07:00")

	// Define the parameters of the object
	params := &s3.PutObjectInput{
		Bucket: aws.String("my-lambda-data"),
		Key:    aws.String(objectKey),
		Body:   aws.ReadSeekCloser(bytes.NewReader(data)), // Use bytes.NewReader to create a ReadSeekCloser from the []byte
	}

	// Upload object to S3
	_, err := svc.PutObject(params)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
