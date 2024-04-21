package main

import (
	"os"

	"github.com/uussoop/llmhelpergo"
	llmhelpergorouter "github.com/uussoop/llmhelpergo-router"
)

func main() {
	lhr := llmhelpergorouter.New()
	lhr.SetDeciderModel("")
	lhr.SetDeciderPrompt("")
	lhr.SetDeciderUrl("")

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
	group1 := llmhelpergorouter.NewGroup("", false)
	group1.UseRoute("greetings and general purpose questions", chain)
	group1.UseRoute("in matters of manipulating the abstract paragraph of the paper ", chain)
	group1.UseRoute("in matters of manipulating the body paragraph of the paper ", chain)
	group1.UseRoute("in matters of manipulating the q and a paragraph of the paper ", chain)
	group1.UseRoute("in matters of manipulating idea generation ", chain)
	lhr.AddGroup(group1)
	lhr.Run("hi", nil)
}
