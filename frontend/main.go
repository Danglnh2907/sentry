package main

import (
	//Import standard library
	"fmt"
	"os"
	"strings"

	//Import user's defined package
	"sentry/auth"
	"sentry/crud"
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
		help.Welcome()
		return
	}

	if len(os.Args) != 1 {
		//Create account
		if strings.ToLower(os.Args[1]) == "crac" {
			auth.CreateAccount()
		}

		//Login
		if strings.ToLower(os.Args[1]) == "login" {
			if os.Getenv("state") == "true" {
				fmt.Println("You have already logged in")
			} else {
				isLogin, username := auth.Login()
				if isLogin {
					utility.SetEnvVar("true", username)
				}
			}
		}

		//Logout
		if strings.ToLower(os.Args[1]) == "logout" {
			utility.SetEnvVar("false", "")
			fmt.Println("Log out successfully!")
		}

		//Show profile
		if strings.ToLower(os.Args[1]) == "show-profile" {
			if os.Getenv("state") == "false" {
				fmt.Println("You must log in to use this function")
				return
			}
			auth.ShowProfile(os.Getenv("user"))
		}

		//Add transactions
		if strings.ToLower(os.Args[1]) == "add" {
			if os.Getenv("state") == "true" {
				if len(os.Args) == 3 {
					crud.AddTransByFile(os.Args[2])
				} else {
					crud.AddTrans(os.Getenv("user"))
				}
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
