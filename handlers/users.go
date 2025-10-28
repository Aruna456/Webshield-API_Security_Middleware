package handlers

import (
	"fmt"
	"net/http"
)

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Users API!")
}
