package main

import (
	//Import standard library
	"fmt"
	"net/http"

	//Import standard package
	"sentry/delete"
	"sentry/get"
	"sentry/post"
	"sentry/put"
	"sentry/utility"
)

type Transaction struct {
	Name        string  `json:"name"`
	Descripiton string  `json:"description"`
	Category    string  `json:"category"`
	Date        string  `json:"date"`
	Cost        float64 `json:"cost"`
}

/*
var transactions []Transaction = make([]Transaction, 0)

	func handleGet(w http.ResponseWriter , r *http.Request) {
		//Data for testing
		//transactions = append(transactions, Transaction{Name: "Test", Descripiton: "This is some text", Category: "Test", Date: "29/06/2024", Cost: 0.0})

		//Marshal data into json
		jsonData, err := json.MarshalIndent(transactions, "", " ")

		//Handle error
		utility.SendServerError(err, w, http.StatusInternalServerError)

		//Set header to application/json for sending json data
		w.Header().Set("Content-Type", "application/json")

		//Write json data to respone writer
		_, err = w.Write(jsonData)

		//Print error
		if err != nil {
			fmt.Println("Error sending data to client")
		}
	}

	func handleGetFile(w http.ResponseWriter , r *http.Request) {
		//Open the file for sending
		imgFilePath := "main.html"
		img, err := os.Open(imgFilePath)
		utility.HandleInternalError(err)
		defer img.Close()

		//Set content type to the appropriate file exttension
		w.Header().Set("Content-Type", "text/html")
		//The content disposition set whether the file is an attachment (can be downloaded) or display inline in browser
		w.Header().Set("Content-Disposition", "inline")

		//Copy the file content to respond writer
		_, err = io.Copy(w, img)
		utility.SendServerError(err, w, http.StatusInternalServerError)
	}
*/

func main() {
	//Create multiplexer for routing
	mux := http.NewServeMux()

	//handle GET method here
	mux.HandleFunc("/get-profile", get.HandleGetProfile)
	mux.HandleFunc("/get-transactions", get.HandleGetTransactions)
	mux.HandleFunc("/get-report", get.HandleGetReport)

	//handle POST method
	mux.HandleFunc("/transactions", post.HandlePostTransactions)

	//handle PUT method
	mux.HandleFunc("/put", put.Test)

	//handle DELETE method
	mux.HandleFunc("/delete-transactions", delete.HandleDeleteTransactions)

	//Run server at port 8080
	fmt.Println("Sentry running at http://localhost:8080")

	err := http.ListenAndServe("localhost:8080", mux)

	//Sending error to terminal
	utility.HandleInternalError(err)

}
