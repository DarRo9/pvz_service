package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func RequestToken(t *testing.T, role string) string {
	data := map[string]string{"role": role}
	body, _ := json.Marshal(data)

	resp, err := http.Post("http://localhost:8080/dummyLogin", "application/json", bytes.NewReader(body))
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var token string
	_ = json.NewDecoder(resp.Body).Decode(&token)

	return token
}

func Post(t *testing.T, path string, token string, payload map[string]string) map[string]interface{} {
	var req *http.Request
	if payload == nil {
		req, _ = http.NewRequest("POST", "http://localhost:8080"+path, nil)
	} else {
		body, _ := json.Marshal(payload)
		req, _ = http.NewRequest("POST", "http://localhost:8080"+path, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)

	var result map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&result)
	return result
}

func Test_Integration(t *testing.T) {
	moderatorToken := RequestToken(t, "moderator")
	fmt.Println(moderatorToken)

	pvzResp := Post(t, "/pvz", moderatorToken, map[string]string{"city": "Санкт-Петербург"})
	pvzID := pvzResp["id"].(string)

	employeeToken := RequestToken(t, "employee")
	fmt.Println(employeeToken)

	recResp := Post(t, "/receptions", employeeToken, map[string]string{"pvzId": pvzID})
	assert.Equal(t, "in_progress", recResp["status"])

	for i := 0; i < 50; i++ {
		Post(t, "/products", employeeToken, map[string]string{
			"pvzId": pvzID,
			"type":  "обувь",
		})
	}

	Post(t, fmt.Sprintf("/pvz/%s/close_last_reception", pvzID), employeeToken, nil)
}
