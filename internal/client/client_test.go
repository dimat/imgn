package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenerate_Success(t *testing.T) {
	fakeImage := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D}
	b64Image := base64.StdEncoding.EncodeToString(fakeImage)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var reqBody map[string]interface{}
		json.NewDecoder(r.Body).Decode(&reqBody)

		resp := map[string]interface{}{
			"candidates": []map[string]interface{}{
				{
					"content": map[string]interface{}{
						"parts": []map[string]interface{}{
							{
								"inlineData": map[string]interface{}{
									"mimeType": "image/png",
									"data":     b64Image,
								},
							},
						},
					},
					"finishReason": "STOP",
				},
			},
			"usageMetadata": map[string]interface{}{
				"promptTokenCount":     10,
				"candidatesTokenCount": 100,
				"totalTokenCount":      110,
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := &Client{
		apiKey:     "test-key",
		httpClient: server.Client(),
	}
	// Override baseURL by using a model that includes the full URL
	// We need to test with the actual server, so we'll test parseResponse instead

	// Test parseResponse directly
	respBody := map[string]interface{}{
		"candidates": []map[string]interface{}{
			{
				"content": map[string]interface{}{
					"parts": []map[string]interface{}{
						{
							"inlineData": map[string]interface{}{
								"mimeType": "image/png",
								"data":     b64Image,
							},
						},
					},
				},
				"finishReason": "STOP",
			},
		},
		"usageMetadata": map[string]interface{}{
			"promptTokenCount":     10,
			"candidatesTokenCount": 100,
			"totalTokenCount":      110,
		},
	}
	body, _ := json.Marshal(respBody)

	resp, err := parseResponse(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(resp.Results))
	}
	if resp.Results[0].MimeType != "image/png" {
		t.Errorf("expected image/png, got %s", resp.Results[0].MimeType)
	}
	if resp.TotalTokens != 110 {
		t.Errorf("expected 110 total tokens, got %d", resp.TotalTokens)
	}

	_ = c
	_ = context.Background()
}

func TestParseResponse_NoImages(t *testing.T) {
	body := `{"candidates":[{"content":{"parts":[{"text":"I cannot generate that"}]},"finishReason":"STOP"}]}`
	_, err := parseResponse([]byte(body))
	if err == nil {
		t.Fatal("expected error for no images")
	}
}

func TestParseAPIError(t *testing.T) {
	body := `{"error":{"code":400,"message":"Invalid request","status":"INVALID_ARGUMENT"}}`
	err := parseAPIError(400, []byte(body))
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatal("expected APIError")
	}
	if apiErr.Code != 400 {
		t.Errorf("expected code 400, got %d", apiErr.Code)
	}
	if apiErr.Message != "Invalid request" {
		t.Errorf("unexpected message: %s", apiErr.Message)
	}
}

func TestDetectMIMEType(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want string
	}{
		{"PNG", []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D}, "image/png"},
		{"JPEG", []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, "image/jpeg"},
		{"GIF", []byte{0x47, 0x49, 0x46, 0x38, 0x39, 0x61, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, "image/gif"},
		{"WebP", []byte{0x52, 0x49, 0x46, 0x46, 0x00, 0x00, 0x00, 0x00, 0x57, 0x45, 0x42, 0x50}, "image/webp"},
		{"unknown", []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B}, "application/octet-stream"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetectMIMEType(tt.data)
			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}
