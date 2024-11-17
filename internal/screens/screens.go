package screens

import (
	"github.com/m87wheeler/golang-vercel-cli/internal/environment"
	"github.com/m87wheeler/golang-vercel-cli/pkg/vercel"
)

type RenderResult struct {
	Data map[string]any
	Err  error
}

type RenderArgs struct {
	Env *environment.Environment
}

type RenderDeploymentsArgs struct {
	VercelAPI       *vercel.VercelAPI
	DeploymentsList vercel.DeploymentsList
}

type RenderDeploymentArgs struct {
	VercelAPI    *vercel.VercelAPI
	DeploymentID string
}

type Screen[T any] struct {
	Render func(args T) RenderResult
}

type ScreensList struct {
	Project     func(e *environment.Environment) RenderResult
	States      func(e *environment.Environment) RenderResult
	Deployments func(v *vercel.VercelAPI, d vercel.DeploymentsList) RenderResult
	Deployment  func(v *vercel.VercelAPI, id string) RenderResult
	Actions     func() RenderResult
}

type Screens struct {
	List          ScreensList
	CurrentScreen string
}

func NewScreens(screens ScreensList, initialScreen string) *Screens {
	return &Screens{
		List:          screens,
		CurrentScreen: initialScreen,
	}
}
