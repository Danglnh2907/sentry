package get

import (
	"encoding/json"
	_ "encoding/json"
	_ "fmt"
	"net/http"
	"sentry/dataStructure"
	"sentry/utility"
	_ "sentry/utility"
)

func HandleGetProfile(w http.ResponseWriter, r *http.Request) {

}

func HandleGetTransactions(w http.ResponseWriter, r *http.Request) {
	//Parse data to JSON
	data, err := json.MarshalIndent(dataStructure.Transactions, "", " ")
	if err != nil {
		utility.LogError(err, "Error parsing data to json", false)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func HandleGetReport(w http.ResponseWriter, r *http.Request) {

}
