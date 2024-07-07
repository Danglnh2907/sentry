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
	"strconv"
	"strings"

	//Import user's defined package
	"sentry/utility"
)

type User struct {
	Username     string   `json:"username"`
	Password     string   `json:"password"`
	Fullname     string   `json:"fullname"`
	Email        string   `json:"email"`
	Budget       float64  `json:"budget"`
	Transactions []string `json:"transactions"`
}

func CreateAccount() {
	reader := bufio.NewReader(os.Stdin)
	isValid := false
	var (
		username, fullname, email, password string
		budget                              float64
		err                                 error
	)

	//Get username
	for !isValid {
		//Ask user for username
		fmt.Print("Enter your username: ")
		username, err = reader.ReadString('\n')
		if err != nil {
			utility.LogError(err, "Error at: CreateAccount -> Error reading username")
			return
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
			utility.LogError(err, "Error at: CreateAccount -> Error reading fullname")
			return
		}
		fullname = strings.TrimSpace(fullname)

		//Check if fullname is a non empty string
		isValid = len(fullname) > 0
		if !isValid {
			fmt.Println("Fullname cannot be empty!")
		}
	}

	isValid = false
	//Get email
	for !isValid {
		//Ask user for email
		fmt.Print("Enter your email: ")
		email, err = reader.ReadString('\n')
		if err != nil {
			utility.LogError(err, "Error at: CreateAccount -> Error reading email")
			return
		}
		email = strings.TrimSpace(email)

		//Check if fullname is a non empty string
		isValid = len(email) > 0
		if !isValid {
			fmt.Println("Email cannot be empty!")
			continue
		}

		//Check if email contain space and @ character or not
		isValid = strings.Contains(email, "@") && !strings.Contains(email, " ")
		if !isValid {
			fmt.Println(err, "Invalid email address")
		}
	}

	isValid = false
	//Get password
	for !isValid {
		//Ask user for password
		fmt.Print("Enter password: ")
		password, err = reader.ReadString('\n')
		if err != nil {
			utility.LogError(err, "Error at: CreateAccount -> Error reading password")
			return
		}
		password = strings.TrimSpace(password)

		//Check if password is valid
		isValid = len(password) >= 11
		if !isValid {
			fmt.Println("Password length must be at least 11 characters!")
		}

		isValid = !strings.Contains(password, " ")
		if !isValid {
			fmt.Println("Password must not contain any space!")
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
		budget, temp := 0.0, ""
		temp, err = reader.ReadString('\n')
		if err != nil {
			utility.LogError(err, "Error at: CreateAccount -> Error reading budget")
			return
		}

		//If temp is '\n', break from loop
		if temp = strings.TrimSpace(temp); temp == "" {
			isValid = true
			break
		}

		budget, err = strconv.ParseFloat(temp, 64)
		if err != nil {
			utility.LogError(err, "Error at: CreateAccount -> Error parsing budget")
			isValid = false
		}

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
		Email:        email,
		Budget:       budget,
		Transactions: make([]string, 0)}

	jsonData, err := json.MarshalIndent(newUser, "", " ")
	if err != nil {
		utility.LogError(err, "Error at: CreateAccount -> Error marshal user's info")
		return
	}

	resp, err := http.Post("http://localhost:8080/create-account", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		utility.LogError(err, "Error at: CreateAccount -> Error sending data to server")
		return
	}

	respMessage, err := io.ReadAll(resp.Body)
	if err != nil {
		utility.LogError(err, "Error at: CreateAccount -> Error reading respond body")
		return
	}
	defer resp.Body.Close()
	fmt.Printf("%s\n", respMessage)
}

func Login() (bool, string) {
	reader := bufio.NewReader(os.Stdin)

	//Get username
	fmt.Print("Enter your username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		utility.LogError(err, "Error at: Login -> Error reading username")
		return false, ""
	}
	username = strings.TrimSpace(username)

	//Get password
	fmt.Print("Enter your password: ")
	password, err := reader.ReadString('\n')
	if err != nil {
		utility.LogError(err, "Error at: Login -> Error reading password")
		return false, ""
	}
	password = strings.TrimSpace(password)

	//Send data to server for checking
	loginInfo := make([]map[string]string, 0)
	loginInfo = append(loginInfo, map[string]string{"username": username})
	loginInfo = append(loginInfo, map[string]string{"password": password})

	jsonData, err := json.MarshalIndent(loginInfo, "", " ")
	if err != nil {
		utility.LogError(err, "Error at: Login -> Error marshal user's info")
		return false, ""
	}

	resp, err := http.Post("http://localhost:8080/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		utility.LogError(err, "Error at: Login -> Error sending data to server")
		return false, ""
	}
	defer resp.Body.Close()

	message, err := io.ReadAll(resp.Body)
	if err != nil {
		utility.LogError(err, "Error at: Login -> Error reading message from respond body")
		return false, ""
	}

	fmt.Printf("%s\n", message)
	return resp.StatusCode == 202, username
}

func ShowProfile(username string) {
	//Sending request to server
	req, err := http.NewRequest("GET", "http://localhost:8080/get-profile", nil)
	if err != nil {
		utility.LogError(err, "Error at: ShowProfile -> Error making new request")
		return
	}
	req.Header.Set("Content-Type", "appliction/json")
	req.Header.Set("Identity", username)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utility.LogError(err, "Error at: ShowProfile -> Error sending request to server")
		return
	}
	defer resp.Body.Close()

	//Reading respond body and parsing data
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		utility.LogError(err, "Error at: ShowProfile -> Error reading respond body")
		return
	}

	var user User
	err = json.Unmarshal(data, &user)
	if err != nil {
		utility.LogError(err, "Error at: ShowProfile -> Error unmarshal respond body")
		return
	}

	//Display information
	result := fmt.Sprintf("Username: %s\nPassword: %s\nFullname: %s\nEmail: %s\nBudget: %.2f\n",
		user.Username,
		strings.Repeat("*", len(user.Password)),
		user.Fullname,
		user.Email,
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

func DeleteAccount(username string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Warning: deleting account will be permanent and all of your data will be lost!")
	fmt.Print("Enter your password to reconfirm: ")
	password, err := reader.ReadString('\n')
	if err != nil {
		utility.LogError(err, "Error at: DeleteAccount -> Error reading password")
		return
	}
	password = strings.TrimSpace(password)

	//Get password from server
	req, err := http.NewRequest("GET", "http://localhost:8080/get-profile", nil)
	if err != nil {
		utility.LogError(err, "Error at: ShowProfile -> Error making new request")
		return
	}
	req.Header.Set("Content-Type", "appliction/json")
	req.Header.Set("Identity", username)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utility.LogError(err, "Error at: ShowProfile -> Error sending request to server")
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		utility.LogError(err, "Error at: DeleteAccount -> Error reading respond body")
		return
	}

	var user User
	err = json.Unmarshal(data, &user)
	if err != nil {
		utility.LogError(err, "Error at: DeleteAccount -> Error unmarshal respond body")
		return
	}

	//Check if password match
	if password != user.Password {
		fmt.Println("Password incorrect!")
		return
	}

	//Delete account if password match
	req, err = http.NewRequest("DELETE", "http://localhost:8080/delete-profile", nil)
	if err != nil {
		utility.LogError(err, "Error at: DeleteAccount -> Error making new request")
		return
	}

	req.Header.Set("Identity", username)
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		utility.LogError(err, "Error at: DeleteAccount -> Error sending request to server")
		return
	}
	defer resp.Body.Close()

	//Print message to user
	message, err := io.ReadAll(resp.Body)
	if err != nil {
		utility.LogError(err, "Error at: DeleteAccount -> Error reading respond body")
		return
	}

	fmt.Println(string(message))
}
