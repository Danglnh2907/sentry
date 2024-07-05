package crud

import (
	//Import standard library
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
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
			utility.LogError(err, "Error reading user input")
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
		utility.LogError(err, "Error reading user input")
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
		utility.LogError(err, "Error reading user input")
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
			utility.LogError(err, "Error reading user input")
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
			utility.LogError(err, "Error reading user input")
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
		utility.LogError(err, "Error parsing data to json")
		return
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/transaction", bytes.NewBuffer(jsonData))
	if err != nil {
		utility.LogError(err, "Error making request")
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Identity", username)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utility.LogError(err, "Error sending request to server")
		return
	}
	defer resp.Body.Close()

	message, err := io.ReadAll(resp.Body)
	if err != nil {
		utility.LogError(err, "Error reading respond body")
		return
	}

	fmt.Println(string(message))
}

func AddTransByFile(filePath string) {

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
