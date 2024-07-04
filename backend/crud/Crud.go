package crud

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sentry/utility"
)

type Transaction struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Date        string  `json:"date"`
	Cost        float64 `json:"cost"`
}

func HandlePostTransaction(w http.ResponseWriter, r *http.Request) {
	//Get requester identity
	identity := r.Header.Get("Identity")

	//Reading data from request body
	jsonData, err := io.ReadAll(r.Body)
	if err != nil {
		utility.LogError(err, "Error reading data from request body", false)
		return
	}
	defer r.Body.Close()

	//Append new transaction to list of transactions
	var (
		transaction  Transaction
		transactions []Transaction
		filePath     string = fmt.Sprintf("data/%s/transactions.json", identity)
	)

	err = json.Unmarshal(jsonData, &transaction)
	if err != nil {
		utility.LogError(err, "Error parsing json data", false)
		return
	}

	err = json.Unmarshal(utility.OpenFile(filePath), &transactions)
	if err != nil {
		utility.LogError(err, "Error reading data from transactions.json", false)
		return
	}

	transactions = append(transactions, transaction)
	jsonData, err = json.MarshalIndent(transactions, "", " ")
	if err != nil {
		utility.LogError(err, "Error parsing data to json", false)
		return
	}

	utility.WriteFile(filePath, jsonData)

	//Write message to user
	w.Write([]byte("Adding transaction successfully!"))
}
