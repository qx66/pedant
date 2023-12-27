package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"io"
	"os"
)

type Region string
type Model string

const (
	RegionNorthAmericaNortheast1 Region = "northamerica-northeast1" // 加拿大蒙特利尔
	RegionUsCentral1             Region = "us-central1"             // 爱荷华
	RegionUsWest4                Region = "us-west4"                // 内华达州，拉斯维加
	RegionUsEast4                Region = "us-east4"                // 北弗吉尼亚
	RegionUsWest1                Region = "us-west1"                // 俄勒冈
	RegionAsiaNortheast3         Region = "asia-northeast3"         // 韩国首尔
	RegionAsiaSoutheast1         Region = "asia-southeast1"         // 新加坡
	RegionAsiaNortheast1         Region = "asia-northeast1"         // 日本东京
	
	ModelTextBison     Model = "text-bison"     // 文本
	ModelTextUnicorn   Model = "text-unicorn"   // 文本
	ModelChatBison     Model = "chat-bison"     // 聊天
	ModelCodeBison     Model = "code-bison"     // 编码
	ModelCodeChatBison Model = "codechat-bison" // 代码聊天
	ModelCodeGecko     Model = "code-gecko"
	
	ModelImageGeneration Model = "imagegeneration" // 图片生成
	ModelImageText       Model = "imagetext"       // 图片说明
)

// 如果要调用最新版本的 textembedding-gecko 模型，您必须添加 @latest 作为后缀。例如 textembedding-gecko@latest。

// 使用 CountTokens API 无需付费或配额限制。CountTokens API 和 ComputeTokens API 的最大配额为每分钟 3000 个请求。

// https://github.com/GoogleCloudPlatform/golang-samples/tree/0e5964e6cbbbc139ed10f7fde5eb41b83147b2b7/vertexai

type Client struct {
	cli *genai.Client
}

func NewClient(apiKey string) (*Client, error) {
	cli, err := genai.NewClient(context.Background(), option.WithAPIKey(apiKey),
	)
	
	if err != nil {
		return &Client{}, err
	}
	
	return &Client{
		cli: cli,
	}, nil
}

// 文本

func (client *Client) Text(ctx context.Context, text string) error {
	
	model := client.cli.GenerativeModel("gemini-pro")
	
	resp, err := model.GenerateContent(ctx, genai.Text(text))
	
	fmt.Println(resp)
	
	if err != nil {
		return err
	}
	
	rb, err := json.MarshalIndent(resp, "", "")
	if err != nil {
		return err
	}
	
	fmt.Println(string(rb))
	return nil
}

// Gemini 提供了一个多模态模型 (gemini-pro-vision)，因此您可以同时输入文本和图片

func (client *Client) Multimodal(ctx context.Context, imagePaths []string, text string) error {
	model := client.cli.GenerativeModel("gemini-pro-vision")
	
	var prompt []genai.Part
	
	for _, imagePath := range imagePaths {
		//
		imgF, err := os.Open(imagePath)
		if err != nil {
			return err
		}
		
		//
		imgData, err := io.ReadAll(imgF)
		if err != nil {
			return err
		}
		
		//
		/*
			_, format, err := image.Decode(imgF)
			if err != nil {
				return err
			}
		*/
		
		prompt = append(prompt, genai.ImageData("png", imgData))
	}
	
	prompt = append(prompt, genai.Text(text))
	
	resp, err := model.GenerateContent(ctx, prompt...)
	if err != nil {
		return err
	}
	
	fmt.Println("resp: ", resp)
	rb, err := json.MarshalIndent(resp, "", "")
	if err != nil {
		return err
	}
	
	fmt.Println(string(rb))
	
	return nil
}

// 对话

func (client *Client) Chat(ctx context.Context, content []*genai.Content, text string) error {
	model := client.cli.GenerativeModel("gemini-pro")
	cs := model.StartChat()
	cs.History = content
	resp, err := cs.SendMessage(ctx, genai.Text(text))
	if err != nil {
		return err
	}
	
	fmt.Println("resp: ", resp)
	rb, err := json.MarshalIndent(resp, "", "")
	if err != nil {
		return err
	}
	
	fmt.Println(string(rb))
	return nil
}
