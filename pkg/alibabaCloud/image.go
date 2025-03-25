package alibabaCloud

type ImageModel string

// 文生图模型需要使用一段文字描述生成的图片。提示词（prompt）描述越完整、精确和丰富，生成的图像品质越高，越贴近期望生成的内容。

const (
	TextToImageV21Plus  ImageModel = "wanx2.1-t2i-plus"          // 全面升级版本。生成图像细节更丰富，速度稍慢。对应通义万相官网2.1专业模型。
	TextToImageV21Turbo ImageModel = "wanx2.1-t2i-turbo"         // 全面升级版本。生成速度快、效果全面、综合性价比高。对应通义万相官网2.1极速模型。
	TextToImageV20Turbo ImageModel = "wanx2.1-t2i-turbo"         // 擅长质感人像，速度中等、成本较低。对应通义万相官网2.0极速模型。
	TextToPoster        ImageModel = "wanx-poster-generation-v1" // 创意海报生成，您的创意海报魔法工厂！
	ImageRepaint        ImageModel = "wanx-style-repaint-v1"     // 通义万相-人像风格重绘可以将输入的人物图像进行多种风格化的重绘生成，使新生成的图像在兼顾原始人物相貌的同时，带来不同风格的绘画效果。
	CosplayImage        ImageModel = "wanx-style-cosplay-v1"     // 通义万相-Cosplay动漫人物生成。
)

const (
	getTaskUrl         = "https://dashscope.aliyuncs.com/api/v1/tasks/"                                    // 获取任务API
	generatePosterUrl  = "https://dashscope.aliyuncs.com/api/v1/services/aigc/text2image/image-synthesis"  // 生成海报API
	generateRepaintUrl = "https://dashscope.aliyuncs.com/api/v1/services/aigc/image-generation/generation" // 人像风格重绘API
	generateCosplayUrl = "https://dashscope.aliyuncs.com/api/v1/services/aigc/image-generation/generation" // 生成Cosplay图片
)

/*
prompt example:

prompt = "生成一张新年祝福贺卡，背景有白雪，放鞭炮的小孩，蛇形成文案2025，并写上HAPPY NEW YEAR。"
prompt = "一个用针毡制作的圣诞老人手持礼物，旁边站着一只白色的猫咪，背景中有许多五颜六色的礼物。整个场景应该是可爱、温暖和舒适的，并且背景中还有一些绿色植物。"
prompt = "中国女孩，圆脸，看着镜头，优雅的民族服装，商业摄影，室外，电影级光照，半身特写，精致的淡妆，锐利的边缘。 "
prompt = "高清摄影写真，一只布偶猫慵懒地躺在复古木质窗台上。它拥有一身柔软的银白色长毛，深蓝色宝石般的眼睛，粉嫩的小鼻头和肉垫。猫咪眼神温柔地望向镜头，嘴角似乎带着一抹满足的微笑。午后阳光透过半开的窗户洒在它身上，营造出温馨而宁静的氛围。背景是模糊的绿色植物和老式窗帘，增添了几分生活气息。近景特写，自然光影效果。"
*/

// 文本生成图片

func TextGenerateImage() {
	
}

// 创意海报生成

type TextGeneratePosterReq struct {
	model string                     `json:"model,omitempty"`
	input TextGeneratePosterReqInput `json:"input,omitempty"`
}

type TextGeneratePosterReqInput struct {
	Title        string  `json:"title,omitempty"`
	SubTitle     string  `json:"sub_title,omitempty"`
	BodyText     string  `json:"body_text,omitempty"`
	PromptTextZh string  `json:"prompt_text_zh,omitempty"`
	WhRatios     string  `json:"wh_ratios,omitempty"`
	LoraName     string  `json:"lora_name,omitempty"`
	LoraWeight   float32 `json:"lora_weight,omitempty"`
	CtrlRatio    float32 `json:"ctrl_ratio,omitempty"`
	CtrlStep     float32 `json:"ctrl_step,omitempty"`
	GenerateMode string  `json:"generate_mode,omitempty"`
	GenerateNum  int     `json:"generate_num,omitempty"`
}

func TextGeneratePoster() {
	
}

// 生成Cosplay图片，需要一张脸部图片 & 一张对应的动漫图片

type GenerateCosplayImageReq struct {
	Model string                       `json:"model,omitempty"`
	Input GenerateCosplayImageReqInput `json:"input,omitempty"`
}

type GenerateCosplayImageReqInput struct {
	FaceImageUrl     string `json:"face_image_url,omitempty"`
	TemplateImageUrl string `json:"template_image_url,omitempty"`
}

func GenerateCosplayImage() {
	
}
