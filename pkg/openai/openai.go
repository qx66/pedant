package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

/*
text-davinci-003
*/

type TextDavinci003 struct {
	Model       string  `json:"model,omitempty"`
	Prompt      string  `json:"prompt,omitempty"`      // 问题
	Temperature float32 `json:"temperature,omitempty"` // default 1, 可以调成 0.8试试
	MaxTokens   int     `json:"max_tokens,omitempty"`  // 一般最大支持 4096
	//TopP int `json:"top_p,omitempty"` // default: 1
}

/*
"gpt-3.5-turbo-0301"
*/

type GptTurbo0301 struct {
	Model    string                `json:"model,omitempty"`
	Messages []GptTurbo0301Message `json:"messages,omitempty"`
}

type GptTurbo0301Message struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

const (
	chatApi                  = "https://api.openai.com/v1/chat/completions"
	chatApiContentType       = "application/json"
	completionApi            = "https://api.openai.com/v1/completions"
	completionApiContentType = "application/json"
)

const (
	ChatRoleSystem    = "system"
	ChatRoleUser      = "user"
	ChatRoleAssistant = "assistant"
	
	ChatModuleGpt4       = "gpt-4"
	ChatModuleGpt35Turbo = "gpt-3.5-turbo-0301"
	ChatModuleGpt432K    = "gpt-4-32k"
)

type ChatModuleResponse struct {
	Id      string                      `json:"id,omitempty"`
	Object  string                      `json:"object,omitempty"`
	Created int64                       `json:"created,omitempty"`
	Model   string                      `json:"model,omitempty"`
	Usage   ChatModuleResponseUsage     `json:"usage,omitempty"`
	Choices []ChatModuleResponseChoices `json:"choices,omitempty"`
}

type ChatModuleResponseUsage struct {
	PromptTokens     int64 `json:"prompt_tokens,omitempty"`
	CompletionTokens int64 `json:"completion_tokens,omitempty"`
	TotalTokens      int64 `json:"total_tokens,omitempty"`
}

type ChatModuleResponseChoices struct {
	Message      ChatModuleResponseChoicesMessage `json:"message,omitempty"`
	FinishReason string                           `json:"finish_reason,omitempty"`
	Index        int                              `json:"index,omitempty"`
}

type ChatModuleResponseChoicesMessage struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

type CompletionModuleResponse struct {
	Id      string                            `json:"id,omitempty"`
	Object  string                            `json:"object,omitempty"`
	Created int64                             `json:"created,omitempty"`
	Model   string                            `json:"model,omitempty"`
	Usage   CompletionModuleResponseUsage     `json:"usage,omitempty"`
	Choices []CompletionModuleResponseChoices `json:"choices,omitempty"`
}

type CompletionModuleResponseUsage struct {
	PromptTokens     int64 `json:"prompt_tokens,omitempty"`
	CompletionTokens int64 `json:"completion_tokens,omitempty"`
	TotalTokens      int64 `json:"total_tokens,omitempty"`
}

type CompletionModuleResponseChoices struct {
	Text         string `json:"text,omitempty"`
	Logprobs     int    `json:"logprobs,omitempty"` // 最大值: 5
	FinishReason string `json:"finish_reason,omitempty"`
	Index        int    `json:"index,omitempty"`
}

func SendCompletion(body TextDavinci003, apiKey string) (CompletionModuleResponse, error) {
	var completionModuleResponse CompletionModuleResponse
	b, err := json.Marshal(body)
	if err != nil {
		return completionModuleResponse, err
	}
	req, err := http.NewRequest("POST", completionApi, bytes.NewBuffer(b))
	if err != nil {
		return completionModuleResponse, err
	}
	
	req.Header.Set("Content-Type", completionApiContentType)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return completionModuleResponse, err
	}
	defer resp.Body.Close()
	
	respBodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return completionModuleResponse, err
	}
	
	if resp.StatusCode != http.StatusOK {
		return completionModuleResponse, errors.New(string(respBodyByte))
	}
	
	err = json.Unmarshal(respBodyByte, &completionModuleResponse)
	if err != nil {
		return completionModuleResponse, err
	}
	return completionModuleResponse, nil
}

func SendChat(body GptTurbo0301, apiKey string) (ChatModuleResponse, error) {
	var chatModuleResponse ChatModuleResponse
	b, err := json.Marshal(body)
	if err != nil {
		return chatModuleResponse, err
	}
	
	req, err := http.NewRequest("POST", chatApi, bytes.NewBuffer(b))
	if err != nil {
		return chatModuleResponse, err
	}
	
	req.Header.Set("Content-Type", chatApiContentType)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return chatModuleResponse, err
	}
	defer resp.Body.Close()
	
	respBodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return chatModuleResponse, err
	}
	
	if resp.StatusCode != http.StatusOK {
		return chatModuleResponse, errors.New(string(respBodyByte))
	}
	
	err = json.Unmarshal(respBodyByte, &chatModuleResponse)
	if err != nil {
		return chatModuleResponse, err
	}
	return chatModuleResponse, nil
}
