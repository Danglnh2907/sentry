package main

import (
	//Import standard library
	"fmt"
	"os"
	"strings"

	//Import user's defined package
	"sentry/auth"
	"sentry/help"
	"sentry/utility"

	"github.com/joho/godotenv"
)

func main() {
	//Create new .env file to store the login state of current device
	utility.CreateNewFile(".env")
	//Load env var from .env file
	err := godotenv.Load(".env")
	if err != nil {
		utility.LogError(err, "Error loading enviroment variables!")
		return
	}

	if len(os.Args) == 1 {
		help.Help()
		return
	}

	if len(os.Args) != 1 {
		if strings.ToLower(os.Args[1]) == "crac" {
			auth.CreateAccount()
		}

		if strings.ToLower(os.Args[1]) == "login" {
			isLogin, username := auth.Login()
			if isLogin {
				utility.SetEnvVar("true", username)
			}
		}

		if strings.ToLower(os.Args[1]) == "logout" {
			//utility.SetEnvVar("false", "")
			fmt.Println("Log out successfully!")
		}

		if strings.ToLower(os.Args[1]) == "add" {
			if os.Getenv("state") == "true" {
				fmt.Println("Adding func")
			} else {
				fmt.Println("You must log in to use this function")
			}
		}
	}

	//Write the enviroment variable back to .env file
	envMap := map[string]string{
		"state": os.Getenv("state"),
		"user":  os.Getenv("user"),
	}
	godotenv.Write(envMap, ".env")
}
