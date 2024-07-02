package main

import (
	//Import standard library
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	//Import standard package
	"sentry/dataModel"
	"sentry/get"
	"sentry/post"
	"sentry/utility"
)

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

func prepareData() {
	utility.CreateNewFile("data/users.json")

	err := json.Unmarshal(utility.OpenFile("data/users.json"), &dataModel.Users)
	if err != nil {
		utility.LogError(err, "Error loading data from database", true)
	}
}

func main() {
	//Prepare data for using
	prepareData()

	//Create multiplexer for routing
	mux := http.NewServeMux()

	//handle GET method here
	mux.HandleFunc("/get-username", get.HandleGetUsername)

	//handle POST request
	mux.HandleFunc("/user", post.HandlePostUser)

	//Run server at port 8080
	fmt.Println("Sentry running at http://localhost:8080")
	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		var w http.ResponseWriter
		http.Error(w, "Server error", http.StatusInternalServerError)
		log.Fatal(err)
	}
}
