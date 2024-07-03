package main

import (
	//Import standard library
	"fmt"
	"os"
	"strings"

	//Import user's defined package
	"sentry/auth"
	"sentry/help"
)

func main() {
	if len(os.Args) == 1 {
		help.Help()
		return
	}

	if len(os.Args) != 1 {
		if strings.ToLower(os.Args[1]) == "crac" {
			auth.CreateAccount()
			return
		}

		if strings.ToLower(os.Args[1]) == "login" {
			auth.Login()
			return
		}

		if strings.ToLower(os.Args[1]) == "logout" {
			fmt.Println("|_ Log out successfully!")
			return
		}

		if strings.ToLower(os.Args[1]) == "add" {
			fmt.Println("Adding func")
		}
	}
}
