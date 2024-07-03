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
	IsLogIn      bool     `json:"isLogin"`
}

func HandleCreateAccount(w http.ResponseWriter, r *http.Request) {
	//Get the list of usernames in database
	var usernames []string
	err := json.Unmarshal(utility.OpenFile("data/usernames.json"), &usernames)
	if err != nil {
		utility.LogError(err, "Error reading data from usernames.json", false)
	}

	//Parsing respond body
	data, err := io.ReadAll(r.Body)
	if err != nil {
		utility.LogError(err, "Error reading data from respond body", false)
	}
	defer r.Body.Close()

	//Create new user instance
	var newUser User
	err = json.Unmarshal(data, &newUser)
	if err != nil {
		utility.LogError(err, "Error parsing json data", false)
		return
	}
	newUser.IsLogIn = true

	//Check if user's username already exist in database or not
	for _, name := range usernames {
		if newUser.Username == name {
			_, err = w.Write([]byte("Username already exist!"))
			if err != nil {
				utility.LogError(err, "Error sending message to client", false)
			}
			return
		}
	}

	//Add new username to usernames.json
	usernames = append(usernames, newUser.Username)
	jsonData, err := json.MarshalIndent(usernames, "", " ")
	if err != nil {
		utility.LogError(err, "Error parsing data to json", false)
	}
	utility.WriteFile("data/usernames.json", jsonData)

	//Create new dir for new user, add user's information to user.json
	utility.CreateNewDir(newUser.Username)
	jsonData, err = json.MarshalIndent(newUser, "", " ")
	if err != nil {
		utility.LogError(err, "Error parsing data to json", false)
	}
	utility.WriteFile(fmt.Sprintf("data/%s/user.json", newUser.Username), jsonData)

	//Send successful message to client
	_, err = w.Write([]byte("Account created successfully!"))
	if err != nil {
		utility.LogError(err, "Error sending message to client", false)
	}
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	//Reading data from respond body
	data, err := io.ReadAll(r.Body)
	if err != nil {
		utility.LogError(err, "Error reading data from respond body", false)
		return
	}
	defer r.Body.Close()

	//Parsing data
	info := make([]map[string]string, 0)
	err = json.Unmarshal(data, &info)
	if err != nil {
		utility.LogError(err, "Error parsing data", false)
		return
	}

	//Check if username exist in database
	data = utility.OpenFile("data/usernames.json")
	var usernames []string
	err = json.Unmarshal(data, &usernames)
	if err != nil {
		utility.LogError(err, "Error reading data from usernames.json", false)
		return
	}

	for _, name := range usernames {
		if info[0]["username"] == name {
			//Get password from user
			var user User
			data = utility.OpenFile(fmt.Sprintf("data/%s/user.json", name))
			err = json.Unmarshal(data, &user)
			if err != nil {
				utility.LogError(err, "Error reading data from user.json", false)
				return
			}

			//Check if password is correct or not
			if user.Password == info[1]["password"] {
				user.IsLogIn = true
				w.WriteHeader(http.StatusOK)
				_, err = w.Write([]byte("Login successfully!"))
				if err != nil {
					utility.LogError(err, "Error sending message to client", false)
					return
				}
			} else {
				w.WriteHeader(http.StatusNotAcceptable)
				_, err = w.Write([]byte("Password not correct!"))
				if err != nil {
					utility.LogError(err, "Error sending message to client", false)
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
		utility.LogError(err, "Error sending message to client", false)
		return
	}
}

func GetProfile(w http.ResponseWriter, r *http.Request) {

}
