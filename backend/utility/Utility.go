package utility

import (
	//Import standard library
	"fmt"
	"log"
	"os"
	//Import user's defined package
)

/*Error handling*/
func LogError(err error, message string, isCritical bool) {
	log.Println(message)
	log.Printf("Error: %s\n", err)

	//If error is a critical one, end the server
	if isCritical {
		os.Exit(1)
	}
}

/*File handling*/
func CreateNewDir(dirName string) {
	//Check if directory already exist
	dirPath := fmt.Sprintf("data/%s", dirName)
	if _, err := os.Stat(dirPath); err != nil {
		err = os.Mkdir(dirPath, os.ModePerm)
		if err != nil {
			LogError(err, "Error creating user's directory", false)
			return
		}
	}
	//Add user.json and transactions.json into user directory
	userPath := fmt.Sprintf("data/%s/user.json", dirName)
	CreateNewFile(userPath)
	transactionsPath := fmt.Sprintf("data/%s/transactions.json", dirName)
	CreateNewFile(transactionsPath)

	//Add empty array/object to json file
	WriteFile(userPath, []byte("{}"))
	WriteFile(transactionsPath, []byte("[]"))
}

func CreateNewFile(filePath string) {
	//If file does not exist, create new one
	if _, err := os.Stat(filePath); err != nil {
		_, err = os.Create(filePath)
		//Handle error
		if err != nil {
			LogError(err, fmt.Sprintf("Error create %s file\n", filePath), false)
		}
	}
}

func OpenFile(filePath string) []byte {
	CreateNewFile(filePath)

	//Read data from file
	data, err := os.ReadFile(filePath)
	if err != nil {
		LogError(err, fmt.Sprintf("Error reading data from %s\n", filePath), false)
		return nil
	}

	return data
}

func WriteFile(filePath string, data []byte) {
	CreateNewFile(filePath)

	//Write data to file (overwrite mode)
	err := os.WriteFile(filePath, data, 0666)
	if err != nil {
		LogError(err, fmt.Sprintf("Error writing data to %s file\n", filePath), false)
	}
}
