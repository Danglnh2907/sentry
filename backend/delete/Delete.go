package delete

import (
	"net/http"
)

func HandleDeleteTransactions(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete transactions"))
}
