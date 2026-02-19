// Package client provides an HTTP client for the Google Gemini image generation API.
package client

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "https://generativelanguage.googleapis.com/v1beta"

// Client is a Gemini API client for image generation.
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// New creates a new Gemini API client.
func New(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 180 * time.Second,
		},
	}
}

// GenerateRequest holds parameters for an image generation request.
type GenerateRequest struct {
	Model       string
	Prompt      string
	Negative    string
	AspectRatio string
	ImageSize   string
	Images      []ImageInput // reference images
}

// ImageInput represents an inline image to send with the request.
type ImageInput struct {
	MimeType string
	Data     []byte
}

// GenerateResult holds a single generated image.
type GenerateResult struct {
	ImageData []byte
	MimeType  string
	Index     int
}

// GenerateResponse holds the full API response.
type GenerateResponse struct {
	Results          []GenerateResult
	PromptTokens     int
	CandidatesTokens int
	TotalTokens      int
}

// APIError represents an error from the Gemini API.
type APIError struct {
	StatusCode int
	Code       int    `json:"code"`
	Message    string `json:"message"`
	Status     string `json:"status"`
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("gemini API error (%d): %s", e.Code, e.Message)
	}
	return fmt.Sprintf("gemini API error: HTTP %d", e.StatusCode)
}

// Generate sends an image generation request to the Gemini API.
func (c *Client) Generate(ctx context.Context, req GenerateRequest) (*GenerateResponse, error) {
	endpoint := fmt.Sprintf("%s/models/%s:generateContent?key=%s", baseURL, req.Model, c.apiKey)

	prompt := req.Prompt
	if req.Negative != "" {
		prompt = fmt.Sprintf("%s. Avoid: %s", prompt, req.Negative)
	}

	// Build parts
	parts := make([]map[string]interface{}, 0, len(req.Images)+1)
	for _, img := range req.Images {
		parts = append(parts, map[string]interface{}{
			"inlineData": map[string]interface{}{
				"mimeType": img.MimeType,
				"data":     base64.StdEncoding.EncodeToString(img.Data),
			},
		})
	}
	parts = append(parts, map[string]interface{}{"text": prompt})

	genConfig := map[string]interface{}{
		"responseModalities": []string{"IMAGE"},
	}
	imageConfig := map[string]interface{}{}
	if req.AspectRatio != "" {
		imageConfig["aspectRatio"] = req.AspectRatio
	}
	if req.ImageSize != "" {
		imageConfig["imageSize"] = req.ImageSize
	}
	if len(imageConfig) > 0 {
		genConfig["imageConfig"] = imageConfig
	}

	body := map[string]interface{}{
		"contents": []map[string]interface{}{
			{"parts": parts},
		},
		"generationConfig": genConfig,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, parseAPIError(resp.StatusCode, respBody)
	}

	return parseResponse(respBody)
}

func parseAPIError(statusCode int, body []byte) error {
	var errResp struct {
		Error struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Status  string `json:"status"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error.Message != "" {
		return &APIError{
			StatusCode: statusCode,
			Code:       errResp.Error.Code,
			Message:    errResp.Error.Message,
			Status:     errResp.Error.Status,
		}
	}
	return &APIError{StatusCode: statusCode, Code: statusCode, Message: string(body)}
}

func parseResponse(body []byte) (*GenerateResponse, error) {
	var apiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text       string `json:"text,omitempty"`
					InlineData *struct {
						MimeType string `json:"mimeType"`
						Data     string `json:"data"`
					} `json:"inlineData,omitempty"`
				} `json:"parts"`
			} `json:"content"`
			FinishReason string `json:"finishReason"`
		} `json:"candidates"`
		UsageMetadata *struct {
			PromptTokenCount     int `json:"promptTokenCount"`
			CandidatesTokenCount int `json:"candidatesTokenCount"`
			TotalTokenCount      int `json:"totalTokenCount"`
		} `json:"usageMetadata,omitempty"`
	}

	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	var results []GenerateResult
	idx := 0
	for _, candidate := range apiResp.Candidates {
		for _, part := range candidate.Content.Parts {
			if part.InlineData != nil && part.InlineData.Data != "" {
				imgData, err := base64.StdEncoding.DecodeString(part.InlineData.Data)
				if err != nil {
					continue
				}
				mimeType := part.InlineData.MimeType
				if mimeType == "" {
					mimeType = "image/png"
				}
				results = append(results, GenerateResult{
					ImageData: imgData,
					MimeType:  mimeType,
					Index:     idx,
				})
				idx++
			}
		}
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no images generated in response")
	}

	resp := &GenerateResponse{Results: results}
	if apiResp.UsageMetadata != nil {
		resp.PromptTokens = apiResp.UsageMetadata.PromptTokenCount
		resp.CandidatesTokens = apiResp.UsageMetadata.CandidatesTokenCount
		resp.TotalTokens = apiResp.UsageMetadata.TotalTokenCount
	}
	return resp, nil
}

// DetectMIMEType detects image MIME type from file header bytes.
func DetectMIMEType(data []byte) string {
	if len(data) < 12 {
		return "application/octet-stream"
	}
	if data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
		return "image/png"
	}
	if data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return "image/jpeg"
	}
	if data[0] == 0x52 && data[1] == 0x49 && data[2] == 0x46 && data[3] == 0x46 &&
		data[8] == 0x57 && data[9] == 0x45 && data[10] == 0x42 && data[11] == 0x50 {
		return "image/webp"
	}
	if data[0] == 0x47 && data[1] == 0x49 && data[2] == 0x46 {
		return "image/gif"
	}
	return "application/octet-stream"
}
