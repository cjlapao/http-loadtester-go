package controller

import (
	"encoding/json"
	"net/http"
)

func TestController(w http.ResponseWriter, r *http.Request) {
	response := "testing new controller"

	json.NewEncoder(w).Encode(response)
}
