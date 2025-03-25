package baiduCloud

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	chatV2Url = "https://qianfan.baidubce.com/v2/chat/completions"
)

type ChatReq struct {
	Model               string        `json:"model,omitempty"`
	Messages            []ChatMessage `json:"messages,omitempty"`
	Stream              bool          `json:"stream,omitempty"`
	Temperature         float32       `json:"temperature,omitempty"`           // 选填参数 - 说明：（1）较高的数值会使输出更加随机，而较低的数值会使其更加集中和确定。（2）默认0.95，范围 (0, 1.0]，不能为0。（3）建议该参数和top_p只设置1个。（4）建议top_p和temperature不要同时更改。
	TopP                float32       `json:"top_p,omitempty"`                 // 选填参数 - 说明：（1）影响输出文本的多样性，取值越大，生成文本的多样性越强。（2）默认0.8，取值范围 [0, 1.0]。（3）建议该参数和temperature只设置1个。（4）建议top_p和temperature不要同时更改。
	PenaltyScore        float32       `json:"penalty_score,omitempty"`         // 选填参数 - 通过对已生成的token增加惩罚，减少重复生成的现象。说明：（1）值越大表示惩罚越大。（2）默认1.0，取值范围：[1.0, 2.0]。
	System              string        `json:"system,omitempty"`                // 选填参数 - 模型人设，主要用于人设设定，例如，你是xxx公司制作的AI助手，说明：（1）长度限制1024个字符（2）如果使用functions参数，不支持设定人设system
	MaxCompletionTokens int           `json:"max_completion_tokens,omitempty"` // 选填参数 - 指定模型最大输出token数，范围[2, 2048]
	Functions           string        `json:"functions,omitempty"`             // 选填参数 - 一个可触发函数的描述列表
	Stop                string        `json:"stop,omitempty"`                  // 选填参数 - 生成停止标识，当模型生成结果以stop中某个元素结尾时，停止文本生成
	DisableSearch       bool          `json:"disable_search,omitempty"`        // 选填参数 - 是否强制关闭实时搜索功能，默认false，表示不关闭
	EnableCitation      bool          `json:"enable_citation,omitempty"`       // 选填参数 - 是否开启上角标返回
	
}

type ChatMessage struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

func ChatV2(ctx context.Context, appid, authorization string, reqBody ChatReq) error {
	
	bodyByte, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	
	//
	client := &http.Client{
		Timeout: time.Duration(120) * time.Second,
	}
	
	req, err := http.NewRequest(http.MethodPost, chatV2Url, bytes.NewBuffer(bodyByte))
	//
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authorization))
	req.Header.Set("appid", appid)
	
	resp, err := client.Do(req)
	
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	//
	respByte, err := io.ReadAll(resp.Body)
	
	if resp.StatusCode != 200 {
		fmt.Println("body: ", string(respByte))
		return errors.New(fmt.Sprintf("httpCode: %d", resp.StatusCode))
	}
	
	fmt.Println("body: ", string(respByte))
	//
	
	return nil
}
