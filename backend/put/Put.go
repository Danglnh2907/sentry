package put

import (
	_ "fmt"
	"net/http"
)

func Test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("PUT request"))
}
