package help

import (
	"fmt"
	"os"
	"sentry/utility"

	"github.com/joho/godotenv"
)

func Welcome() {
	err := godotenv.Load()
	if err != nil {
		utility.LogError(err, "Error loading env file")
		return
	}

	if os.Getenv("state") == "true" {
		fmt.Printf("Welcome, %s\n", os.Getenv("user"))
	}
	fmt.Println("This is Sentry, a CLI-based finance management application")
	fmt.Println("If you need some help on how to use this app, enter ./sentry help-command")
}
