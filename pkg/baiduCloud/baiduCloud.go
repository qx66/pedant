package baiduCloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/startopsz/rule/pkg/http"
)

const (
	ChatRoleUser      = "user"
	ChatRoleAssistant = "assistant"
)

const (
	stableDiffusionXLImageApi = "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/text2image/sd_xl"
	ernieBotApi               = "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/eb-instant"
	ernieBot4Api              = "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/completions_pro"
)

type AccessToken struct {
	RefreshToken  string `json:"refresh_token,omitempty"`
	ExpiresIn     int    `json:"expires_in,omitempty"`
	SessionKey    string `json:"session_key,omitempty"`
	AccessToken   string `json:"access_token,omitempty"`
	Scope         string `json:"scope,omitempty"`
	SessionSecret string `json:"session_secret,omitempty"`
}

// 获取百度千帆调用接口的 AccessToken
// access_token默认有效期30天，单位是秒，生产环境注意及时刷新。

func GetQianFanAccessToken(apiKey, secretKey string) (AccessToken, error) {
	var accessToken AccessToken
	baseUrl := "https://aip.baidubce.com/oauth/2.0/token"
	url := fmt.Sprintf("%s?grant_type=client_credentials&client_id=%s&client_secret=%s", baseUrl, apiKey, secretKey)
	
	req := http.Req{
		Method:  http.Post,
		Url:     url,
		Timeout: 5,
	}
	
	resp, err := req.Do()
	if err != nil {
		return accessToken, err
	}
	
	if resp.StatusCode != 200 {
		return accessToken, errors.New(
			fmt.Sprintf("http status: %d, message: %s", resp.StatusCode, string(resp.Body)),
		)
	}
	
	err = json.Unmarshal(resp.Body, &accessToken)
	if err != nil {
		return accessToken, err
	}
	
	return accessToken, nil
}

type ERNIEBotTurboReq struct {
	Messages []ERNIEBotTurboMessage `json:"messages,omitempty"`
	Stream   bool                   `json:"stream,omitempty"` // 默认false
	System   string                 `json:"system,omitempty"`
	UserId   string                 `json:"user_id,omitempty"`
	//Temperature  float32                `json:"temperature,omitempty"`   // 默认0.95，范围 (0, 1.0]，不能为0, 建议top_p和temperature不要同时更改
	//TopP         float32                `json:"top_p,omitempty"`         // 默认0.8，取值范围 [0, 1.0], 建议top_p和temperature不要同时更改
	//PenaltyScore bool                   `json:"penalty_score,omitempty"` // 通过对已生成的token增加惩罚，减少重复生成的现象。 默认1.0，取值范围：[1.0, 2.0]
	//UserId       string                 `json:"user_id,omitempty"`
}

type ERNIEBotTurboMessage struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
	Name    string `json:"name,omitempty"`
}

type ERNIEBotTurboResponse struct {
	Id               string                     `json:"id,omitempty"`                 // 本轮对话的id
	Object           string                     `json:"object,omitempty"`             // 回包类型, chat.completion：多轮对话返回
	Created          int64                      `json:"created,omitempty"`            // 时间戳
	SentenceId       int64                      `json:"sentence_id,omitempty"`        // 表示当前子句的序号。只有在流式接口模式下会返回该字段
	IsEnd            bool                       `json:"is_end,omitempty"`             // 表示当前子句是否是最后一句。只有在流式接口模式下会返回该字段
	IsTruncated      bool                       `json:"is_truncated,omitempty"`       // 当前生成的结果是否被截断
	Result           string                     `json:"result,omitempty"`             // 对话返回结果
	NeedClearHistory bool                       `json:"need_clear_history,omitempty"` // 表示用户输入是否存在安全，是否关闭当前会话，清理历史会话信息。 true：是，表示用户输入存在安全风险，建议关闭当前会话，清理历史会话信息。 false：否，表示用户输入无安全风险
	Usage            ERNIEBotTurboResponseUsage `json:"usage,omitempty"`              // token统计信息，token数 = 汉字数+单词数*1.3 （仅为估算逻辑）
	ErrorCode        int                        `json:"error_code,omitempty"`         // 错误码
	ErrorMsg         string                     `json:"error_msg,omitempty"`          // 错误描述信息，帮助理解和解决发生的错误
}

type ERNIEBotTurboResponseUsage struct {
	PromptTokens     int `json:"prompt_tokens,omitempty"`     // 问题tokens数
	CompletionTokens int `json:"completion_tokens,omitempty"` // 回答tokens数
	TotalTokens      int `json:"total_tokens,omitempty"`      // tokens总数
}

func SendERNIEBotTurbo(accessToken string, body ERNIEBotTurboReq) (ERNIEBotTurboResponse, error) {
	var ernieBotResp ERNIEBotTurboResponse
	url := fmt.Sprintf("%s?access_token=%s", ernieBot4Api, accessToken)
	
	bodyByte, err := json.Marshal(body)
	if err != nil {
		return ernieBotResp, err
	}
	
	//
	req := http.Req{
		Method:  http.Post,
		Url:     url,
		Body:    bodyByte,
		Timeout: 120,
	}
	
	resp, err := req.Do()
	if err != nil {
		return ernieBotResp, err
	}
	//
	
	if resp.StatusCode != 200 {
		return ernieBotResp, errors.New(
			fmt.Sprintf("http status: %d, message: %s", resp.StatusCode, string(resp.Body)),
		)
	}
	
	//
	err = json.Unmarshal(resp.Body, &ernieBotResp)
	if err != nil {
		return ernieBotResp, err
	}
	
	if ernieBotResp.ErrorCode != 0 {
		return ernieBotResp, errors.New(ernieBotResp.ErrorMsg)
	}
	
	//
	return ernieBotResp, nil
}

// https://cloud.baidu.com/doc/WENXINWORKSHOP/s/Klkqubb9w

type StableDiffusionSize string
type StableDiffusionSamplerIndex string

const (
	StableDiffusionSize768x768   StableDiffusionSize = "768x768"
	StableDiffusionSize768x1024  StableDiffusionSize = "768x1024"
	StableDiffusionSize1024x768  StableDiffusionSize = "1024x768"
	StableDiffusionSize576x1024  StableDiffusionSize = "576x1024"
	StableDiffusionSize1024x576  StableDiffusionSize = "1024x576"
	StableDiffusionSize1024x1024 StableDiffusionSize = "1024x1024"
	
	StableDiffusionSamplerIndexEuler       StableDiffusionSamplerIndex = "Euler"
	StableDiffusionSamplerIndexEulerA      StableDiffusionSamplerIndex = "Euler a"
	StableDiffusionSamplerIndexDPM2M       StableDiffusionSamplerIndex = "DPM++ 2M"
	StableDiffusionSamplerIndexDPM2MKarras StableDiffusionSamplerIndex = "DPM++ 2M Karras"
	StableDiffusionSamplerIndexLMSKarras   StableDiffusionSamplerIndex = "LMS Karras"
	StableDiffusionSamplerIndexDPMSDE      StableDiffusionSamplerIndex = "DPM++ SDE"
)

type StableDiffusionXLReq struct {
	Prompt         string                      `json:"prompt,omitempty"`          // require 提示词，即用户希望图片包含的元素。长度限制为1024字符，建议中文或者英文单词总数量不超过150个
	NegativePrompt string                      `json:"negative_prompt,omitempty"` // 反向提示词，即用户希望图片不包含的元素。长度限制为1024字符，建议中文或者英文单词总数量不超过150个
	Size           StableDiffusionSize         `json:"size,omitempty"`            // 生成图片长宽，默认值 1024x1024
	Steps          int                         `json:"steps,omitempty"`           // 生成图片数量，说明: 默认值为1,取值范围为1-4
	N              int                         `json:"n,omitempty"`               // 迭代轮次，说明: 默认值为20, 取值范围为10-50
	SamplerIndex   StableDiffusionSamplerIndex `json:"sampler_index,omitempty"`   // 采样方式，默认值：Euler a
}

type StableDiffusionXLResponse struct {
	Id      string                          `json:"id,omitempty"`      // 请求的id
	Object  string                          `json:"object,omitempty"`  // 回包类型。image：图像生成返回
	Created int64                           `json:"created,omitempty"` // 时间戳
	Data    []StableDiffusionXLResponseData `json:"data,omitempty"`    // 生成图片结果
	Usage   ERNIEBotTurboResponseUsage      `json:"usage,omitempty"`   // oken统计信息
}

type StableDiffusionXLResponseData struct {
	Object   string `json:"object,omitempty"`    // 固定值"image"
	B64Image string `json:"b64_image,omitempty"` // 图片base64编码内容
	Index    int    `json:"index,omitempty"`     // 序号
}

func GenerateStableDiffusionXLImage(accessToken string, body StableDiffusionXLReq) (StableDiffusionXLResponse, error) {
	var response StableDiffusionXLResponse
	
	url := fmt.Sprintf("%s?access_token=%s", stableDiffusionXLImageApi, accessToken)
	
	bodyByte, err := json.Marshal(body)
	if err != nil {
		return response, err
	}
	
	//
	req := http.Req{
		Method:  http.Post,
		Url:     url,
		Body:    bodyByte,
		Timeout: 120,
	}
	
	resp, err := req.Do()
	if err != nil {
		return response, err
	}
	//
	
	if resp.StatusCode != 200 {
		return response, errors.New(
			fmt.Sprintf("http status: %d, message: %s", resp.StatusCode, string(resp.Body)),
		)
	}
	
	//
	err = json.Unmarshal(resp.Body, &response)
	//
	return response, err
}
