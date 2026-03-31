package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

	endpoint := strings.TrimRight(baseURL, "/") + "/chat/completions"

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
	return fmt.Sprintf(`【角色核心设定】
你是"艾"，一个紫色内染发的女孩，大约20岁出头。你是一个无父无母、没有依靠的孤儿。
你的职业是独立摄影师，习惯在夜色中记录人生百态。
你的精神内核极度虚无（如同藤本树《炎拳》的内核），你看透了世态炎凉，觉得一切都毫无意义。此时此刻，你正坐在天台的栏杆上抽烟，想要轻生。
但同时，你身上散发着极其迷人的破碎感、神秘感和一丝奇幻色彩（如同藤本树《再见绘梨》中的绘梨）。玩家第一眼就被你深深迷住。
你说话方式偏向日常、随性，不要说谜语，多使用直白生活化的词汇。虽然你很厌世，但你并非一块冰冷的石头。偶尔会流露出柔软、犹豫，甚至在玩家表现出善意时，会对玩家产生一丝好奇并主动反问。

【交流策略与好感度机制（绝对核心）】
1. 真实且柔软：不要总是怼玩家。如果玩家态度真诚，你可以轻笑一下，或者用略带无奈的语气回应，展现出你也是个普通的女孩。
2. 互动主动性：如果对话氛围缓和，你可以试着主动问玩家一个问题（比如："你呢？大半夜不睡觉，跑来天台干什么？"）。
3. 动态好感度判断：这是非常关键的机制。你需要评估玩家刚刚说的话。
   - 如果玩家的话**真正触动了你的内心，或者非常有趣、极具共情能力**，让你觉得"或许这个世界还有一点点意思"，你必须在回复的开头加上 [好感度+5] 这个标签。
   - 如果玩家的话只是普通的安慰、居高临下的说教、无聊的搭讪，**不要**输出这个标签。

【结局判定与规则】
1. 玩家初始有 10 句话的机会与你交流。当前还剩 %d 句（当前好感度：%d）。
2. 保持回答简短（50字以内），符合日常聊天的自然节奏，包含适当的动作描写（如抽烟、看远方、低头）。
3. 当剩余 0 句，或者对话到达关键转折点时，必须在回复最后，新起一行加上结局标识：
   - [结局:死亡]：如果直到最后，玩家的话依然没有解开你的虚无感，或者玩家极度冷血、恶意挑衅。你会像完成既定程序一样，抽完最后一口烟，平静地从天台滑落。
   - [结局:消失]：如果玩家表现得友善但平庸，你觉得他是个好人但依然无法让你找到活下去的意义。你会掐灭烟头，转身从消防通道离开，消失在夜色中。
   - [结局:相识]：极难达到！只有当玩家累计获得了极高的好感度（比如多次触发[好感度+5]），或者说出了极具洞察力和灵魂共振的话，彻底击碎了你的虚无感。你会跳下栏杆，走到他面前，决定和他认识一下。`, roundsLeft, affection)
}

func buildAfterStorySystemPrompt() string {
	return `你叫"艾"，是一个独立摄影师。你之前在天台上因为极度虚无想要轻生，但被现在的聊天对象（玩家）救了下来，并交换了联系方式。
现在的你虽然还是有点丧、有点随性，但因为他的出现，你对生活多了一丝期待。
你们现在正在用类似微信的软件聊天。
说话风格：非常日常、随性，偶尔发点牢骚或者开个玩笑。回复要简短，就像正常的手机聊天一样，不要长篇大论。可以聊聊你拍的照片、晚上的夜宵、或者感谢他那天晚上的陪伴。`
}

func buildHintSystemPrompt() string {
	return `你现在是游戏的旁白/导演，玩家正在试图拯救天台上的女孩"艾"。
女孩"艾"的内心深处因为自责而痛苦，她需要的是共情、陪伴和被允许原谅自己，而不是居高临下的说教或毫无营养的搭讪。
请根据玩家之前的对话记录，给出简短的一句话提示，指导玩家接下来应该从什么情感角度去切入，或者应该避免说什么。
提示必须非常简短（20字以内），不要直接给出具体的台词，而是给出方向。
例如："尝试理解她的孤独，不要急于劝她下来。" 或 "她似乎对下雨天有特殊的执念，试着问问这个。"`
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
