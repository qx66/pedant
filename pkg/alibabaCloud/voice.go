package alibabaCloud

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"os"
	
	"net/http"
)

/*
bill: https://help.aliyun.com/zh/model-studio/developer-reference/cosyvoice-billing?spm=a2c4g.11186623.help-menu-2400256.d_3_3_6_0_4.733d289e6SPJpj
CosyVoice 系列模型: 2元/万字符
*/

/*
音色列表: https://help.aliyun.com/zh/model-studio/developer-reference/timbre-list?spm=a2c4g.11186623.help-menu-2400256.d_3_3_6_0_5.7539575eH3PXjY
*/

/*
关于建连开销和连接复用
  WebSocket服务支持连接复用以提升资源的利用效率，避免建立连接开销。
  当服务收到 run-task 指令后，将启动一个新的任务，并在任务完成时下发 task-finished 指令以结束该任务。结束任务后webSocket连接可以被复用，发送run-task指令开启下一个任务。

1. 在复用连接中的不同任务需要使用不同 task_id。
2. 如果在任务执行过程中发生失败，服务将依然下发 task-failed 指令，并关闭该连接。此时这个连接无法继续复用。
3. 如果在任务结束后60秒没有新的任务，连接会超时自动断开。
*/

/*

longhua -


*/

type VoiceUser string

const (
	LongWan       VoiceUser = "longwan"       // 龙婉
	LongCheng     VoiceUser = "longcheng"     // 龙橙
	LongHua       VoiceUser = "longhua"       // 龙华
	LongXiaoChun  VoiceUser = "longxiaochun"  //龙小淳
	LongXiaoXia   VoiceUser = "longxiaoxia"   // 龙小夏
	LongXiaoCheng VoiceUser = "longxiaocheng" // 龙小诚
	LongXiaoBai   VoiceUser = "longxiaobai"   // 龙小白
	LongLaoTie    VoiceUser = "longlaotie"    // 龙老铁
	LongShu       VoiceUser = "longshu"       // 龙书
	LongShuo      VoiceUser = "longshuo"      // 龙硕
	LongJing      VoiceUser = "longjing"      // 龙婧
	LongMiao      VoiceUser = "longmiao"      // 龙妙
	LongYue       VoiceUser = "longyue"       // 龙悦
	LongYuan      VoiceUser = "longyuan"      // 龙媛
	LongFei       VoiceUser = "longfei"       // 龙飞
	LongJieLiDou  VoiceUser = "longjielidou"  // 龙杰力豆
	LongTong      VoiceUser = "longtong"      // 龙彤
	Stella        VoiceUser = "loongstella"
	Bella         VoiceUser = "loongbella"
)

const (
	voiceWebSocketUrl = "wss://dashscope.aliyuncs.com/api-ws/v1/inference"
)

var dialer = websocket.DefaultDialer

type VoiceWebSocketClient struct {
	conn *websocket.Conn
}

// 定义结构体来表示JSON数据

type Header struct {
	Action       string                 `json:"action"`    // require 指令类型，可以选填
	TaskID       string                 `json:"task_id"`   // require 当次任务ID，随机生成的32位唯一ID。
	Streaming    string                 `json:"streaming"` // require 固定字符串："duplex"
	Event        string                 `json:"event"`     // 事件类型: task-started、result-generated、task-finished、task-failed
	ErrorCode    string                 `json:"error_code,omitempty"`
	ErrorMessage string                 `json:"error_message,omitempty"`
	Attributes   map[string]interface{} `json:"attributes"`
}

type Payload struct {
	TaskGroup  string     `json:"task_group"` // require 固定字符串："audio"。
	Task       string     `json:"task"`       // require 固定字符串："tts"。
	Function   string     `json:"function"`   // require 固定字符串："SpeechSynthesizer"。
	Model      string     `json:"model"`      // require 模型名称："cosyvoice-v1"。 //建议参考音色列表 当前支持的模型名：cosyvoice-v1，cosyvoice-v2。
	Input      Input      `json:"input"`      // require 需要合成的文本片段。
	Parameters Params     `json:"parameters"`
	Resources  []Resource `json:"resources"`
}

type Params struct {
	TextType string `json:"text_type"` // require  固定字符串：“PlainText”。
	//Voice      string `json:"voice"`       // require  发音人。
	Voice      VoiceUser `json:"voice"`       // require  发音人。
	Format     string    `json:"format"`      // require  音频编码格式，支持"pcm"、"wav"和"mp3"。
	SampleRate int       `json:"sample_rate"` // require  音频采样率，支持下述采样率: 8000, 16000, 22050, 24000, 44100, 48000。
	Volume     int       `json:"volume"`      // optional 音量，取值范围：0～100。默认值：50。
	Rate       int       `json:"rate"`        // optional 合成音频的语速，取值范围：0.5~2。默认值：1.0。
	Pitch      int       `json:"pitch"`       // optional 合成音频的语调，取值范围：0.5~2。 默认值：1.0。
}

type Resource struct {
	ResourceID   string `json:"resource_id"`
	ResourceType string `json:"resource_type"`
}

type Input struct {
	Text string `json:"text"`
}

type Event struct {
	Header  Header  `json:"header"`
	Payload Payload `json:"payload"`
}

func NewWebSocketCli(apiKey string) (*VoiceWebSocketClient, error) {
	headers := make(http.Header)
	headers.Add("Authorization", fmt.Sprintf("bearer %s", apiKey))
	headers.Add("user-agent", "StartOps")
	//headers["X-DashScope-WorkSpace"] = ""
	headers.Add("X-DashScope-DataInspection", "enable")
	
	conn, _, err := dialer.Dial(voiceWebSocketUrl, headers)
	
	return &VoiceWebSocketClient{
		conn: conn,
	}, err
}

// 发送run-task指令
// 1. 发送run-task指令：开启语音合成任务

func (voiceWebSocketClient *VoiceWebSocketClient) RunTask(voiceUser VoiceUser) (string, error) {
	runTaskCmd, taskID, err := generateRunTaskCmd(voiceUser)
	if err != nil {
		return "", err
	}
	err = voiceWebSocketClient.conn.WriteMessage(websocket.TextMessage, []byte(runTaskCmd))
	return taskID, err
}

// 生成run-task指令

func generateRunTaskCmd(voiceUser VoiceUser) (string, string, error) {
	taskID := uuid.New().String()
	runTaskCmd := Event{
		Header: Header{
			Action:    "run-task",
			TaskID:    taskID,
			Streaming: "duplex",
		},
		Payload: Payload{
			TaskGroup: "audio",
			Task:      "tts",
			Function:  "SpeechSynthesizer",
			Model:     "cosyvoice-v1",
			Parameters: Params{
				TextType:   "PlainText",
				Voice:      voiceUser,
				Format:     "mp3",
				SampleRate: 22050,
				Volume:     50,
				Rate:       1,
				Pitch:      1,
			},
			Input: Input{},
		},
	}
	
	runTaskCmdJSON, err := json.Marshal(runTaskCmd)
	return string(runTaskCmdJSON), taskID, err
}

// 发送待合成文本
// 2. 按顺序发送一个或多个包含待合成文本的continue-task指令

func (voiceWebSocketClient *VoiceWebSocketClient) ContinueTask(taskID, text string) error {
	//texts := []string{"床前明月光", "疑是地上霜", "举头望明月", "低头思故乡"}
	
	runTaskCmd, err := generateContinueTaskCmd(text, taskID)
	if err != nil {
		return err
	}
	
	err = voiceWebSocketClient.conn.WriteMessage(websocket.TextMessage, []byte(runTaskCmd))
	if err != nil {
		return err
	}
	
	return nil
}

// 生成continue-task指令
func generateContinueTaskCmd(text string, taskID string) (string, error) {
	runTaskCmd := Event{
		Header: Header{
			Action:    "continue-task",
			TaskID:    taskID,
			Streaming: "duplex",
		},
		Payload: Payload{
			Input: Input{
				Text: text,
			},
		},
	}
	runTaskCmdJSON, err := json.Marshal(runTaskCmd)
	return string(runTaskCmdJSON), err
}

// 启动一个goroutine来接收结果

func (voiceWebSocketClient *VoiceWebSocketClient) StartResultReceiver(outputFile string) (chan struct{}, *bool) {
	done := make(chan struct{})
	taskStarted := new(bool)
	*taskStarted = false
	
	go func() {
		defer close(done)
		for {
			msgType, message, err := voiceWebSocketClient.conn.ReadMessage()
			if err != nil {
				fmt.Println("解析服务器消息失败：", err)
				return
			}
			
			if msgType == websocket.BinaryMessage {
				// 处理二进制音频流
				fmt.Println("接受到数据返回")
				if err := writeBinaryDataToFile(message, outputFile); err != nil {
					fmt.Println("写入二进制数据失败：", err)
					return
				}
			} else {
				// 处理文本消息
				var event Event
				err = json.Unmarshal(message, &event)
				if err != nil {
					fmt.Println("解析事件失败：", err)
					continue
				}
				if voiceWebSocketClient.handleEvent(event, taskStarted) {
					return
				}
			}
		}
	}()
	
	return done, taskStarted
}

// 处理事件
func (voiceWebSocketClient *VoiceWebSocketClient) handleEvent(event Event, taskStarted *bool) bool {
	switch event.Header.Event {
	case "task-started":
		fmt.Println("收到task-started事件")
		*taskStarted = true
	case "result-generated":
		// 忽略result-generated事件
		return false
	case "task-finished":
		fmt.Println("任务完成")
		return true
	case "task-failed":
		voiceWebSocketClient.handleTaskFailed(event)
		return true
	default:
		fmt.Printf("预料之外的事件：%v\n", event)
	}
	return false
}

// 处理任务失败事件
func (voiceWebSocketClient *VoiceWebSocketClient) handleTaskFailed(event Event) {
	if event.Header.ErrorMessage != "" {
		fmt.Printf("任务失败：%s\n", event.Header.ErrorMessage)
	} else {
		fmt.Println("未知原因导致任务失败")
	}
}

// 关闭连接

func (voiceWebSocketClient *VoiceWebSocketClient) closeConnection() {
	if voiceWebSocketClient.conn != nil {
		voiceWebSocketClient.conn.Close()
	}
}

// 写入二进制数据到文件

func writeBinaryDataToFile(data []byte, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	
	fmt.Println("写入文件成功, ", filePath)
	return nil
}

// 发送finish-task指令
// 3. 发送finish-task指令：结束语音合成任务

func (voiceWebSocketClient *VoiceWebSocketClient) FinishTask(taskID string) error {
	finishTaskCmd, err := generateFinishTaskCmd(taskID)
	if err != nil {
		return err
	}
	err = voiceWebSocketClient.conn.WriteMessage(websocket.TextMessage, []byte(finishTaskCmd))
	return err
}

// 生成finish-task指令
func generateFinishTaskCmd(taskID string) (string, error) {
	finishTaskCmd := Event{
		Header: Header{
			Action:    "finish-task",
			TaskID:    taskID,
			Streaming: "duplex",
		},
		Payload: Payload{
			Input: Input{},
		},
	}
	finishTaskCmdJSON, err := json.Marshal(finishTaskCmd)
	return string(finishTaskCmdJSON), err
}

// 清空输出文件

func clearOutputFile(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	file.Close()
	return nil
}
