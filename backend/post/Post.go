package post

import (
	"net/http"
)

func HandlePostTransactions(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Post Transactions"))
}
