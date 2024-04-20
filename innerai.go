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

func decide(p *string, t *Tasks) (int16, error) {
	// tasks := Tasks{Task{Id: "", Description: ""}}
	sysprompt, err := llmhelpergo.TemplateRender(t, routerChooser)
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
		URL:   "https://api.openai.com/v1/chat/completions",
		Model: "gpt-3.5-turbo-0125",
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
