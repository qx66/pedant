package alibabaCloud

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	fullModelApi       = "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions"
	defaultModel       = "qwen-omni-turbo"
	defaultContentType = "application/json"
)

type Client struct {
	authorization string
	cli           *http.Client
}

func NewClient(apiKey string) *Client {
	if apiKey == "" {
		panic("apiKey is null")
	}
	
	return &Client{
		authorization: fmt.Sprintf("Bearer %s", apiKey),
		cli:           http.DefaultClient,
	}
}

type ImageModelReq struct {
	Model          string              `json:"model,omitempty"`
	Messages       []ImageModelMessage `json:"messages,omitempty"`
	ResponseFormat ChatResponseFormat  `json:"response_format,omitempty"`
	Stream         bool                `json:"stream,omitempty"`
	StreamOptions  StreamOptions       `json:"stream_options"`
}

type ChatResponseFormat struct {
	Type string `json:"type,omitempty"`
}

type StreamOptions struct {
	IncludeUsage bool `json:"include_usage"`
}

type ImageModelMessage struct {
	Role    string                     `json:"role,omitempty"`
	Content []ImageModelMessageContent `json:"content,omitempty"`
}

type ImageModelMessageContent struct {
	Type     string                           `json:"type,omitempty"` // text or image_url
	Text     string                           `json:"text,omitempty"`
	ImageUrl ImageModelMessageContentImageUrl `json:"image_url,omitempty"`
}

type ImageModelMessageContentImageUrl struct {
	Url string `json:"url,omitempty"`
}

type StreamResponse struct {
	Choices []StreamResponseChoices `json:"choices"`
	Object  string                  `json:"object"`
	Usage   StreamResponseUsage     `json:"usage"`
	Created int                     `json:"created"`
	Model   string                  `json:"model"`
	Id      string                  `json:"id"`
	//SystemFingerprint interface{}             `json:"system_fingerprint"`
}

type StreamResponseChoices struct {
	FinishReason string                     `json:"finish_reason"`
	Delta        StreamResponseChoicesDelta `json:"delta"`
	Index        int                        `json:"index"`
	//Logprobs     interface{}                `json:"logprobs"`
}

type StreamResponseChoicesDelta struct {
	Content string `json:"content"`
}

type StreamResponseUsage struct {
	PromptTokens            int                                        `json:"prompt_tokens"`
	CompletionTokens        int                                        `json:"completion_tokens"`
	TotalTokens             int                                        `json:"total_tokens"`
	CompletionTokensDetails StreamResponseUsageCompletionTokensDetails `json:"completion_tokens_details"`
	PromptTokensDetails     StreamResponseUsagePromptTokensDetails     `json:"prompt_tokens_details"`
}

type StreamResponseUsageCompletionTokensDetails struct {
	TextTokens int `json:"text_tokens"`
}

type StreamResponseUsagePromptTokensDetails struct {
	TextTokens  int `json:"text_tokens"`
	ImageTokens int `json:"image_tokens"`
}

const (
	lotteryTicketSystem = `
你需要提取出彩票类型名字(name,为string类型)、期号(issue,为string类型)、期号数字(issueNumber,为int类型)、开奖日期(drawDate,为string类型)、单式票(tickets,为array string类型)、金额(amount,为string类型)、金额值(amountNumber,为int类型)
示例：
Q:图片为大乐透彩票
A:{"name": "超级大乐透", "issue": "第25030期","issueNumber": 25030, "drawDate": "2025年03月22日", "tickets": ["05 08 10 22 34 + 02 08","08 11 12 15 29 + 01 04","01 02 05 20 31 + 07 08","07 13 15 18 24 + 05 09"],"amount": "合计8元", "amountNumber": "8"}
`
)

type LotteryTicket struct {
	Name         string   `json:"name"`
	Issue        string   `json:"issue"`
	IssueNumber  int      `json:"issueNumber"`
	DrawDate     string   `json:"drawDate"`
	Tickets      []string `json:"tickets"`
	Amount       string   `json:"amount"`
	AmountNumber int      `json:"amountNumber"`
}

func (client *Client) ImageCompletions(imageUrl string, userTextContent string) {
	imageModelReq := ImageModelReq{
		Model: defaultModel,
		Messages: []ImageModelMessage{
			{
				Role: "system",
				Content: []ImageModelMessageContent{
					{
						Type: "text",
						Text: lotteryTicketSystem,
					},
				},
			},
			{
				Role: "user",
				Content: []ImageModelMessageContent{
					{
						Type: "image_url",
						ImageUrl: ImageModelMessageContentImageUrl{
							Url: imageUrl,
						},
					},
					{
						Type: "text",
						Text: userTextContent,
					},
				},
			},
		},
		ResponseFormat: ChatResponseFormat{
			Type: "json_object",
		},
		Stream: true,
		StreamOptions: StreamOptions{
			IncludeUsage: true,
		},
	}
	
	imageModelReqByte, err := json.Marshal(&imageModelReq)
	if err != nil {
		return
	}
	
	req, err := http.NewRequest(http.MethodPost, fullModelApi, bytes.NewBuffer(imageModelReqByte))
	if err != nil {
		return
	}
	
	req.Header.Add("Content-Type", defaultContentType)
	req.Header.Add("Authorization", client.authorization)
	
	resp, err := client.cli.Do(req)
	if err != nil {
		return
	}
	
	defer resp.Body.Close()
	
	reader := bufio.NewReader(resp.Body)
	var contents []string
	
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break // 读取完毕或发生错误
		}
		// 解析返回的 JSON 数据
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "data:") {
			data := strings.TrimPrefix(line, "data:")
			//fmt.Println(data) // 输出 AI 响应
			r := StreamResponse{}
			err = json.Unmarshal([]byte(data), &r)
			if err != nil {
				fmt.Println("json unmarshal failed.")
				return
			}
			
			if r.Choices[0].FinishReason == "stop" {
				break
			}
			
			//fmt.Print(r.Choices[0].Delta.Content)
			contents = append(contents, r.Choices[0].Delta.Content)
		}
		
		// 处理 OpenAI 的流式结束标记
		if line == "data: [DONE]" {
			break
		}
	}
	
	content := strings.Join(contents, "")
	var lotteryTicket LotteryTicket
	
	err = json.Unmarshal([]byte(content), &lotteryTicket)
	if err != nil {
		return
	}
	
	fmt.Println("lotteryTicket: ", lotteryTicket)
	return
}
