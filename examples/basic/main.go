// Example: Basic usage of the TempMailChecker Go SDK
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	tempmailchecker "github.com/Fushey/go-disposable-email-checker"
)

func main() {
	// Get API key from environment
	apiKey := os.Getenv("TEMPMAILCHECKER_API_KEY")
	if apiKey == "" {
		log.Fatal("TEMPMAILCHECKER_API_KEY environment variable is required")
	}

	// Create client with options
	checker, err := tempmailchecker.New(apiKey,
		tempmailchecker.WithEndpoint(tempmailchecker.EndpointEU),
		tempmailchecker.WithTimeout(10*time.Second),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Println("TempMailChecker Go SDK Example")
	fmt.Println("==============================\n")

	// Test emails
	emails := []string{
		"user@gmail.com",
		"test@10minutemail.com",
		"hello@yahoo.com",
		"fake@tempmail.org",
	}

	for _, email := range emails {
		result, err := checker.Check(email)
		if err != nil {
			if tempmailchecker.IsRateLimitError(err) {
				fmt.Println("\n⚠️  Rate limit reached! Try again later.")
				break
			}
			fmt.Printf("❌ Error checking %s: %v\n", email, err)
			continue
		}

		status := "✅ Legitimate"
		if result.Temp {
			status = "⚠️  Disposable"
		}
		fmt.Printf("%s: %s\n", email, status)
	}

	// Check a domain directly
	fmt.Println("\n--- Domain Check ---")
	domainResult, err := checker.CheckDomain("guerrillamail.com")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		if domainResult.Temp {
			fmt.Println("guerrillamail.com: ⚠️  Disposable domain")
		} else {
			fmt.Println("guerrillamail.com: ✅ Legitimate domain")
		}
	}

	// Check usage
	fmt.Println("\n--- API Usage ---")
	usage, err := checker.GetUsage()
	if err != nil {
		fmt.Printf("Error getting usage: %v\n", err)
	} else {
		fmt.Printf("Requests today: %d / %d\n", usage.UsageToday, usage.Limit)
		fmt.Printf("Resets at: %s\n", usage.Reset)
	}
}

