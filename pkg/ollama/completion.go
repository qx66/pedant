package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// 两种API都可以进行结构化输出，即，可以根据需求生成Json等数据

const (
	generateCompletionUri     = "/api/generate"
	generateChatCompletionUri = "/api/chat"
)

type Client struct {
	BasicUrl string `json:"basicUrl,omitempty"` // 基础Url，比如: http://127.0.0.1:11434
}

func NewClient(basicUrl string) *Client {
	return &Client{BasicUrl: basicUrl}
}

// GenerateCompletion

type GenerateCompletionReq struct {
	Model  string `json:"model,omitempty"` // require
	Prompt string `json:"prompt,omitempty"`
	Suffix string `json:"suffix,omitempty"`
	Images string `json:"images,omitempty"`
	//Format CompletionFormat `json:"format,omitempty"` // the format to return a response in. Format can be json or a JSON schema
	Stream bool `json:"stream"`        // default: true
	Raw    bool `json:"raw,omitempty"` // raw - if true no formatting will be applied to the prompt. You may choose to use the raw parameter if you are specifying a full templated prompt in your request to the API
	// options - additional model parameters listed in the documentation for the Modelfile such as temperature
	// system - system message to
	// template - the prompt template to use
	// keep_alive: controls how long the model will stay loaded into memory following the request (default: 5m)
}

type GenerateCompletionReturnStructReq struct {
	Model  string           `json:"model,omitempty"` // require
	Prompt string           `json:"prompt,omitempty"`
	Suffix string           `json:"suffix,omitempty"`
	Images string           `json:"images,omitempty"`
	Format CompletionFormat `json:"format,omitempty"` // the format to return a response in. Format can be json or a JSON schema
	Stream bool             `json:"stream"`           // default: true
	Raw    bool             `json:"raw,omitempty"`    // raw - if true no formatting will be applied to the prompt. You may choose to use the raw parameter if you are specifying a full templated prompt in your request to the API
}

type CompletionFormat struct {
	Type       string   `json:"type,omitempty"`
	Properties any      `json:"properties,omitempty"`
	Required   []string `json:"required,omitempty"`
}

type CompletionResponse struct {
	Model              string `json:"model,omitempty"`
	CreatedAt          string `json:"created_at,omitempty"`
	Response           string `json:"response,omitempty"`
	Done               bool   `json:"done,omitempty"`
	DoneReason         string `json:"done_reason,omitempty"`
	Context            []int  `json:"context,omitempty"`              // 此响应中使用的对话编码，可在下一个请求中发送以保留对话记忆
	TotalDuration      int64  `json:"total_duration,omitempty"`       // 生成响应所花费的时间
	LoadDuration       int64  `json:"load_duration,omitempty"`        // 加载模型所用的时间（以纳秒为单位）
	PromptEvalCount    int    `json:"prompt_eval_count,omitempty"`    // 提示中的标记数
	PromptEvalDuration int64  `json:"prompt_eval_duration,omitempty"` // 评估提示所花费的时间（以纳秒为单位）
	EvalCount          int    `json:"eval_count,omitempty"`           // 响应中的令牌数量
	EvalDuration       int64  `json:"eval_duration,omitempty"`        // 生成响应所用的时间（以纳秒为单位）
}

func (client *Client) GenerateCompletion(ctx context.Context, model string, prompt string) error {
	req := GenerateCompletionReq{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	}
	
	url := fmt.Sprintf("%s%s", client.BasicUrl, generateCompletionUri)
	
	reqByte, err := json.Marshal(req)
	if err != nil {
		return err
	}
	
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqByte))
	if err != nil {
		return err
	}
	
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("httpCodeError, %d", resp.StatusCode))
	}
	
	defer resp.Body.Close()
	respBodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	
	fmt.Println("respBody: ", string(respBodyByte))
	return nil
}

// 解析格式
// AI 可以根据你的内容，解析成Json格式返回

func (client *Client) GenerateCompletionReturnStruct(ctx context.Context, reqParams GenerateCompletionReturnStructReq) (string, error) {
	completionResponse := CompletionResponse{}
	url := fmt.Sprintf("%s%s", client.BasicUrl, generateCompletionUri)
	
	reqByte, err := json.Marshal(reqParams)
	if err != nil {
		return "", err
	}
	
	cli := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqByte))
	
	// 如果需要中文响应，需要设置header language 为 zh-CN -- 测试下来，效果不是很好
	// 问题中如果没有英文，可以生成中文内容。
	req.Header.Set("language", "zh-CN")
	//resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqByte))
	if err != nil {
		return "", err
	}
	
	resp, err := cli.Do(req)
	if err != nil {
		return "", err
	}
	
	if resp.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("httpCodeError, %d", resp.StatusCode))
	}
	
	defer resp.Body.Close()
	respBodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	
	err = json.Unmarshal(respBodyByte, &completionResponse)
	if err != nil {
		return "", err
	}
	
	fmt.Println("respBodyByte: ", string(respBodyByte))
	if completionResponse.Done != true {
		return "", errors.New("done is not true, incomplete")
	}
	
	if completionResponse.DoneReason != "stop" {
		return "", errors.New(fmt.Sprintf("done reason: %s", completionResponse.DoneReason))
	}
	
	return completionResponse.Response, nil
}

// ChatCompletion

type ChatCompletionReq struct {
	Model    string                  `json:"model,omitempty"`
	Messages []ChatCompletionMessage `json:"messages,omitempty"`
	Format   string                  `json:"format,omitempty"` // the format to return a response in. Format can be json or a JSON schema.
	Stream   bool                    `json:"stream"`           // if false the response will be returned as a single response object, rather than a stream of objects
	Tools    []byte                  `json:"tools,omitempty"`
	//options // additional model parameters listed in the documentation for the Modelfile such as temperature
	//keep_alive // controls how long the model will stay loaded into memory following the request (default: 5m)
}

type ChatCompletionMessage struct {
	Role      string `json:"role,omitempty"`       // the role of the message, either system, user, assistant, or tool
	Content   string `json:"content,omitempty"`    // the content of the message
	Images    []byte `json:"images,omitempty"`     // (optional): a list of images to include in the message (for multimodal models such as llava)
	ToolCalls []byte `json:"tool_calls,omitempty"` // (optional): a list of tools in JSON that the model wants to use
}

type ChatCompletionResponse struct {
	Model              string                        `json:"model,omitempty"`
	CreatedAt          string                        `json:"created_at,omitempty"`
	Message            ChatCompletionResponseMessage `json:"message,omitempty"`
	DoneReason         string                        `json:"done_reason,omitempty"`
	Done               bool                          `json:"done,omitempty"`
	TotalDuration      int64                         `json:"total_duration,omitempty"`       // 生成响应所花费的时间
	LoadDuration       int64                         `json:"load_duration,omitempty"`        // 加载模型所用的时间（以纳秒为单位）
	PromptEvalCount    int                           `json:"prompt_eval_count,omitempty"`    // 提示中的标记数
	PromptEvalDuration int64                         `json:"prompt_eval_duration,omitempty"` // 评估提示所花费的时间（以纳秒为单位）
	EvalCount          int                           `json:"eval_count,omitempty"`           // 响应中的令牌数量
	EvalDuration       int64                         `json:"eval_duration,omitempty"`        // 生成响应所用的时间（以纳秒为单位）
}

type ChatCompletionResponseMessage struct {
	Role    string `json:"role,omitempty"` // assistant
	Content string `json:"content,omitempty"`
}

// 该方法不适合 Stream 模式, Stream 需要单独处理每一条请求进行合并，并判断最后一条记录

func (client *Client) ChatCompletion(ctx context.Context, req *ChatCompletionReq) (ChatCompletionResponse, error) {
	result := ChatCompletionResponse{}
	req.Stream = false
	reqByte, err := json.Marshal(req)
	
	if err != nil {
		return result, err
	}
	
	url := fmt.Sprintf("%s%s", client.BasicUrl, generateChatCompletionUri)
	
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqByte))
	if err != nil {
		return result, err
	}
	
	if resp.StatusCode != 200 {
		return result, errors.New(fmt.Sprintf("httpCodeError, %d", resp.StatusCode))
	}
	
	defer resp.Body.Close()
	respBodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}
	
	err = json.Unmarshal(respBodyByte, &result)
	if err != nil {
		//fmt.Println("respBody: ", string(respBodyByte))
		return result, err
	}
	
	if result.Done == false {
		return result, errors.New(result.DoneReason)
	}
	
	return result, nil
}

// 解析格式
// AI 可以根据你的内容，解析成Json格式返回
