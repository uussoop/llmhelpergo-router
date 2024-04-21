package llmhelpergorouter

import (
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/uussoop/llmhelpergo"
)

type Task struct {
	Id          string
	Description string
}
type Tasks []Task

func decide(p *string, t *Tasks, extraConfig *map[string]string) (int16, error) {
	// tasks := Tasks{Task{Id: "", Description: ""}}
	url, found := (*extraConfig)["url"]
	if !found {
		url = "https://api.openai.com/v1/chat/completions"
	}
	model, found := (*extraConfig)["model"]
	if !found {
		model = "gpt-3.5-turbo-0125"
	}
	prompt, found := (*extraConfig)["prompt"]
	if !found {
		prompt = routerChooser
	}

	sysprompt, err := llmhelpergo.TemplateRender(t, prompt)
	logrus.Info("here is the decder system prompt: ", sysprompt)
	if err != nil {
		logrus.Error(err)

		return -3, ErrTemplateRender
	}
	llm := &llmhelpergo.GeneralLlm{

		SystemPrompt: sysprompt,
		Messages: &llmhelpergo.Messages{llmhelpergo.Message{
			Content: p,
			Role:    "user",
		}},
		URL:   url,
		Model: model,
	}
	prediction, err := llm.Predict()
	if err != nil {
		logrus.Error(err)

		return -3, ErrMakingPrediction
	}
	predictionint, err := strconv.ParseInt(strings.TrimSpace(*prediction), 10, 16)
	if err != nil {
		logrus.Error(err)

		return -3, ErrMakingConversionToInt
	}
	return int16(predictionint), nil

}
