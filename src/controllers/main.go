package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/cjlapao/common-go/helper/http_helper"
	"github.com/cjlapao/http-loadtester-go/entities"
	"github.com/cjlapao/http-loadtester-go/usecases"
)

func LoadController(w http.ResponseWriter, r *http.Request) {
	var body entities.LoadTest

	err := http_helper.MapRequestBody(r, &body)

	if err != nil {
		WriteError(w, err)
		return
	}

	results, err := usecases.ExecuteLoadTest(body)

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

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=\"api.json\"")

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

func StartLoadFileController(w http.ResponseWriter, r *http.Request) {
	var body entities.LoadTest

	err := http_helper.MapRequestBody(r, &body)

	if err != nil {
		WriteError(w, err)
		return
	}

	results, err := usecases.ExecuteLoadTest(body)

	if err != nil {
		WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+body.DisplayName+".md\"")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(results[0].MarkDown()))
}

func StartLoadMarkdownController(w http.ResponseWriter, r *http.Request) {
	var body entities.LoadTest

	err := http_helper.MapRequestBody(r, &body)

	if err != nil {
		WriteError(w, err)
		return
	}

	results, err := usecases.ExecuteLoadTest(body)

	if err != nil {
		WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "plain/text")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+body.DisplayName+".md\"")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(results[0].MarkDown()))
}
