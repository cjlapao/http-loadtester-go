package controller

import (
	"encoding/json"
	"net/http"

	"github.com/cjlapao/common-go/helper/http_helper"
	"github.com/cjlapao/http-loadtester-go/jobs"
)

func LoadController(w http.ResponseWriter, r *http.Request) {
	var body jobs.LoadTest

	err := http_helper.MapRequestBody(r, &body)

	if err != nil {
		WriteError(w, err)
		return
	}

	results, err := jobs.ExecuteLoadTest(body)

	if err != nil {
		WriteError(w, err)
		return
	}

	response := make([]LoadTestResponse, 0)
	for _, result := range results {
		responseItem := LoadTestResponse{
			ID:            result.ID,
			Name:          result.Name,
			Type:          result.Type,
			OperationType: result.OperationType,
			Target:        result.Target,
			Options:       result.Options,
			Result:        result.Result,
		}
		response = append(response, responseItem)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
