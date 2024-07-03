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
func CreateNewFile(filePath string) {
	//If file does not exist, create new file
	if _, err := os.Stat(filePath); err != nil {
		//Create new file
		file, err := os.Create(filePath)
		if err != nil {
			LogError(err, fmt.Sprintf("Error create %s file", filePath))
		}
		defer file.Close()

		//Set initial value
		file.WriteString("state=false\n")
		file.WriteString("user=\n")
	}
}

func SetEnvVar(state string, username string) {
	os.Setenv("state", state)
	os.Setenv("user", username)
}
