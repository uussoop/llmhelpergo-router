package llmhelpergorouter

import "errors"

var ErrTemplateRender = errors.New("template render went wrong")
var ErrMakingPrediction = errors.New("predicting the prompt went wrong")
var ErrMakingConversionToInt = errors.New("converting string to int went wrong")
var ErrMakingDecisions = errors.New("making decisions went wrong")
var ErrHandlerMakingPrediction = errors.New("handler making predictions went wrong")
