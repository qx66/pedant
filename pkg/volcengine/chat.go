package volcengine

import (
	"context"
	"fmt"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"time"
)

// 火山引擎
// https://www.volcengine.com/docs/82379/1319853

type Client struct {
	cli *arkruntime.Client
}

func NewClient(apiKey string, timeout int) *Client {
	cli := arkruntime.NewClientWithApiKey(apiKey, arkruntime.WithTimeout(time.Duration(timeout)*time.Second))
	return &Client{
		cli: cli,
	}
}

func (client *Client) ChatCompletions(ctx context.Context) {
	req := model.ChatCompletionRequest{
		Model: "<YOUR_ENDPOINT_ID>",
		Messages: []*model.ChatCompletionMessage{
			{
				Role: model.ChatMessageRoleSystem,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String("你是豆包，是由字节跳动开发的 AI 人工智能助手"),
				},
			},
			{
				Role: model.ChatMessageRoleUser,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String("常见的十字花科植物有哪些？"),
				},
			},
		},
	}
	
	resp, err := client.cli.CreateChatCompletion(ctx, req)
	if err != nil {
		fmt.Printf("standard chat error: %v\n", err)
		return
	}
	
	fmt.Println(*resp.Choices[0].Message.Content.StringValue)
	
}
