package main

import (
	"os"

	"github.com/uussoop/llmhelpergo"
	llmhelpergorouter "github.com/uussoop/llmhelpergo-router"
)

func main() {
	lhr := llmhelpergorouter.New()
	os.Setenv("OPENAI_KEY", "")
	// message := "hello world"
	llm := &llmhelpergo.GeneralLlm{
		SystemPrompt: "you are an ai",
		Messages:     nil,
		URL:          "https://api.openai.com/v1/chat/completions",
		Model:        "gpt-4",
	}
	chain := llmhelpergo.Chain(llm, 450)
	chain.Use(llmhelpergo.SampleAgent)
	chain.Use(llmhelpergo.SampleAgent2)
	chain.Use(llmhelpergo.SampleAgent3)

	lhr.UseRoute("", chain)
	lhr.Run("", nil)
}
