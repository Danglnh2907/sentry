package main

import (
	"os"
	"strings"
)

var State bool = false

func main() {
	if len(os.Args) == 1 {
		Help()
		return
	}

	if len(os.Args) != 1 {
		if strings.ToLower(os.Args[1]) == "crac" {
			CreateAccount()
			return
		}

		if strings.ToLower(os.Args[1]) == "login" {
			Login()
			return
		}
	}
}
