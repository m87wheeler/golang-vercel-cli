package helpers

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/m87wheeler/golang-vercel-cli/internal/environment"
	"github.com/m87wheeler/golang-vercel-cli/internal/screens"
	"github.com/m87wheeler/golang-vercel-cli/pkg/menu"
	"github.com/m87wheeler/golang-vercel-cli/pkg/utils"
	"github.com/m87wheeler/golang-vercel-cli/pkg/vercel"
)

// handleConfig loads and configures the environment, then retrieves Vercel credentials from environment variables.
func HandleConfig(e *environment.Environment) (string, string, string, error) {
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

// renderProjectScreen displays a menu to select a project and returns the selected project name.
func RenderProjectScreen(e *environment.Environment) screens.RenderResult {
	m := menu.NewMenu("Select a project")
	for n := range e.Projects {
		m.AddItem(n, n)
	}
	projectName := m.Display()
	return screens.RenderResult{
		Data: map[string]any{"projectName": projectName},
	}
}

// renderStatesScreen displays a menu to select deployment statuses and returns the selected states.
func RenderStatesScreen(e *environment.Environment) screens.RenderResult {
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

// renderDeploymentsScreen displays a menu to select a deployment and returns the selected deployment ID.
func RenderDeploymentsScreen(v *vercel.VercelAPI, deployments vercel.DeploymentsList) screens.RenderResult {
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

// renderDeploymentScreen displays detailed information about a deployment and returns the deployment data.
func RenderDeploymentScreen(v *vercel.VercelAPI, deploymentId string) screens.RenderResult {
	m := menu.NewMenu("")
	d, err := v.GetDeployment(deploymentId)
	if err != nil {
		return screens.RenderResult{Err: errors.New("No deployment found for " + deploymentId)}
	}
	m.DisplayInfoTable(FormatDeploymentTable(d))
	return screens.RenderResult{
		Data: map[string]any{"deployment": d},
	}
}

// renderDeploymentActionsScreen displays a menu to select a deployment action and returns the selected action.
func RenderDeploymentActionsScreen() screens.RenderResult {
	m := menu.NewMenu("Deployment Actions")
	for k, v := range vercel.DeploymentActionsMap {
		m.AddItem(string(k), v)
	}
	action := m.Display()
	return screens.RenderResult{Data: map[string]any{
		"action": action,
	}}
}

// formatDeploymentTable formats the deployment data into a slice of InfoTableData.
func FormatDeploymentTable(d vercel.DeploymentData) []menu.InfoTableData {
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

// deploymentAction performs the selected action on the specified deployment.
func DeploymentAction(v *vercel.VercelAPI, action, deploymentId string, deployment vercel.DeploymentData) {
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
