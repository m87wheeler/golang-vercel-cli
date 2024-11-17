package main

import (
	"fmt"
	"log"
	"os"

	"github.com/m87wheeler/golang-vercel-cli/internal/environment"
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
	e.Load()

	if !e.EnvFound {
		log.Default().Println("No env file found")
		e.Configure() // exits from app
	}

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")

	}

	vercelEndpoint := os.Getenv("VERCEL_ENDPOINT")
	vercelAuthKey := os.Getenv("VERCEL_AUTH_KEY")
	vercelTeamID := os.Getenv("VERCEL_TEAM_ID")

	if vercelEndpoint == "" || vercelAuthKey == "" || vercelTeamID == "" {
		log.Fatal("Missing credentials")
	}

	// Project Name Menu
	m := menu.NewMenu("Select a project")
	for n := range e.Projects {
		m.AddItem(n, n)
	}
	projectName := m.Display()

	// Status Multi-choice Menu
	m = menu.NewMenu("Select deployment status'")
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
		os.Exit(0)
	}

	// Fetch Deployments
	c := http_client.NewHttpClient()
	v := vercel.NewVercelAPI(c, vercelEndpoint, vercelAuthKey, vercelTeamID, e.Projects)
	ds, err := v.GetDeployments(projectName, 10, 24, states)

	if err != nil {
		log.Panic(err)
	} else if len(ds.Deployments) < 1 {
		fmt.Println("No deployments to display")
		os.Exit(0)
	}

	// Deployment Menu
	m = menu.NewMenu("Select a deployment")
	for _, d := range ds.Deployments {
		elapsed := utils.ElapsedTime(int64(d.Created) / 1000)
		m.AddItem(d.UID, fmt.Sprintf("%-20s\t%-25s\t%-10s\t%-10s\t%-10s", d.Name, d.Creator.Username, d.Meta.CommitRef, elapsed, d.ReadyState))
	}
	deploymentId := m.Display()

	// Deployment Data
	m = menu.NewMenu("")
	d, err := v.GetDeployment(deploymentId)
	if err != nil {
		log.Fatal("No deployment found for " + deploymentId)
	}
	m.DisplayInfoTable(formatDeploymentTable(d))

	// Deployment Actions
	m = menu.NewMenu("Deployment Actions")
	for k, v := range vercel.DeploymentActionsMap {
		m.AddItem(string(k), v)
	}
	action := m.Display()

	switch action {
	case string(vercel.CANCEL):
		_, err := v.CancelDeployment(deploymentId)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	case string(vercel.REDEPLOY):
		_, err := v.CreateRedeployment(d)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	case string(vercel.EXIT):
	default:
		os.Exit(0)
	}
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
