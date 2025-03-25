package tools

import (
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
	"io"
	"time"
)

func PlayVoice(content io.ReadCloser) error {
	// 解码 MP3 文件
	streamer, format, err := mp3.Decode(content)
	defer streamer.Close()
	
	if err != nil {
		return err
	}
	
	// 初始化音频播放器
	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		return err
	}
	
	// 播放音频
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	
	// 等待播放完成
	<-done
	return nil
}
