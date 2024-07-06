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
func CreateNewDir(dirName string) error {
	//Check if directory already exist
	dirPath := fmt.Sprintf("data/%s", dirName)
	if _, err := os.Stat(dirPath); err != nil {
		err = os.Mkdir(dirPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	//Add user.json and transactions.json into user directory
	userPath := fmt.Sprintf("data/%s/user.json", dirName)
	if err := CreateNewFile(userPath); err != nil {
		return err
	}
	transactionsPath := fmt.Sprintf("data/%s/transactions.json", dirName)
	if err := CreateNewFile(transactionsPath); err != nil {
		return err
	}

	//Add empty array/object to json file
	if err := WriteFile(userPath, []byte("{}")); err != nil {
		return err
	}
	if err := WriteFile(transactionsPath, []byte("[]")); err != nil {
		return err
	}

	return nil
}

func CreateNewFile(filePath string) error {
	//If file does not exist, create new one
	if _, err := os.Stat(filePath); err != nil {
		file, err := os.Create(filePath)
		//Handle error
		if err != nil {
			return err
		}
		defer file.Close()
	}
	return nil
}

func OpenFile(filePath string) ([]byte, error) {
	if err := CreateNewFile(filePath); err != nil {
		return []byte(""), err
	}

	//Read data from file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return []byte(""), err
	}

	return data, nil
}

func WriteFile(filePath string, data []byte) error {
	if err := CreateNewFile(filePath); err != nil {
		return err
	}

	//Write data to file (overwrite mode)
	err := os.WriteFile(filePath, data, 0666)
	if err != nil {
		return err
	}

	return nil
}
