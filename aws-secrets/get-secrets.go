package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func parseJSONString(jsonStr string) (map[string]string, error) {
	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &jsonData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	result := make(map[string]string)
	for k, v := range jsonData {
		strValue, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("value for key '%s' is not a string", k)
		}
		result[k] = strValue
	}

	return result, nil
}

// func writeEnvVars()

func main() {
	secretName := os.Args[1]
	// pathToWrite := os.Args[2]

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-2"),
	)
	if err != nil {
		log.Fatalf("Failed to load AWS configuration: %v", err)
	}

	// Create Secrets Manager client
	smClient := secretsmanager.NewFromConfig(cfg)

	// Get the secret value
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := smClient.GetSecretValue(context.TODO(), input)
	if err != nil {
		log.Fatalf("Failed to get secret value: %v", err)
	}

	// Extract and use the secret value
	secretString := *result.SecretString
	// fmt.Printf("Retrieved secret value: %s\n", secretString)
	parsedData, err := parseJSONString(secretString)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	for k, v := range parsedData {
		fmt.Printf("export %s=%s\n", k, v)
	}

	// Use the secret value in your application...
}
