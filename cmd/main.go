package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/m87wheeler/golang-vercel-cli/internal/environment"
	"github.com/m87wheeler/golang-vercel-cli/internal/screens"
	"github.com/m87wheeler/golang-vercel-cli/pkg/http_client"
	"github.com/m87wheeler/golang-vercel-cli/pkg/menu"
	"github.com/m87wheeler/golang-vercel-cli/pkg/utils"
	"github.com/m87wheeler/golang-vercel-cli/pkg/vercel"

	"github.com/joho/godotenv"
)

var version = "development"

func main() {
	fmt.Printf("Version: %s\n", version)

	// Load or configure environment
	e := environment.NewEnvironment()
	vercelEndpoint, vercelAuthKey, vercelTeamID, err := handleConfig(e)
	if err != nil {
		log.Fatal(err)
	}

	// Define Screens
	scr := screens.ScreensList{
		Project:     renderProjectScreen,
		States:      renderStatesScreen,
		Deployments: renderDeploymentsScreen,
		Deployment:  renderDeploymentScreen,
		Actions:     renderDeploymentActionsScreen,
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

	deploymentAction(v, action, deploymentId, deployment)
}

// handleConfig loads and configures the environment, then retrieves Vercel credentials from environment variables.
func handleConfig(e *environment.Environment) (string, string, string, error) {
	e.Load()

	if !e.EnvFound {
		log.Default().Println("No env file found")
		e.Configure() // exits from app
	}

	// Load environment variables
	fmt.Println("Loading env config from " + e.EnvLoadFrom)
	err := godotenv.Load(e.EnvLoadFrom)
	if err != nil {
		fmt.Println("Error loading .env file")
		return "", "", "", err
	}

	vercelEndpoint := os.Getenv("VERCEL_ENDPOINT")
	vercelAuthKey := os.Getenv("VERCEL_AUTH_KEY")
	vercelTeamID := os.Getenv("VERCEL_TEAM_ID")

	if vercelEndpoint == "" || vercelAuthKey == "" || vercelTeamID == "" {
		fmt.Println("Missing credentials")
		return "", "", "", err
	}

	return vercelEndpoint, vercelAuthKey, vercelTeamID, nil
}

func renderProjectScreen(e *environment.Environment) screens.RenderResult {
	m := menu.NewMenu("Select a project")
	for n := range e.Projects {
		m.AddItem(n, n)
	}
	projectName := m.Display()
	return screens.RenderResult{
		Data: map[string]any{"projectName": projectName},
	}
}

func renderStatesScreen(e *environment.Environment) screens.RenderResult {
	m := menu.NewMenu("Select deployment status'")
	for _, s := range vercel.DeploymentStates {
		ss := string(s)
		m.AddItem(ss, ss)
	}
	states := []string{string(vercel.READY), string(vercel.BUILDING)}
	m.DisplayMultiChoice(func(choice string) []string {
		states = utils.ToggleState(states, choice)
		return states
	})

	if len(states) < 1 {
		fmt.Println("Must choose at least 1 state")
		return screens.RenderResult{Err: errors.New("Must choose at least 1 state")}
	}

	return screens.RenderResult{
		Data: map[string]any{"states": states},
	}
}

func renderDeploymentsScreen(v *vercel.VercelAPI, deployments vercel.DeploymentsList) screens.RenderResult {
	m := menu.NewMenu("Select a deployment")
	for _, d := range deployments.Deployments {
		elapsed := utils.ElapsedTime(int64(d.Created) / 1000)
		m.AddItem(d.UID, fmt.Sprintf("%-20s\t%-25s\t%-10s\t%-10s\t%-10s", d.Name, d.Creator.Username, d.Meta.CommitRef, elapsed, d.ReadyState))
	}
	deploymentId := m.Display()
	return screens.RenderResult{
		Data: map[string]any{"deploymentId": deploymentId},
	}
}

func renderDeploymentScreen(v *vercel.VercelAPI, deploymentId string) screens.RenderResult {
	m := menu.NewMenu("")
	d, err := v.GetDeployment(deploymentId)
	if err != nil {
		return screens.RenderResult{Err: errors.New("No deployment found for " + deploymentId)}
	}
	m.DisplayInfoTable(formatDeploymentTable(d))
	return screens.RenderResult{
		Data: map[string]any{"deployment": d},
	}
}

func renderDeploymentActionsScreen() screens.RenderResult {
	m := menu.NewMenu("Deployment Actions")
	for k, v := range vercel.DeploymentActionsMap {
		m.AddItem(string(k), v)
	}
	action := m.Display()
	return screens.RenderResult{Data: map[string]any{
		"action": action,
	}}
}

// Formats the deployment data into a slice of InfoTableData.
func formatDeploymentTable(d vercel.DeploymentData) []menu.InfoTableData {
	var url string
	if len(d.Alias) > 0 {
		url = d.Alias[0]
	}
	elapsed := utils.ElapsedTime(int64(d.BuildingAt) / 1000)

	return []menu.InfoTableData{
		{Label: "ID", Data: d.ID},
		{Label: "Name", Data: d.Name},
		{Label: "Creator", Data: d.Creator.Username},
		{Label: "State", Data: vercel.FormatStateString(d.ReadyState)},
		{Label: "Started", Data: elapsed},
		{Label: "URL", Data: url},
		{Label: "Branch", Data: d.GitSource.Branch},
		{Label: "Commit SHA", Data: d.GitSource.CommitSHA},
	}
}

func deploymentAction(v *vercel.VercelAPI, action, deploymentId string, deployment vercel.DeploymentData) {
	switch action {
	case string(vercel.CANCEL):
		_, err := v.CancelDeployment(deploymentId)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	case string(vercel.REDEPLOY):
		_, err := v.CreateRedeployment(deployment)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	case string(vercel.EXIT):
	default:
		os.Exit(0)
	}
}
