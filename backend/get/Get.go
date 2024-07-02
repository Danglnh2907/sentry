package get

import (
	"encoding/json"
	_ "fmt"
	"net/http"
	"sentry/dataModel"
	"sentry/utility"
)

func HandleGetProfile(w http.ResponseWriter, r *http.Request) {

}

func HandleGetUsername(w http.ResponseWriter, r *http.Request) {
	usernames := make([]string, 0)

	for _, val := range dataModel.Users {
		usernames = append(usernames, val.Username)
	}

	data, err := json.MarshalIndent(usernames, "", " ")
	if err != nil {
		utility.LogError(err, "Error parsing data to json", false)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
