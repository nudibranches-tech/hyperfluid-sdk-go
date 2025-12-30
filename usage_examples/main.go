package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🚀 Bifrost SDK - Fluent API Demo")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("⚠️  Warning: .env not loaded: %v\n", err)
	} else {
		fmt.Println("✓ .env file loaded successfully")
	}

	// TODO: progressive examples
	fmt.Println()
	fmt.Println("══════════════════════════════════════════════")
	fmt.Println("📊 FLUENT EXAMPLE")
	fmt.Println("══════════════════════════════════════════════")
	fmt.Println()

	runS3Example()

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🎉 All examples completed!")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
