package main

import (
	"fmt"
	"log"
	"os"

	"github.com/m87wheeler/golang-vercel-cli/internal/environment"
	"github.com/m87wheeler/golang-vercel-cli/internal/helpers"
	"github.com/m87wheeler/golang-vercel-cli/internal/screens"
	"github.com/m87wheeler/golang-vercel-cli/pkg/http_client"
	"github.com/m87wheeler/golang-vercel-cli/pkg/vercel"
)

var version = "development"

// main is the entry point of the application.
func main() {
	fmt.Printf("Version: %s\n", version)

	// Load or configure environment
	e := environment.NewEnvironment()
	vercelEndpoint, vercelAuthKey, vercelTeamID, err := helpers.HandleConfig(e)
	if err != nil {
		log.Fatal(err)
	}

	// Define Screens
	scr := screens.ScreensList{
		Project:     helpers.RenderProjectScreen,
		States:      helpers.RenderStatesScreen,
		Deployments: helpers.RenderDeploymentsScreen,
		Deployment:  helpers.RenderDeploymentScreen,
		Actions:     helpers.RenderDeploymentActionsScreen,
	}

	// Project Name Menu
	sc := scr.Project(e)
	if sc.Err != nil {
		log.Fatal(sc.Err)
	}
	projectName, ok := sc.Data["projectName"].(string)
	if !ok {
		log.Fatal("missing project name")
	}

	// Status Multi-choice Menu
	sc = scr.States(e)
	if sc.Err != nil {
		log.Fatal(sc.Err)
	}
	states, ok := sc.Data["states"].([]string)
	if !ok {
		log.Fatal("missing states")
	}

	// Fetch Deployments
	c := http_client.NewHttpClient()
	v := vercel.NewVercelAPI(c, vercelEndpoint, vercelAuthKey, vercelTeamID, e.Projects)
	dl, err := v.GetDeployments(projectName, 10, 24, states)
	if err != nil {
		log.Panic(err)
	} else if len(dl.Deployments) < 1 {
		fmt.Println("No deployments to display")
		os.Exit(0)
	}

	// Deployment Menu
	sc = scr.Deployments(v, dl)
	if sc.Err != nil {
		log.Fatal(sc.Err)
	}
	deploymentId, ok := sc.Data["deploymentId"].(string)
	if !ok {
		log.Fatal("no deployment id")
	}

	// Deployment Data
	sc = scr.Deployment(v, deploymentId)
	if sc.Err != nil {
		log.Fatal(sc.Err)
	}
	deployment, ok := sc.Data["deployment"].(vercel.DeploymentData)
	if !ok {
		log.Fatal("no deployment")
	}

	// Deployment Actions
	sc = scr.Actions()
	if sc.Err != nil {
		log.Fatal(sc.Err)
	}
	action, ok := sc.Data["action"].(string)
	if !ok {
		log.Fatal("no deployment id")
	}

	helpers.DeploymentAction(v, action, deploymentId, deployment)
}
