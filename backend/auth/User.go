package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sentry/utility"
)

// User model
type User struct {
	Username     string   `json:"username"`
	Password     string   `json:"password"`
	Fullname     string   `json:"fullname"`
	Budget       float64  `json:"budget"`
	Transactions []string `json:"transactions"`
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	//Get the list of usernames in database
	data, err := utility.OpenFile("data/usernames.json")
	if err != nil {
		utility.LogError(err, "Error at: CreateAccount -> Error open data/usernames.json", false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//Unmarshal list of usernames
	var usernames []string
	err = json.Unmarshal(data, &usernames)
	if err != nil {
		utility.LogError(err, "Error at: CreateAccount -> Error unmarshal list of usernames from usernames.json", false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//Reading respond body
	data, err = io.ReadAll(r.Body)
	if err != nil {
		utility.LogError(err, "Error at: CreateAccount -> Error reading data from respond body", false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	//Create new user instance from respond body
	var newUser User
	err = json.Unmarshal(data, &newUser)
	if err != nil {
		utility.LogError(err, "Error at: CreateAccount -> Error unmarshal respond body data", false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//Check if user's username already exist in database or not
	for _, name := range usernames {
		if newUser.Username == name {
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write([]byte("Username already exist!"))
			if err != nil {
				utility.LogError(err, "Error at: CreateAccount -> Error sending message to client", false)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}
	}

	//Parsing list of usernames
	usernames = append(usernames, newUser.Username)
	jsonData, err := json.MarshalIndent(usernames, "", " ")
	if err != nil {
		utility.LogError(err, "Error at: CreateAccount -> Error marshal list of usernames", false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//Writing usernames to data.usernames.json
	err = utility.WriteFile("data/usernames.json", jsonData)
	if err != nil {
		utility.LogError(err, "Error at: CreateAccount -> Error wrting data to data/usernames.json", false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//Create new directory for new user, add user's information to user.json and transactions.json
	err = utility.CreateNewDir(newUser.Username)
	if err != nil {
		utility.LogError(err, "Error at: CreateAccount -> Error create new directory", false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	jsonData, err = json.MarshalIndent(newUser, "", " ")
	if err != nil {
		utility.LogError(err, "Error at: CreateAccount -> Error marshal user", false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	userPath := fmt.Sprintf("data/%s/user.json", newUser.Username)
	transactionsPath := fmt.Sprintf("data/%s/transactions.json", newUser.Username)

	err = utility.WriteFile(userPath, jsonData)
	if err != nil {
		utility.LogError(err, "Error at: CreateAccount -> Error writing to data "+userPath, false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = utility.WriteFile(transactionsPath, []byte("[]"))
	if err != nil {
		utility.LogError(err, "Error at: CreateAccount -> Error writing to data "+transactionsPath, false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//Send successful message to client
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte("Account created successfully!"))
	if err != nil {
		utility.LogError(err, "Error at: CreateAccount -> Error sending message to client", false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	//Reading data from respond body
	data, err := io.ReadAll(r.Body)
	if err != nil {
		utility.LogError(err, "Error at: Login -> Error reading data from respond body", false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	//Parsing user's login info
	info := make([]map[string]string, 0)
	err = json.Unmarshal(data, &info)
	if err != nil {
		utility.LogError(err, "Error at: Login -> Error unmarshal user's login info", false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//Check if username exist in database
	data, err = utility.OpenFile("data/usernames.json")
	if err != nil {
		utility.LogError(err, "Error at: Login -> Error reading data from data/usernames.json", false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var usernames []string
	err = json.Unmarshal(data, &usernames)
	if err != nil {
		utility.LogError(err, "Error at: Login -> Error unmarshal data from data/usernames.json", false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	for _, username := range usernames {
		if info[0]["username"] == username {
			//Reading data from user.json
			userPath := fmt.Sprintf("data/%s/user.json", username)
			data, err = utility.OpenFile(userPath)
			if err != nil {
				utility.LogError(err, "Error at: Login -> Error reading data from"+userPath, false)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			//Unmarshal user's data
			var user User
			err = json.Unmarshal(data, &user)
			if err != nil {
				utility.LogError(err, "Error at: Login -> Error unmarshal data from user.json", false)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			//Check if password is correct or not
			if user.Password == info[1]["password"] {
				w.WriteHeader(http.StatusAccepted)
				_, err = w.Write([]byte("Login successfully!"))
				if err != nil {
					utility.LogError(err, "Error at: Login -> Error sending message to client", false)
					http.Error(w, "Internal server error", http.StatusInternalServerError)
					return
				}
			} else {
				w.WriteHeader(http.StatusNotAcceptable)
				_, err = w.Write([]byte("Password not correct!"))
				if err != nil {
					utility.LogError(err, "Error at: Login -> Error sending message to client", false)
					http.Error(w, "Internal server error", http.StatusInternalServerError)
					return
				}
			}
			return
		}
	}

	//If username does not exist in database
	w.WriteHeader(http.StatusNotAcceptable)
	_, err = w.Write([]byte("Username incorrect! There is no such username"))
	if err != nil {
		utility.LogError(err, "Error at: Login -> Error sending message to client", false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	//Reading request body
	identity := r.Header.Get("identity")

	//Read data from usernames.json
	var usernames []string
	data, err := utility.OpenFile("data/usernames.json")
	if err != nil {
		utility.LogError(err, "Error at: GetProfile -> Error reading data from data/usernames.json", false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//Parsing usernames.json
	err = json.Unmarshal(data, &usernames)
	if err != nil {
		utility.LogError(err, "Error at: GetProfile -> Error unmarshal data from usernames.json", false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//Check for identity
	var isValid bool = false
	for _, name := range usernames {
		if name == string(identity) {
			isValid = true
		}
	}

	if !isValid {
		w.WriteHeader(http.StatusNotAcceptable)
		_, err = w.Write([]byte("Cannot verify the identity"))
		if err != nil {
			utility.LogError(err, "Error at: GetProfile -> Error sending message to client", false)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	//If identity is valid, get the data and send back to client
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	data, err = utility.OpenFile(fmt.Sprintf("data/%s/user.json", identity))
	if err != nil {
		utility.LogError(err, "Error at: GetProfile -> Error reading data from user.json", false)
		http.Error(w, "Internal server", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		utility.LogError(err, "Error at: GetProfile -> Error sending data to client", false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
