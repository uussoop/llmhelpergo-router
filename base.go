package llmhelpergorouter

import (
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/uussoop/llmhelpergo"
)

type route struct {
	Description string
	handlerFunc *llmhelpergo.ChainType
}

type routes []*route

type group struct {
	Routes *routes
	Groups []*group
	Nested bool
}

func validateRoutes(d *string, r *routes) {
	if len(*r) == 5 {
		panic("cannot have more than 5 routes per group ")
	}
	if *d == "" {
		panic("the descriptions must not be empty!")

	}
}
func (g *group) AddGroup(ng *group) {
	g.Groups = append(g.Groups, ng)
}
func (g *group) UseRoute(d string, h *llmhelpergo.ChainType) {
	validateRoutes(&d, g.Routes)
	*g.Routes = append(*g.Routes, &route{
		Description: d,
		handlerFunc: h,
	})
}

type groups []*group

type LlmEngine interface {
	UseRoute(string, llmhelpergo.ChainType)
	AddGroup(group)
	Run() error
}

type Engine struct {
	Routes      routes
	Groups      groups
	ExtraConfig map[string]string
}

func New() *Engine {
	return &Engine{}
}

func (e *Engine) AddGroup(g *group) {
	e.Groups = append(e.Groups, g)
}

func (e *Engine) UseRoute(d string, h *llmhelpergo.ChainType) {
	validateRoutes(&d, &e.Routes)

	e.Routes = append(e.Routes, &route{
		Description: d,
		handlerFunc: h,
	})
}

func (e *Engine) Run(p string, h *llmhelpergo.Messages) (string, error) {
	tasks := Tasks{}
	for i, v := range e.Routes {

		(*v.handlerFunc.Llm).ReplaceMessages(h)
		tasks = append(tasks, Task{
			Id:          strconv.Itoa(i),
			Description: v.Description,
		})
	}
	//here i get to decide which route will we choose to use
	handlerInt, err := decide(&p, &tasks)
	if err != nil {
		logrus.Error(err)
		return "", ErrMakingDecisions
	}
	answer, err := e.Routes[handlerInt].handlerFunc.Predict(&p)
	if err != nil {
		logrus.Error(err)
		return "", ErrHandlerMakingPrediction
	}

	return *answer, nil
}
