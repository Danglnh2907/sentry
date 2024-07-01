package utility

import (
	_ "fmt"
	"log"
	"net/http"
	"os"
)

func HandleInternalError(err error) {
	if err != nil {
		log.Printf("Error at: %s", err)
	}
}

func SendServerError(err error, w http.ResponseWriter, statusCode int) {
	HandleInternalError(err)
	if err != nil {
		http.Error(w, "Server error", statusCode)
		os.Exit(1)
	}

}
