package oapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type APIRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Data    interface{}
	page    interface{}
}

func NewRequest(method, url string) *APIRequest {
	return &APIRequest{
		Method: method,
		URL:    url,
	}
}

func (apiReq *APIRequest) Do() (*APIResponse, error) {
	var request *http.Request
	apiResp := new(APIResponse)

	if apiReq.Data != nil {
		DataJSON, err := json.Marshal(apiReq.Data)
		if err != nil {
			return apiResp, fmt.Errorf("Failed marshal Data of %s %s request. Error: %v\n", apiReq.Method, apiReq.URL, err)
		}
		request, _ = http.NewRequest(apiReq.Method, apiReq.URL, bytes.NewBuffer(DataJSON))
	} else {
		request, _ = http.NewRequest(apiReq.Method, apiReq.URL, nil)
	}
	request.Close = true
	request.Header.Set("Content-Type", "application/json")
	for key, value := range apiReq.Headers {
		request.Header.Set(key, value)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return apiResp, fmt.Errorf("%s %s request failed with error: %v\n", apiReq.Method, apiReq.URL, err)
	}

	apiResp.Response = response

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		code, _ := strconv.Atoi(response.Header.Get("Error-Code"))
		apiResp.Code = code
		apiResp.ErrMessage = string(body)
		return apiResp, fmt.Errorf("%s %s request failed. Status: %v Err: %s", apiReq.Method, apiReq.URL, response.StatusCode, string(body))
	}

	if apiReq.page != nil {
		apiResp.Data = &apiReq.page
		if err := json.NewDecoder(response.Body).Decode(&apiResp.Data); err != nil {
			return apiResp, err
		}
	}

	return apiResp, nil
}
