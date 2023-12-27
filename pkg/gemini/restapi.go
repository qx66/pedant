package gemini

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/startopsz/rule/pkg/http"
	"os"
)

const (
	Api = "https://generativelanguage.googleapis.com/v1beta/models/"
)

type Contents struct {
	Contents []Content `json:"contents,omitempty"`
}

type Content struct {
	Role  string        `json:"role,omitempty"`
	Parts []interface{} `json:"parts,omitempty"`
}

type ContentText struct {
	Text string `json:"text,omitempty"`
}

type ContentInlineData struct {
	InlineData ContentImg `json:"inline_data,omitempty"`
}

type ContentImg struct {
	MimeType string `json:"mime_type,omitempty"`
	Data     string `json:"data,omitempty"` // image base64
}

type Response struct {
	Candidates     []Candidates   `json:"candidates,omitempty"`
	PromptFeedback PromptFeedback `json:"promptFeedback,omitempty"`
}

type Candidates struct {
	Content       CandidatesContent                `json:"content,omitempty"`
	FinishReason  string                           `json:"finishReason,omitempty"`
	Index         int                              `json:"index,omitempty"`
	SafetyRatings []CandidatesContentSafetyRatings `json:"safetyRatings,omitempty"`
}

type CandidatesContent struct {
	Parts []ContentText `json:"parts"`
	Role  string        `json:"role,omitempty"`
}

type CandidatesContentSafetyRatings struct {
	Category    string `json:"category,omitempty"`
	Probability string `json:"probability,omitempty"`
}

type PromptFeedback struct {
	SafetyRatings []CandidatesContentSafetyRatings `json:"safetyRatings,omitempty"`
}

type ApiKey string

func (apiKey ApiKey) Text(text string) (Response, error) {
	var response Response
	
	var parts []interface{}
	
	parts = append(parts, ContentText{
		Text: text,
	})
	
	contents := Contents{
		Contents: []Content{
			{
				Parts: parts,
			},
		},
	}
	
	contentsByte, err := json.Marshal(contents)
	if err != nil {
		return response, err
	}
	
	realUrl := fmt.Sprintf("%s%s:generateContent?key=%s", Api, "gemini-pro", apiKey)
	
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	
	req := http.Req{
		Method:  http.Post,
		Url:     realUrl,
		Body:    contentsByte,
		Headers: header,
		Timeout: 60,
	}
	
	resp, err := req.Do()
	if err != nil {
		return response, err
	}
	
	if resp.StatusCode != 200 {
		return response, errors.New(fmt.Sprintf("status: %d, body: %s", resp.StatusCode, string(resp.Body)))
	}
	
	//fmt.Println("resp.Body: ", string(resp.Body))
	err = json.Unmarshal(resp.Body, &response)
	return response, err
}

// 经过测试，imagePaths 建议只传一张图片，多张图片测试效果不是很好

func (apiKey ApiKey) TextAndImage(text string, imagePaths ...string) (Response, error) {
	var response Response
	
	//
	var parts []interface{}
	parts = append(parts, ContentText{
		Text: text,
	})
	
	//
	for _, imagePath := range imagePaths {
		//
		imgData, err := os.ReadFile(imagePath)
		if err != nil {
			return response, err
		}
		
		parts = append(parts, ContentInlineData{
			InlineData: ContentImg{
				MimeType: "image/jpeg",
				Data:     base64.StdEncoding.EncodeToString(imgData),
			},
		})
	}
	
	//
	contents := Contents{
		Contents: []Content{
			{
				Parts: parts,
			},
		},
	}
	
	contentsByte, err := json.Marshal(contents)
	if err != nil {
		return response, err
	}
	
	realUrl := fmt.Sprintf("%s%s:generateContent?key=%s", Api, "gemini-pro-vision", apiKey)
	
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	
	req := http.Req{
		Method:  http.Post,
		Url:     realUrl,
		Body:    contentsByte,
		Headers: header,
		Timeout: 60,
	}
	
	resp, err := req.Do()
	if err != nil {
		return response, err
	}
	
	if resp.StatusCode != 200 {
		return response, errors.New(fmt.Sprintf("status: %d, body: %s", resp.StatusCode, string(resp.Body)))
	}
	
	err = json.Unmarshal(resp.Body, &response)
	return response, err
}

func (apiKey ApiKey) Chat(history []Content, text string) (Response, error) {
	var response Response
	
	//
	var parts []interface{}
	parts = append(parts, ContentText{
		Text: text,
	})
	
	history = append(history, Content{
		Role:  "user",
		Parts: parts,
	})
	//
	contents := Contents{
		Contents: history,
	}
	
	contentsByte, err := json.Marshal(contents)
	if err != nil {
		return response, err
	}
	
	realUrl := fmt.Sprintf("%s%s:generateContent?key=%s", Api, "gemini-pro", apiKey)
	
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	
	req := http.Req{
		Method:  http.Post,
		Url:     realUrl,
		Body:    contentsByte,
		Headers: header,
		Timeout: 60,
	}
	
	resp, err := req.Do()
	if err != nil {
		return response, err
	}
	
	if resp.StatusCode != 200 {
		return response, errors.New(fmt.Sprintf("status: %d, body: %s", resp.StatusCode, string(resp.Body)))
	}
	
	err = json.Unmarshal(resp.Body, &response)
	return response, err
}
