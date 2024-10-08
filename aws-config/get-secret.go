// Package awsconfig contains utilities for retrieving secrets from AWS Secrets Manager.
package awsconfig

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// GetSecret retrieves a secret from Secrets Manager by ID and decodes the JSON value.
//
// It takes a single parameter, the secret ID, and returns the decoded secret map
// and any error.
//
// Example:
//
//	secret, err := GetSecret("mySecretId")
//		if err != nil {
//		log.Fatal(err)
//		}
//
// password := secret["password"]
func GetSecret(secretId string) (map[string]string, error) {

	// Configure SDK and create client
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("REGION")))
	if err != nil {
		panic("configuration error")
	}
	client := secretsmanager.NewFromConfig(cfg)

	// Call GetSecretValue
	output, err := client.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: &secretId,
	})

	if err != nil {
		return nil, err
	}

	// Decode JSON secret
	var secret map[string]string
	if err := json.Unmarshal([]byte(*output.SecretString), &secret); err != nil {
		return nil, err
	}
	return secret, nil
}
