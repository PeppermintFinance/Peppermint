package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/plaid/plaid-go/v21/plaid"
)

// Helper function, gets key from .env file
func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

// Helper function, creates Plaid API client
func createPlaidClient() *plaid.APIClient {
	clientID := goDotEnvVariable("PLAID_CLIENT_ID")
	if clientID == "" {
		log.Fatal("PLAID_CLIENT_ID not found in environment variables")
	}

	secret := goDotEnvVariable("PLAID_PRODUCTION_SECRET")
	if secret == "" {
		log.Fatal("PLAID_PRODUCTION_SECRET not found in environment variables")
	}

	config := plaid.NewConfiguration()
	config.AddDefaultHeader("PLAID-CLIENT-ID", clientID)
	config.AddDefaultHeader("PLAID-SECRET", secret)
	config.UseEnvironment(plaid.Production)
	client := plaid.NewAPIClient(config)
	return client
}
