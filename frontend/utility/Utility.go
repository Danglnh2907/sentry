package utility

import (
	"fmt"
	"log"
	"os"
)

func LogError(err error, message string) {
	fmt.Println("Fatal Error! ", message)
	log.Fatal(err)
}

/*File Handling*/
func CreateEnvFile(filePath string) error {
	//If file does not exist, create new file
	if _, err := os.Stat(filePath); err != nil {
		//Create new file
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		//Set initial value
		_, err = file.WriteString("state=false\n")
		if err != nil {
			return err
		}
		_, err = file.WriteString("user=\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateNewDir(filePath string) error {
	//If dir does not exist, create new dir
	if _, err := os.Stat(filePath); err != nil {
		//Create new dir
		err := os.Mkdir(filePath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateNewFile(filePath string) error {
	//If file does not exist, create new file
	if _, err := os.Stat(filePath); err != nil {
		//Create new file
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	return nil
}

func OpenFile(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func WriteFile(filepath string, data []byte) error {
	//Create new file (if not exist)
	err := CreateNewFile(filepath)
	if err != nil {
		return err
	}

	//Write data to file
	err = os.WriteFile(filepath, data, 0666)
	if err != nil {
		return err
	}
	return nil
}

func SetEnvVar(state string, username string) error {
	err := os.Setenv("state", state)
	if err != nil {
		return err
	}
	err = os.Setenv("user", username)
	if err != nil {
		return err
	}
	return nil
}
