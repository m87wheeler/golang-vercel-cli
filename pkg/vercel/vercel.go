package vercel

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/m87wheeler/golang-vercel-cli/pkg/http_client"
)

func NewVercelAPI(httpClient *http_client.HttpClient, endpoint, authToken, teamId string, projectIds map[string]string) *VercelAPI {
	return &VercelAPI{
		HttpCLient: httpClient,
		Endpoint:   endpoint,
		AuthToken:  authToken,
		ProjectIDs: projectIds,
		TeamID:     teamId,
	}
}

func (v *VercelAPI) GetDeployments(projectName string, limit, hoursSince int, states []string) (DeploymentsList, error) {
	var deployments DeploymentsList
	projectId, ok := v.ProjectIDs[projectName]
	if !ok {
		return deployments, errors.New("No project ID found for " + projectName)
	}

	// ensure hoursSince is a negative number
	if hoursSince > 0 {
		hoursSince = hoursSince * -1
	}

	// ensure hoursSince is at least -1
	if hoursSince > -1 {
		hoursSince = -1
	}

	var st []DeploymentState
	for _, s := range states {
		ds, err := ToDeploymentState(s)
		if err == nil {
			fmt.Println(err)
			st = append(st, ds)
		}
	}

	url, err := DeploymentsEndpoint(v.Endpoint, DeploymentListOpts{
		Limit:      limit,
		ProjectId:  projectId,
		HoursSince: -24,
		States:     st,
		// States:     []DeploymentState{BUILDING, READY, CANCELED, ERROR},
	})
	if err != nil {
		fmt.Println(err)
		return deployments, err
	}

	resp, err := http.NewRequest(http.MethodGet, url, nil)
	resp.Header.Add("Authorization", fmt.Sprintf("Bearer %s", v.AuthToken))

	response, err := http.DefaultClient.Do(resp)
	if err != nil {
		fmt.Println(err)
		return deployments, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return deployments, err
	}

	err = json.Unmarshal(body, &deployments)
	if err != nil {
		fmt.Println(err)
		return deployments, err
	}

	return deployments, nil
}

func (v *VercelAPI) GetDeployment(deploymentId string) (DeploymentData, error) {
	var deployment DeploymentData

	url, err := DeploymentEndpoint(v.Endpoint, DeploymentOpts{
		ID: deploymentId,
	})
	if err != nil {
		fmt.Println(err)
		return deployment, err
	}

	resp, err := http.NewRequest(http.MethodGet, url, nil)
	resp.Header.Add("Authorization", fmt.Sprintf("Bearer %s", v.AuthToken))

	response, err := http.DefaultClient.Do(resp)
	if err != nil {
		fmt.Println(err)
		return deployment, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return deployment, err
	}

	err = json.Unmarshal(body, &deployment)
	if err != nil {
		fmt.Println(err)
		return deployment, err
	}

	return deployment, nil
}

func (v *VercelAPI) CancelDeployment(deploymentId string) (DeploymentData, error) {
	var deployment DeploymentData

	url, err := CancelDeploymentEndpoint(v.Endpoint, DeploymentOpts{
		ID: deploymentId,
	})
	if err != nil {
		fmt.Println(err)
		return deployment, err
	}

	resp, err := http.NewRequest(http.MethodPatch, url, nil)
	resp.Header.Add("Authorization", fmt.Sprintf("Bearer %s", v.AuthToken))

	response, err := http.DefaultClient.Do(resp)
	if err != nil {
		fmt.Println(err)
		return deployment, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return deployment, err
	}

	err = json.Unmarshal(body, &deployment)
	if err != nil {
		fmt.Println(err)
		return deployment, err
	}

	fmt.Printf("Cancelling %s", deployment.ID)
	return deployment, nil
}

func (v *VercelAPI) CreateRedeployment(sourceDeployment DeploymentData) (DeploymentData, error) {
	var deployment DeploymentData

	url, err := CreateDeploymentEndpoint(
		v.Endpoint,
		RedeploymentParams{
			ForceNew: true,
			TeamID:   v.TeamID,
		},
	)
	if err != nil {
		fmt.Println(err)
		return deployment, err
	}
	fmt.Println(url)

	bodyData := RedeploymentBody{
		Name:         sourceDeployment.Name,
		DeploymentId: sourceDeployment.ID,
		ProjectSettings: RedeployProjectSettings{
			CommandForIgnoringBuildStep: "",
		},
	}
	jsonBody, err := json.Marshal(bodyData)
	if err != nil {
		fmt.Println(err)
		return deployment, err
	}

	resp, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	resp.Header.Add("Content-Type", "application/json")
	resp.Header.Add("Authorization", fmt.Sprintf("Bearer %s", v.AuthToken))

	response, err := http.DefaultClient.Do(resp)
	if err != nil {
		fmt.Println(err)
		return deployment, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return deployment, err
	}

	err = json.Unmarshal(body, &deployment)
	if err != nil {
		fmt.Println(err)
		return deployment, err
	}

	fmt.Printf("Redeploying %s", deployment.ID)
	return deployment, nil
}
