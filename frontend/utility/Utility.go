package utility

import (
	"fmt"
	"log"
)

func LogError(err error, message string) {
	fmt.Println("Fatal Error! ", message)
	log.Fatal(err)
}
