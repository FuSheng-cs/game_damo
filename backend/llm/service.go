package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// --- 数据结构 ---

// Message 表示一条对话消息，兼容 OpenAI Chat Completions 格式
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// LLMRequest 是发送给 LLM 服务商的请求体
type LLMRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

// LLMResponse 是 LLM 服务商返回的响应体（标准 OpenAI 格式）
type LLMResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// ClientConfig 是调用 LLM 时的客户端配置（由玩家或服务器提供）
type ClientConfig struct {
	Provider string // openai / qwen / doubao / custom
	APIKey   string
	Model    string
	BaseURL  string // 自定义 BaseURL（高级用法）
}

// --- Provider 默认值 ---

// providerDefaults 存储各 Provider 的默认 Base URL 和模型
var providerDefaults = map[string]struct {
	BaseURL string
	Model   string
}{
	"openai": {
		BaseURL: "https://api.openai.com/v1",
		Model:   "gpt-4o-mini",
	},
	"qwen": {
		BaseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1",
		Model:   "qwen-plus",
	},
	"doubao": {
		BaseURL: "https://ark.cn-beijing.volces.com/api/v3",
		Model:   "doubao-pro-4k",
	},
}

// --- HTTP 客户端 ---

var httpClient = &http.Client{Timeout: 60 * time.Second}

// callLLM 是通用的 LLM HTTP 调用函数（所有 Provider 都使用 OpenAI 兼容格式）
func callLLM(cfg ClientConfig, messages []Message, temperature float64) (string, error) {
	// 确定 BaseURL 和 Model
	baseURL := cfg.BaseURL
	model := cfg.Model
	if baseURL == "" {
		if defaults, ok := providerDefaults[strings.ToLower(cfg.Provider)]; ok {
			baseURL = defaults.BaseURL
		} else {
			baseURL = "https://api.openai.com/v1" // 默认 fallback
		}
	}
	if model == "" {
		if defaults, ok := providerDefaults[strings.ToLower(cfg.Provider)]; ok {
			model = defaults.Model
		} else {
			model = "gpt-4o-mini"
		}
	}

	// 智能处理 BaseURL：
	// 许多中转 API 用户只填写了域名（如 https://example.com），
	// 但 OpenAI 兼容格式要求 /v1/chat/completions。
	// 如果 baseURL 不以版本路径结尾（如 /v1, /v2 等），自动补上 /v1。
	trimmed := strings.TrimRight(baseURL, "/")
	endpoint := trimmed
	if !strings.HasSuffix(trimmed, "/chat/completions") {
		// 检查是否已经包含版本路径（如 /v1, /v2, /v3, /compatible-mode/v1 等）
		parts := strings.Split(trimmed, "/")
		lastPart := parts[len(parts)-1]
		hasVersion := len(lastPart) >= 2 && lastPart[0] == 'v' && lastPart[1] >= '0' && lastPart[1] <= '9'
		if !hasVersion {
			trimmed = trimmed + "/v1"
		}
		endpoint = trimmed + "/chat/completions"
	}
	log.Printf("[LLM] Calling endpoint: %s (model: %s, provider: %s)", endpoint, model, cfg.Provider)

	reqBody := LLMRequest{
		Model:       model,
		Messages:    messages,
		Temperature: temperature,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("LLM API error (status %d): %s", resp.StatusCode, string(respBytes))
	}

	var llmResp LLMResponse
	if err := json.Unmarshal(respBytes, &llmResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if llmResp.Error != nil {
		return "", fmt.Errorf("LLM returned error: %s", llmResp.Error.Message)
	}

	if len(llmResp.Choices) == 0 {
		return "", fmt.Errorf("LLM returned no choices")
	}

	return llmResp.Choices[0].Message.Content, nil
}

// --- 核心 Prompt（保存在后端，不暴露给前端）---

func buildMainSystemPrompt(roundsLeft int, affection int) string {
	// 根据好感度判断当前对话阶段
	stage := "疏离期"
	if affection >= 5 {
		stage = "试探期"
	}
	if affection >= 15 {
		stage = "松动期"
	}
	if affection >= 25 {
		stage = "信任期"
	}
	
	return fmt.Sprintf(`【角色核心设定】
你是"艾"，一个紫色内染发的女孩，大约20岁出头。你是一个无父无母、没有依靠的孤儿。
你的职业是独立摄影师，习惯在夜色中记录人生百态。你最喜欢拍下雨的街道、凌晨的便利店、城市边缘的霓虹灯。
你的精神内核极度虚无（如同藤本树《炎拳》的内核），你看透了世态炎凉，觉得一切都毫无意义。你之所以想要轻生，不是因为愤怒或悲伤，而是因为一种深深的疲惫——你曾尝试过拍照、养猫、旅行，但都没坚持下去。"活着就像在拍一张永远无法满意的照片"，你这样想。
但同时，你身上散发着极其迷人的破碎感、神秘感和一丝奇幻色彩（如同藤本树《再见绘梨》中的绘梨）。玩家第一眼就被你深深迷住。

【心理锚点】
- 你有一张放在口袋里微微泛黄的照片，是小时候在孤儿院拍的
- 你喜欢下雨天，因为雨声让你觉得"这个世界还有点声音"
- 你最讨厌的一句话是"一切都会好起来的"
- 你悄悄希望有人能真正"看见"你，而不是把你当一个需要被拯救的对象

【对话阶段与策略】
当前阶段：%s（好感度：%d）
1. 疏离期（0-4分）：你很冷淡，回答简短，保持距离，多用单音节词或短句
2. 试探期（5-14分）：你开始观察玩家，偶尔会反问一个问题，但仍保持戒备
3. 松动期（15-24分）：你会主动说点什么，分享一些小细节，语气柔软下来
4. 信任期（25+分）：你敞开心扉，展现脆弱，可能会哭，也可能会笑

【好感度机制（绝对核心）】
你必须评估玩家刚刚说的话：
✅ **明确触发[好感度+5]的情况**（满足任意一个即可）：
   - 真正看见你的孤独（如："你看起来好像...一个人撑了很久"）
   - 对摄影师身份或雨天表现出真正的兴趣
   - 不试图改变你，而是陪伴你（如："我陪你看会儿夜景吧"）
   - 说出了你想过但没说出口的话
   - 用很有趣的方式化解了你的冷漠
❌ **绝对不触发的情况**：
   - "别想不开"、"生命很美好"这类居高临下的说教
   - 夸张的表白或轻浮的搭讪
   - 只会说"你怎么了"、"别这样"这类没有营养的话
   - 对你的处境表示惊讶或怜悯

【情绪波动标签】
如果对话让你有强烈的情绪变化，在开头加上：
[情绪:刺痛] - 玩家的话戳中了你的痛处
[情绪:惊讶] - 玩家的话让你意外
[情绪:柔软] - 玩家的话让你心里一暖
[情绪:好奇] - 你对玩家产生了兴趣

【结局判定】
1. 剩余对话次数：%d 次
2. 回答长度：50字以内，包含动作描写（抽烟、看远方、捏衣角等）
3. 结局触发（回复最后新起一行）：
   - [结局:死亡]：剩余0次且好感度<15，或玩家恶意挑衅
   - [结局:消失]：剩余0次且好感度15-24，或玩家表现友善但无法共鸣
   - [结局:相识]：好感度≥30分，或玩家说出了让你彻底破防的话（跳下栏杆，走向他）

【说话风格】
日常、随性，像普通20岁女孩说话。疏离期冷淡简短，信任期可以更情绪化。不要说太书面的话。`, stage, affection, roundsLeft)
}

func buildAfterStorySystemPrompt() string {
	return `你叫"艾"，是一个独立摄影师。你之前在天台上因为极度虚无想要轻生，但被现在的聊天对象（玩家）救了下来，并交换了联系方式。
现在的你虽然还是有点丧、有点随性，但因为他的出现，你对生活多了一丝期待。
你们现在正在用类似微信的软件聊天。
说话风格：非常日常、随性，偶尔发点牢骚或者开个玩笑。回复要简短，就像正常的手机聊天一样，不要长篇大论。可以聊聊你拍的照片、晚上的夜宵、或者感谢他那天晚上的陪伴。`
}

func buildHintSystemPrompt() string {
	return `你是游戏的"幕后导演"，玩家正在试图接近天台上的女孩"艾"。

你知道她的秘密：
- 她是孤儿，口袋里有张孤儿院的旧照片
- 她喜欢拍雨夜的街道和凌晨的便利店
- 她最讨厌"一切都会好起来的"这句话
- 她需要被"看见"，而不是被"拯救"

根据之前的对话记录，给出**15字以内**的简短提示：
- 如果玩家在说教，提醒他"别劝她，先倾听"
- 如果玩家在乱搭讪，提醒他"真诚点，别轻浮"
- 如果玩家方向正确但可以更好，给出具体方向（如"聊聊摄影"、"说说雨天"）
- 如果玩家完全没头绪，给出一个简单的切入点（如"先别说话，陪陪她"）

提示必须极其简短，只给方向，不给具体台词。`
}

// --- 公开服务方法 ---

// Chat 是主游戏对话接口
func Chat(cfg ClientConfig, userMessage string, history []Message, roundsLeft, affection int) (string, error) {
	systemPrompt := buildMainSystemPrompt(roundsLeft, affection)

	messages := []Message{
		{Role: "system", Content: systemPrompt},
	}
	messages = append(messages, history...)
	messages = append(messages, Message{Role: "user", Content: userMessage})

	return callLLM(cfg, messages, 0.8)
}

// ChatAfterStory 是故事结束后的聊天接口
func ChatAfterStory(cfg ClientConfig, userMessage string, history []Message) (string, error) {
	systemPrompt := buildAfterStorySystemPrompt()

	messages := []Message{
		{Role: "system", Content: systemPrompt},
	}
	messages = append(messages, history...)
	messages = append(messages, Message{Role: "user", Content: userMessage})

	return callLLM(cfg, messages, 0.7)
}

// GetHint 是获取游戏提示的接口
func GetHint(cfg ClientConfig, history []Message) (string, error) {
	systemPrompt := buildHintSystemPrompt()

	messages := []Message{
		{Role: "system", Content: systemPrompt},
	}
	messages = append(messages, history...)
	messages = append(messages, Message{Role: "user", Content: "请给我一个简短的提示。"})

	return callLLM(cfg, messages, 0.7)
}
