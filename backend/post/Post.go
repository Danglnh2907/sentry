package post

import (
	//Import standard library
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	//Import user's defined package
	"sentry/dataModel"
	"sentry/utility"
)

func HandlePostUser(w http.ResponseWriter, r *http.Request) {
	//Parse request body into slice of bytes
	data, err := io.ReadAll(r.Body)
	if err != nil {
		utility.LogError(err, "Error parsing request body", false)
		return
	}

	//Parse json data to Users
	var newUser dataModel.User
	json.Unmarshal(data, &newUser)
	dataModel.Users = append(dataModel.Users, newUser)

	//Write new data to json file
	jsonData, err := json.MarshalIndent(dataModel.Users, "", " ")
	if err != nil {
		utility.LogError(err, "Error parsing data to json", false)
		return
	}

	filePath := fmt.Sprintf("data/%s.json", newUser.Username)
	utility.CreateNewFile(filePath)
	utility.WriteFile(filePath, jsonData)

	//Send succesful message to client
	w.Write([]byte("Post user successfully!"))
}

func HandlePostTransactions(w http.ResponseWriter, r *http.Request) {
	//Parse request body to slice of bytes
	data, err := io.ReadAll(r.Body)
	if err != nil {
		utility.LogError(err, "Error parsing request body", false)
		return
	}

	//Parse json data to Transactions
	var newTransaction dataModel.Transaction
	json.Unmarshal(data, &newTransaction)
	//dataModel.Transactions = append(dataModel.Transactions, newTransaction)

	//Send succesful message to client
	w.Write([]byte("Post transaction successfully!"))
}
