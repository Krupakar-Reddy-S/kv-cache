package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

type PutRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Key     string `json:"key,omitempty"`
	Value   string `json:"value,omitempty"`
}

func TestBasicOperations(t *testing.T) {
	baseURL := "http://localhost:7171"
	maxRetries := 5
	var err error

	t.Log("Checking if service is running...")
	for i := 0; i < maxRetries; i++ {
		_, err = http.Get(baseURL + "/get?key=test")
		if err == nil {
			break
		}
		t.Logf("Service not ready, retrying in 2 seconds (attempt %d/%d)", i+1, maxRetries)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		t.Fatalf("Service is not running on %s", baseURL)
	}

	t.Log("Service is running! Starting tests...")

	// Test 1: Put operation
	t.Run("Put Operation", func(t *testing.T) {
		putReq := PutRequest{
			Key:   "test-key",
			Value: "test-value",
		}
		jsonData, _ := json.Marshal(putReq)

		resp, err := http.Post(baseURL+"/put", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Put request failed: %v", err)
		}

		var putResp Response
		json.NewDecoder(resp.Body).Decode(&putResp)
		resp.Body.Close()

		t.Logf("Put test result: %s - %s", putResp.Status, putResp.Message)
	})

	// Test 2: Get operation
	t.Run("Get Operation", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/get?key=test-key")
		if err != nil {
			t.Fatalf("Get request failed: %v", err)
		}

		var getResp Response
		json.NewDecoder(resp.Body).Decode(&getResp)
		resp.Body.Close()

		t.Logf("Get test result: %s - Value: %s", getResp.Status, getResp.Value)
	})

	// Test 3: Get non-existent key
	t.Run("Get Non-existent Key", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/get?key=non-existent")
		if err != nil {
			t.Fatalf("Get request failed: %v", err)
		}

		var notFoundResp Response
		json.NewDecoder(resp.Body).Decode(&notFoundResp)
		resp.Body.Close()

		t.Logf("Get non-existent key test result: %s - %s", notFoundResp.Status, notFoundResp.Message)
	})
}