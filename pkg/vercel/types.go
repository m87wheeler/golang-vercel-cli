package vercel

import "github.com/m87wheeler/golang-vercel-cli/pkg/http_client"

type VercelAPI struct {
	TeamID     string
	HttpCLient *http_client.HttpClient
	Endpoint   string
	AuthToken  string
	ProjectIDs map[string]string
}

// deployment states
type DeploymentState string

const (
	BUILDING     DeploymentState = "BUILDING"
	ERROR        DeploymentState = "ERROR"
	INITIALIZING DeploymentState = "INITIALIZING"
	QUEUED       DeploymentState = "QUEUED"
	READY        DeploymentState = "READY"
	CANCELED     DeploymentState = "CANCELED"
)

var DeploymentStates = []DeploymentState{
	BUILDING, ERROR, INITIALIZING, QUEUED, READY, CANCELED,
}

// deployment actions
type DeploymentAction string

const (
	EXIT     DeploymentAction = "EXIT"
	CANCEL   DeploymentAction = "CANCEL"
	REDEPLOY DeploymentAction = "REDEPLOY"
)

var DeploymentActionsMap = map[DeploymentAction]string{
	EXIT:     "Exit",
	CANCEL:   "Cancel",
	REDEPLOY: "Redeploy",
}

type DeploymentCreator struct {
	Username string `json:"username"`
}

type DeploymentGitSource struct {
	Branch    string `json:"ref"`
	CommitSHA string `json:"sha"`
}

type DeploymentMeta struct {
	CommitRef string `json:"githubCommitRef"`
}

type DeploymentData struct {
	ID           string              `json:"id"`
	UID          string              `json:"uid"`
	Name         string              `json:"name"`
	Alias        []string            `json:"alias"`
	URL          string              `json:"url"`
	Created      int                 `json:"created"`
	BuildingAt   int                 `json:"buildingAt"`
	Source       string              `json:"source"`
	ReadyState   string              `json:"readyState"`
	Type         string              `json:"type"`
	Creator      DeploymentCreator   `json:"creator"`
	InspectorURL string              `json:"inspectorUrl"`
	GitSource    DeploymentGitSource `json:"gitSource"`
	Meta         DeploymentMeta      `json:"meta"`
}

type DeploymentsList struct {
	Deployments []DeploymentData `json:"deployments"`
}

type DeploymentOpts struct {
	ID string `json:"id"`
}

type DeploymentListOpts struct {
	Limit      int
	ProjectId  string
	HoursSince int
	States     []DeploymentState
}

type RedeploymentParams struct {
	ForceNew bool
	TeamID   string
}

type RedeployProjectSettings struct {
	CommandForIgnoringBuildStep string `json:"commandForIgnoringBuildStep"`
}

type RedeploymentBody struct {
	Name            string                  `json:"name"`
	DeploymentId    string                  `json:"deploymentId"`
	ProjectSettings RedeployProjectSettings `json:"projectSettings"`
}
