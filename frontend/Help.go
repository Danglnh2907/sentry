package main

import (
	"fmt"
)

func Help() {
	fmt.Println("This is Sentry, a CLI-based finance management application")
	if !State {
		fmt.Println("It seems like you haven't logged in! To log in, enter ./sentry login")
		fmt.Println("If you haven't had any account, create account by ./sentry crac")
	}
	fmt.Println("If you need some help on how to use this app, enter ./sentry help-command")
}
