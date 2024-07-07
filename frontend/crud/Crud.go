package crud

import (
	//Import standard library
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	//Import user's defined package
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

func AddTrans(username string) {
	/*Add for transaction's information*/
	reader := bufio.NewReader(os.Stdin)
	var (
		name, description, category, date string
		cost                              float64
		transaction                       Transaction
		err                               error
		isValid                           bool
	)

	//Ask for transaction's name
	isValid = false
	for !isValid {
		fmt.Print("Enter transaction's name: ")
		name, err = reader.ReadString('\n')
		if err != nil {
			utility.LogError(err, "Error at: AddTrans -> Error reading transaction's name")
			return
		}
		name = strings.TrimSpace(name)

		isValid = len(name) > 0
		if !isValid {
			fmt.Println("Transaction's name must not be empty")
			//continue
		}

		transaction.Name = name
	}

	//Ask for transaction's description
	fmt.Print("Enter transaction's description (optional): ")
	description, err = reader.ReadString('\n')
	if err != nil {
		utility.LogError(err, "Error at: AddTrans -> Error reading transaction's description")
		return
	}
	description = strings.TrimSpace(description)

	if len(description) == 0 {
		transaction.Description = "none"
	} else {
		transaction.Description = description
	}

	//Ask for transaction's category
	fmt.Print("Enter transaction's category (optional): ")
	category, err = reader.ReadString('\n')
	if err != nil {
		utility.LogError(err, "Error at: AddTrans -> Error reading transaction's category")
		return
	}
	category = strings.TrimSpace(category)

	if len(category) == 0 {
		transaction.Category = "others"
	} else {
		transaction.Category = category
	}

	//Ask for transaction's date
	isValid = false
	for !isValid {
		fmt.Print("Enter date of transaction (dd/mm/yyyy): ")
		date, err = reader.ReadString('\n')
		if err != nil {
			utility.LogError(err, "Error at: AddTrans -> Error reading transaction's date")
			return
		}

		date = strings.TrimSpace(date)

		dateValue, err := time.Parse("02/01/2006", date)
		if err != nil {
			fmt.Println("Date invalid!")
			isValid = false
			continue
		}

		transaction.Date = dateValue.Format("02/01/2006")
		isValid = true
	}

	//Ask for transaction's cost
	isValid = false
	for !isValid {
		fmt.Print("Enter the cost of transaction: ")
		_, err = fmt.Scanf("%f", &cost)
		if err != nil {
			utility.LogError(err, "Error at: AddTrans -> Error reading transaction's cost")
			return
		}

		isValid = cost > 0
		if !isValid {
			fmt.Println("Transaction's cost cannot be negative")
			continue
		}

		transaction.Cost = cost
	}

	//Generate transaction's ID
	transaction.ID = generateID(transaction, username)

	//Make a http request to server
	jsonData, err := json.MarshalIndent(transaction, "", " ")
	if err != nil {
		utility.LogError(err, "Error at: AddTrans -> Error marshal transaction")
		return
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/transaction", bytes.NewBuffer(jsonData))
	if err != nil {
		utility.LogError(err, "Error at: AddTrans -> Error making new request")
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Identity", username)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utility.LogError(err, "Error at: AddTrans -> Error sending request to server")
		return
	}
	defer resp.Body.Close()

	message, err := io.ReadAll(resp.Body)
	if err != nil {
		utility.LogError(err, "Error at: AddTrans -> Error reading respond body")
		return
	}

	fmt.Println(string(message))
}

func AddTransByFile(filePath, username string) {
	//Check if file exist
	_, err := os.Stat(filePath)
	if err != nil {
		utility.LogError(err, "Error at: AddTransByFile -> File does not exist")
		return
	}

	//If file exist, add ID for each transaction
	csvFile, err := os.Open(filePath)
	if err != nil {
		utility.LogError(err, "Error at: AddTransByFile -> Error open csv file")
		return
	}
	csvReader := csv.NewReader(csvFile)

	records, err := csvReader.ReadAll()
	if err != nil {
		utility.LogError(err, "Error at: AddTransByFile -> Error reading data from csv file")
		return
	}

	newRecords := make([][]string, len(records))
	for index, record := range records {
		tempRecord := record
		if index == 0 {
			record = make([]string, 0)
			record = append(record, "ID")
			record = append(record, tempRecord...)
		} else {
			cost, err := strconv.ParseFloat(record[4], 64)
			if err != nil {
				utility.LogError(err, "Error at: AddTransByFile -> Error parsing number from csv file")
				return
			}
			transaction := Transaction{Name: record[0], Description: record[1], Category: record[2], Date: record[3], Cost: cost}
			record = make([]string, 0)
			record = append(record, generateID(transaction, username))
			record = append(record, tempRecord...)
		}
		newRecords[index] = record
	}

	//Write newRecords to []byte using buffer
	buffer := new(bytes.Buffer)
	writer := csv.NewWriter(buffer)
	err = writer.WriteAll(newRecords)
	if err != nil {
		utility.LogError(err, "Error at: AddTransByFile -> Error writing csv data to buffer")
		return
	}
	writer.Flush()
	if err = writer.Error(); err != nil {
		utility.LogError(err, "Error at: AddTransByFile -> Error writing csv data after flush")
		return
	}
	csvData := buffer.Bytes()

	//Send data to server
	url := "http://localhost:8080/transactions"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(csvData))
	if err != nil {
		utility.LogError(err, "Error at: AddTransByFileError making new request")
		return
	}

	req.Header.Set("Identity", username)
	req.Header.Set("Content-Type", "text/csv")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utility.LogError(err, "Error at: AddTransByFile -> Error sending request to server")
		return
	}
	defer resp.Body.Close()

	//Print message to user
	message, err := io.ReadAll(resp.Body)
	if err != nil {
		utility.LogError(err, "Error at: AddTransByFile -> Error reading respond body")
		return
	}

	fmt.Println(string(message))
}

func generateID(transaction Transaction, username string) string {
	//Get current time in nanosecond
	timemstamp := time.Now().UnixNano()

	//Get random integer
	random := rand.Int63()

	//Get all transaction's information
	info := fmt.Sprintf("%s%f", transaction.Name, transaction.Cost)

	//Combine all data
	data := fmt.Sprintf("%s%d%d%s", info, random, timemstamp, username)

	//Hash data to get ID
	hash := sha256.Sum256([]byte(data))

	//Return the hash to hexa string
	return hex.EncodeToString(hash[:])
}

func GetTransactions(username string) {
	//Make request to server
	url := "http://localhost:8080/get-transactions"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		utility.LogError(err, "Error at: GetTransactions -> Error making new request")
		return
	}

	req.Header.Set("Identity", username)
	req.Header.Set("Content-Type", "text/csv")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utility.LogError(err, "Error at: GetTransactions -> Error sending request to server")
		return
	}
	defer resp.Body.Close()

	//Reading data from resp body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		utility.LogError(err, "Error at: GetTransactions -> Error reading data from respond body")
		return
	}

	//Writting data to csv file
	csvPath := "./transactions.csv"
	err = utility.WriteFile(csvPath, data)
	if err != nil {
		utility.LogError(err, "Error at: GetTransactions -> Error writing data to csv file")
		return
	}

	//Send successful message to user
	fmt.Printf("%s\n", fmt.Sprintf("Transactions data has been store at %s", csvPath))
}
