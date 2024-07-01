package utility

import (
	//Import standard library

	_ "encoding/json"
	"fmt"
	"log"
	"os"

	//Import user's defined package
	_ "sentry/dataStructure"
)

/*Error handling*/
func LogError(err error, message string, isCritical bool) {
	log.Println(message)
	log.Printf("Error at %e\n", err)

	//If error is a critical one, end the server
	if isCritical {
		os.Exit(1)
	}
}

/*File handling*/
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

func WriteFile(filePath string, data []byte) bool {
	CreateNewFile(filePath)

	//Write data to file (overwrite mode)
	err := os.WriteFile(filePath, data, 0666)
	if err != nil {
		LogError(err, fmt.Sprintf("Error writing data to %s file\n", filePath), false)
		return false
	}

	return true
}
