package auth

import (
	//Import standard library
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	//Import user's defined package
	"sentry/utility"
)

type User struct {
	Username     string   `json:"username"`
	Password     string   `json:"password"`
	Fullname     string   `json:"fullname"`
	Budget       float64  `json:"budget"`
	Transactions []string `json:"transactions"`
}

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
			utility.LogError(err, "Error reading input from user")
		}
		username = strings.TrimSpace(username)

		//Check if username is a non-empty string
		isValid = len(username) > 0
		if !isValid {
			fmt.Println("Username cannot be empty!")
			continue
		}
	}

	isValid = false
	//Get fullname
	for !isValid {
		//Ask user for fullname
		fmt.Print("Enter your fullname: ")
		fullname, err = reader.ReadString('\n')
		if err != nil {
			utility.LogError(err, "Error reading input from user")
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
			utility.LogError(err, "Error reading input from user")
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
		Username:     username,
		Password:     password,
		Fullname:     fullname,
		Transactions: make([]string, 0)}

	jsonData, err := json.MarshalIndent(newUser, "", " ")
	if err != nil {
		utility.LogError(err, "Error parsing for sending data to server")
	}

	resp, err := http.Post("http://localhost:8080/create-account", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		utility.LogError(err, "Error sending data to server")
	}

	respMessage, err := io.ReadAll(resp.Body)
	if err != nil {
		utility.LogError(err, "Error reading respond message")
	}
	fmt.Printf("%s\n", respMessage)

	defer resp.Body.Close()
}

func Login() (bool, string) {
	reader := bufio.NewReader(os.Stdin)

	//Get username
	fmt.Print("Enter your username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		utility.LogError(err, "Error reading input from user")
		return false, ""
	}
	username = strings.TrimSpace(username)

	//Get password
	fmt.Print("Enter your username: ")
	password, err := reader.ReadString('\n')
	if err != nil {
		utility.LogError(err, "Error reading input from user")
		return false, ""
	}
	password = strings.TrimSpace(password)

	//Send data to server for checking
	loginInfo := make([]map[string]string, 0)
	loginInfo = append(loginInfo, map[string]string{"username": username})
	loginInfo = append(loginInfo, map[string]string{"password": password})

	jsonData, err := json.MarshalIndent(loginInfo, "", " ")
	if err != nil {
		utility.LogError(err, "Error parsing data to json")
	}

	resp, err := http.Post("http://localhost:8080/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		utility.LogError(err, "Error sending data to server")
		return false, ""
	}
	defer resp.Body.Close()

	message, err := io.ReadAll(resp.Body)
	if err != nil {
		utility.LogError(err, "Error reading message from respond body")
		return false, ""
	}

	fmt.Printf("%s\n", message)
	return resp.StatusCode == 200, username
}

func ShowProfile(username string) {
	//Sending request to server
	req, err := http.NewRequest("GET", "http://localhost:8080/get-profile", nil)
	if err != nil {
		utility.LogError(err, "Error making request")
		return
	}
	req.Header.Set("Content-Type", "appliction/json")
	req.Header.Set("Identity", username)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utility.LogError(err, "Error sending request to server")
	}
	defer resp.Body.Close()

	//Reading respond body and parsing data
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		utility.LogError(err, "Error reading respond body")
		return
	}

	var user User
	err = json.Unmarshal(data, &user)
	if err != nil {
		utility.LogError(err, "Error parsing json data")
		return
	}

	//Display information
	result := fmt.Sprintf("Username: %s\nPassword: %s\nFullname: %s\nBudget: %.2f\n",
		user.Username,
		strings.Repeat("*", len(user.Password)),
		user.Fullname,
		user.Budget)
	if len(user.Transactions) == 0 {
		result += "Transactions: You currently have no transactions\n"
	} else {
		for _, val := range user.Transactions {
			result += fmt.Sprintf("%s ", val)
		}
		result += "\n"
	}
	fmt.Print(result)
}
