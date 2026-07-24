// cmd/claw/main.go
package main

import (
	"github.com/SparkChase/go-tiny-claw/internal/engine"
	"github.com/SparkChase/go-tiny-claw/internal/feishu"
	"github.com/SparkChase/go-tiny-claw/internal/provider"
	"github.com/SparkChase/go-tiny-claw/internal/tools"
	"github.com/larksuite/oapi-sdk-go/v3/core/httpserverext"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	// 自动加载 .env 文件
	_ = godotenv.Load()

	// 确保设置了 ZHIPU_API_KEY
	if os.Getenv("ZHIPU_API_KEY") == "" {
		log.Fatal("请先导出 ZHIPU_API_KEY 环境变量")
	}

	workDir, _ := os.Getwd()
	model := os.Getenv("MODEL")
	if model == "" {
		model = "glm-4.5-air" // fallback
	}
	llmProvider := provider.NewZhipuOpenAIProvider(model)
	registry := tools.NewRegistry()

	registry.Register(tools.NewReadFileTool(workDir))
	registry.Register(tools.NewWriteFileTool(workDir))
	registry.Register(tools.NewBashTool(workDir))
	registry.Register(tools.NewEditFileTool(workDir))

	// 开启慢思考
	eng := engine.NewAgentEngine(llmProvider, registry, workDir, true)

	// 2. 初始化飞书 Bot 调度器
	bot := feishu.NewFeishuBot(eng)
	handler := httpserverext.NewEventHandlerFunc(bot.GetEventDispatcher())

	// 3. 注册路由并启动 HTTP 服务
	http.HandleFunc("/webhook/event", handler)

	port := ":48080"
	log.Printf("🚀 go-tiny-claw 飞书服务端已启动，正在监听 %s 端口\n", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
