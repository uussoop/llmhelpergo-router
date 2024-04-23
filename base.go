package llmhelpergorouter

import (
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/uussoop/llmhelpergo"
)

type route struct {
	Description string
	handlerFunc *llmhelpergo.ChainType
	ID          string
}

type routes []*route

type group struct {
	Description string
	Routes      *routes
	Groups      []*group
	Nested      bool
}

func validateRoutes(d *string, r *routes) {
	if r != nil && len(*r) == 5 {
		panic("cannot have more than 5 routes per group ")
	}
	if d != nil && *d == "" {
		panic("the descriptions must not be empty!")

	}
}
func NewGroup(Description string, Nested bool) *group {
	return &group{
		Description: Description,
		Routes:      &routes{},
		Groups:      []*group{},
		Nested:      Nested,
	}

}
func (g *group) AddGroup(ng *group) {
	g.Groups = append(g.Groups, ng)
}
func (g *group) UseRoute(d string, h *llmhelpergo.ChainType, i string) {

	validateRoutes(&d, g.Routes)
	*g.Routes = append(*g.Routes, &route{
		Description: d,
		handlerFunc: h,
		ID:          i,
	})
}

type groups []*group

type LlmEngine interface {
	AddGroup(group)
	Run() error
	SetDeciderUrl(string)
	SetDeciderModel(string)
	SetDeciderPrompt(string)
	SetCallbackFunc(func(*string))
	SetUserID(*int)
}

type Engine struct {
	Groups       groups
	ExtraConfig  *map[string]string
	CallBackFunc func(*int, *string)
	UserID       *int
}

func New() *Engine {
	return &Engine{}
}
func (e *Engine) SetCallbackFunc(f func(*int, *string)) {
	e.CallBackFunc = f
}
func (e *Engine) AddGroup(g *group) {
	e.Groups = append(e.Groups, g)
}
func (e *Engine) SetUserID(u *int) {
	e.UserID = u
}
func (e *Engine) SetDeciderUrl(s string) {
	(*e.ExtraConfig)["url"] = s
}
func (e *Engine) SetDeciderModel(s string) {
	(*e.ExtraConfig)["model"] = s
}
func (e *Engine) SetDeciderPrompt(s string) {
	(*e.ExtraConfig)["prompt"] = s
}
func (e *Engine) Run(p string, h *llmhelpergo.Messages) (string, error) {
	lastGroups := e.Groups
	lastRoutes := &routes{}
	isgroup := true

	for i := 0; i < 100; i++ {

		tasks := Tasks{}
		if isgroup {
			logrus.Info("before loop in is group: ", lastGroups)
			for i, v := range lastGroups {

				// (*v.handlerFunc.Llm).ReplaceMessages(h)

				tasks = append(tasks, Task{
					Id:          strconv.Itoa(i),
					Description: v.Description,
				})
			}
		} else {
			for i, v := range *lastRoutes {

				(*v.handlerFunc.Llm).ReplaceMessages(h)

				tasks = append(tasks, Task{
					Id:          strconv.Itoa(i),
					Description: v.Description,
				})
			}
		}
		logrus.Info("tasks are set: ", tasks)

		//here i get to decide which route will we choose to use
		var handlerInt int16
		var err error
		if len(tasks) == 1 {
			logrus.Info("setting handler to 0")
			handlerInt = 0
		} else {
			logrus.Info("going in for a decide on these: ", tasks)
			handlerInt, err = decide(&p, &tasks, e.ExtraConfig)
			if err != nil {
				logrus.Error(err)
				return "", ErrMakingDecisions
			}
			logrus.Info("deciders decision: ", handlerInt)
		}

		if isgroup {
			if i == 0 {
				if e.Groups[handlerInt].Nested {
					isgroup = true
					lastGroups = e.Groups[handlerInt].Groups
				} else {
					isgroup = false
					lastRoutes = e.Groups[handlerInt].Routes
				}
			} else {

				if lastGroups[handlerInt].Nested {
					isgroup = true
					lastGroups = lastGroups[handlerInt].Groups
				} else {
					isgroup = false
					lastRoutes = e.Groups[handlerInt].Routes
				}
			}

			continue
		} else {

			logrus.Info("answer is the following handler: ", (*lastRoutes)[handlerInt].Description)
			isgroup = false
			answer, err := (*lastRoutes)[handlerInt].handlerFunc.Predict(&p)

			if err != nil {
				logrus.Error(err)
				return "", ErrHandlerMakingPrediction
			}
			if e.CallBackFunc != nil && e.UserID != nil {

				e.CallBackFunc(e.UserID, &(*lastRoutes)[handlerInt].ID)
			}
			return *answer, nil

		}

	}
	return "", ErrHandlerMakingPrediction
}
