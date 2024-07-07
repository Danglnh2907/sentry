package crud

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sentry/utility"
	"strconv"
)

type Transaction struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Date        string  `json:"date"`
	Cost        float64 `json:"cost"`
}

func PostTransaction(w http.ResponseWriter, r *http.Request) {
	//Get requester identity
	identity := r.Header.Get("Identity")

	//Reading data from request body
	jsonData, err := io.ReadAll(r.Body)
	if err != nil {
		utility.LogError(err, "Error at: PostTransaction -> Error reading data from request body", false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	/*Append new transaction to list of transactions*/
	var (
		transaction  Transaction
		transactions []Transaction
		filePath     string = fmt.Sprintf("data/%s/transactions.json", identity)
	)

	//Unmarshal request body
	err = json.Unmarshal(jsonData, &transaction)
	if err != nil {
		utility.LogError(err, "Error at: PostTransaction -> Error unmarshal request body", false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//Read data from transactions.json
	data, err := utility.OpenFile(filePath)
	if err != nil {
		utility.LogError(err, "Error at: PostTransaction -> Error reading transactions.json", false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//Unmarshal transactions.json
	err = json.Unmarshal(data, &transactions)
	if err != nil {
		utility.LogError(err, "Error at: PostTransaction -> Error unmarshal from transactions.json", false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	transactions = append(transactions, transaction)

	//Marshal new transactions
	jsonData, err = json.MarshalIndent(transactions, "", " ")
	if err != nil {
		utility.LogError(err, "Error at: Error marshal list of transactions", false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//Write data to transactions.json
	err = utility.WriteFile(filePath, jsonData)
	if err != nil {
		utility.LogError(err, "Error at: PostTransaction -> Error writing data to transactions.json", false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//Write successful message to user
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte("Adding transaction successfully!"))
	if err != nil {
		utility.LogError(err, "Error at: PostTransaction -> Error sending message to client", false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func PostTransactions(w http.ResponseWriter, r *http.Request) {
	//Create csv reader
	csvReader := csv.NewReader(r.Body)
	defer r.Body.Close()

	//Read the current transactions from database for appending
	transactionsPath := fmt.Sprintf("data/%s/transactions.json", r.Header.Get("Identity"))
	data, err := utility.OpenFile(transactionsPath)
	if err != nil {
		utility.LogError(err, "Error at: PostTransactions -> Error reading data from transactions.json", false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//Unmarshal data from transactions.json
	var transactions []Transaction
	err = json.Unmarshal(data, &transactions)
	if err != nil {
		utility.LogError(err, "Error at: PostTransactions -> Error unmarshal data from transactions.json", false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//Read all records in csv file
	records, err := csvReader.ReadAll()
	if err != nil {
		utility.LogError(err, "Error at: PostTransactions -> Error reading data from csv file", false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//Read each record (skip header line)
	var transaction Transaction
	for _, val := range records[1:] {
		//Check if csv file is correctly formatted
		if len(val) != 6 {
			utility.LogError(err, "csv file not formatted correctly", false)
			return
		}

		//Parse cost
		cost, err := strconv.ParseFloat(val[5], 64)
		if err != nil {
			utility.LogError(err, "Error parsing cost", false)
			return
		}
		transaction = Transaction{ID: val[0], Name: val[1], Description: val[2], Category: val[3], Date: val[4], Cost: cost}

		//Append new transaction to transactions
		transactions = append(transactions, transaction)
	}

	//Marshal list of transactions
	data, err = json.MarshalIndent(transactions, "", " ")
	if err != nil {
		utility.LogError(err, "Error at: PostTransactions -> Error marshal list of transactions", false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//Write new data to transactions.json
	err = utility.WriteFile(transactionsPath, data)
	if err != nil {
		utility.LogError(err, "Error at: PostTransactions -> Error writing data to transactions.json", false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//Send successful message to client
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte("Transactions added to database successfully!"))
	if err != nil {
		utility.LogError(err, "Error at: PostTransactions -> Error sending message to client", false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func GetTransactions(w http.ResponseWriter, r *http.Request) {
	identity := r.Header.Get("Identity")

	//Read data from transactions.json
	transactionPath := fmt.Sprintf("data/%s/transactions.json", identity)
	data, err := utility.OpenFile(transactionPath)
	if err != nil {
		utility.LogError(err, "Error at: GetTransactions -> Error reading data from transactions.json", false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//Unmarshal json data
	var transactions []Transaction
	err = json.Unmarshal(data, &transactions)
	if err != nil {
		utility.LogError(err, "Error at: GetTransactions -> Error unmarshal transactions.json", false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//Generate csv data
	records := make([][]string, 0)
	//Set header
	header := []string{"Name", "Descripiton", "Category", "Date", "Cost"}
	records = append(records, header)
	//Set each record
	for _, transaction := range transactions {
		//Set each record
		cost := strconv.FormatFloat(transaction.Cost, 'f', 2, 64)
		record := []string{transaction.Name, transaction.Description, transaction.Category, transaction.Date, cost}
		records = append(records, record)

	}

	//Writing csv data before sending
	buffer := new(bytes.Buffer)
	csvWriter := csv.NewWriter(buffer)
	err = csvWriter.WriteAll(records)
	if err != nil {
		utility.LogError(err, "Error at: GetTransactions -> Error writing csvData to buffer", false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		utility.LogError(err, "Error at: GetTransactions -> Error writing csv data after flush", false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	csvData := buffer.Bytes()

	//Sending data to client
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/csv")
	_, err = w.Write(csvData)
	if err != nil {
		utility.LogError(err, "Error at: GetTransactions -> Error sending data to client", false)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
