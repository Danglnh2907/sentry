package post

import (
	//Import standard library
	"encoding/json"
	"io"
	"net/http"

	//Import user's defined package
	"sentry/dataStructure"
	"sentry/utility"
)

func HandlePostTransactions(w http.ResponseWriter, r *http.Request) {
	//Parse request body to slice of bytes
	data, err := io.ReadAll(r.Body)
	if err != nil {
		utility.LogError(err, "Error parsing request body", false)
		return
	}

	//Parse json data to Transactions
	var newTransaction dataStructure.Transaction
	json.Unmarshal(data, &newTransaction)
	dataStructure.Transactions = append(dataStructure.Transactions, newTransaction)

	//Send succesful message to client
	w.Write([]byte("Post transaction successfully!"))
}
