package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func CreateAccount() {
	reader := bufio.NewReader(os.Stdin)
	isValid := false
	var (
		username, fullname, password string
		budget                       float64
		err                          error
	)

	//Get username
	for !isValid {
		//Ask user for username
		fmt.Print("Enter your username: ")
		username, err = reader.ReadString('\n')
		if err != nil {
			LogError(err, "Error reading input from user")
		}
		username = strings.TrimSpace(username)

		//Check if username is a non-empty string
		isValid = len(username) > 0
		if !isValid {
			fmt.Println("Username cannot be empty!")
			continue
		}

		//Get all usernames from server
		resp, err := http.Get("http://localhost:8080/get-username")
		if err != nil {
			LogError(err, "Error getting data from server")
		}
		defer resp.Body.Close()

		//Read data from respond body
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			LogError(err, "Error reading data from server")
		}
		var usernames []string
		err = json.Unmarshal(data, &usernames)
		if err != nil {
			LogError(err, "Error parsing data from server")
		}

		//Check if the input username exist or not
		for _, val := range usernames {
			if val == username {
				fmt.Println("Username already exist! Please try another one")
				isValid = false
				continue
			}
		}
	}

	isValid = false
	//Get fullname
	for !isValid {
		//Ask user for fullname
		fmt.Print("Enter your full name: ")
		fullname, err = reader.ReadString('\n')
		if err != nil {
			LogError(err, "Error reading input from user")
		}
		fullname = strings.TrimSpace(fullname)

		//Check if fullname is a non empty string
		isValid = len(fullname) > 0
		if !isValid {
			fmt.Println("Fullname cannot be empty!")
		}
	}

	isValid = false
	//Get password
	for !isValid {
		//Ask user for password
		fmt.Print("Enter password: ")
		password, err = reader.ReadString('\n')
		if err != nil {
			LogError(err, "Error reading input from user")
		}
		password = strings.TrimSpace(password)

		//Check if password is valid
		isValid = len(password) >= 11
		if !isValid {
			fmt.Println("Password length must be at least 11 characters")
		}

		regex, _ := regexp.Compile(".*[A-Z]+.*")
		isValid = regex.MatchString(password)
		if !isValid {
			fmt.Println("Password must contain at least one uppercase character!")
		}

		regex, _ = regexp.Compile(".*[a-z]+.*")
		isValid = regex.MatchString(password)
		if !isValid {
			fmt.Println("Password must contain at least one lowercase character!")
		}

		regex, _ = regexp.Compile(".*[0-9]+.*")
		isValid = regex.MatchString(password)
		if !isValid {
			fmt.Println("Password must contain at least one number!")
		}

		regex, _ = regexp.Compile(".*[^A-Za-z0-9]+.*")
		isValid = regex.MatchString(password)
		if !isValid {
			fmt.Println("Password must contain at least one special character!")
		}
	}

	isValid = false
	//Get budget
	for !isValid {
		//Ask for budget
		fmt.Print("Enter your budget. If you don't want to add budget now, press ENTER. You can update later ")
		budget = 0.0
		fmt.Scanf("%d\n", &budget)

		//Check if budget is valid
		isValid = budget >= 0
		if !isValid {
			fmt.Println("Budget cannot be a negative number!")
		}
	}

	//Send data back to server
	newUser := User{
		Username: username,
		Password: password,
		Fullname: fullname,
		Budget:   budget, PreferCurrency: "USD",
		Transactions: make([]Transaction, 0)}

	jsonData, err := json.MarshalIndent(newUser, "", " ")
	if err != nil {
		LogError(err, "Error parsing for sending data to server")
	}
	resp, err := http.Post("http://localhost:8080/user", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		LogError(err, "Error sending data to server")
	}
	defer resp.Body.Close()

	//Print message to user
	fmt.Println("Account created successfully! Please log in to use our services!")
}
