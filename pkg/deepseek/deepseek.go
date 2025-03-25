package deepseek

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// https://api-docs.deepseek.com/zh-cn/

const (
	baseUrl = "https://api.deepseek.com"
)

// model:
// deepseek-chat -> DeepSeek-V3
// deepseek-reasoner -> DeepSeek-R1

type CompletionRequest struct {
	Model            string              `json:"model,omitempty"`
	Messages         []CompletionMessage `json:"messages,omitempty"`
	Stream           bool                `json:"stream,omitempty"`            // 如果设置为 True，将会以 SSE（server-sent events）的形式以流式发送消息增量。消息流以 data: [DONE] 结尾。
	FrequencyPenalty int                 `json:"frequency_penalty,omitempty"` // default:0, 介于 -2.0 和 2.0 之间的数字。如果该值为正，那么新 token 会根据其在已有文本中的出现频率受到相应的惩罚，降低模型重复相同内容的可能性。
	MaxTokens        int                 `json:"max_tokens,omitempty"`        // default:4096, 介于 1 到 8192 间的整数，限制一次请求中模型生成 completion 的最大 token 数。输入 token 和输出 token 的总长度受模型的上下文长度的限制。
	PresencePenalty  int                 `json:"presence_penalty,omitempty"`  // default:0, 介于 -2.0 和 2.0 之间的数字。如果该值为正，那么新 token 会根据其是否已在已有文本中出现受到相应的惩罚，从而增加模型谈论新主题的可能性。
	ResponseFormat   string              `json:"response_format,omitempty"`   // default: text, 一个 object，指定模型必须输出的格式。 Must be one of text or json_object.
	Temperature      int                 `json:"temperature,omitempty"`       // default: 1, 采样温度，介于 0 和 2 之间。更高的值，如 0.8，会使输出更随机，而更低的值，如 0.2，会使其更加集中和确定。 我们通常建议可以更改这个值或者更改 top_p，但不建议同时对两者进行修改。
	TopP             int                 `json:"top_p,omitempty"`             // default: 1, 作为调节采样温度的替代方案，模型会考虑前 top_p 概率的 token 的结果。所以 0.1 就意味着只有包括在最高 10% 概率中的 token 会被考虑。 我们通常建议修改这个值或者更改 temperature，但不建议同时对两者进行修改。
	Logprobs         bool                `json:"logprobs,omitempty"`          // 是否返回所输出 token 的对数概率。如果为 true，则在 message 的 content 中返回每个输出 token 的对数概率。
	TopLogprobs      int                 `json:"top_logprobs,omitempty"`      // 一个介于 0 到 20 之间的整数 N，指定每个输出位置返回输出概率 top N 的 token，且返回这些 token 的对数概率。指定此参数时，logprobs 必须为 true。
	// stop
	// stream_options // 流式输出相关选项。只有在 stream 参数为 true 时，才可设置此参数。
	// tools
	// tool_choice
}

type CompletionMessage struct {
	Role    string `json:"role,omitempty"` // system, user
	Content string `json:"content,omitempty"`
}

type CompletionRequestResponseFormat struct {
	Type string `json:"type,omitempty"` // text
}

type CompletionResponse struct {
	Id      string                     `json:"id,omitempty"`
	Choices []CompletionResponseChoice `json:"choices,omitempty"`
	Created int64                      `json:"created,omitempty"`
	Model   string                     `json:"model,omitempty"`
	Object  string                     `json:"object,omitempty"`
	Usage   CompletionResponseUsage    `json:"usage,omitempty"`
}

type CompletionResponseChoice struct {
	FinishReason string            `json:"finish_reason,omitempty"`
	Index        int               `json:"index,omitempty"`
	Message      CompletionMessage `json:"message,omitempty"`
}

type CompletionResponseUsage struct {
	CompletionTokens int `json:"completion_tokens,omitempty"`
	PromptTokens     int `json:"prompt_tokens,omitempty"`
	TotalTokens      int `json:"total_tokens,omitempty"`
}

func Completion(req CompletionRequest, apiKey string) {
	url := fmt.Sprintf("%s/chat/completions", baseUrl)
	
	header := make(http.Header)
	header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	
	reqByte, err := json.Marshal(req)
	if err != nil {
		return
	}
	
	r, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqByte))
	r.Header = header
	
	cli := http.DefaultClient
	resp, err := cli.Do(r)
	if err != nil {
		return
	}
	
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return
	}
	
	respByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	
	fmt.Println("resp: ", string(respByte))
}
