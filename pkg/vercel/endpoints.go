package vercel

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func DeploymentsEndpoint(endpoint string, options DeploymentListOpts) (string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}
	u.Path = "/v6/deployments"

	now := time.Now()
	from := now.Add(time.Duration(options.HoursSince) * time.Hour)
	since := fmt.Sprintf("%d", from.Unix()*1000)

	limit := strconv.Itoa(options.Limit)

	q := u.Query()
	q.Add("name", "ved-front")
	q.Add("limit", limit)
	q.Add("projectId", options.ProjectId)
	q.Add("since", since)

	var stateStrings []string
	for _, state := range options.States {
		stateStrings = append(stateStrings, string(state))
	}
	states := strings.Join(stateStrings, ",")
	q.Add("state", states)

	u.RawQuery = q.Encode()
	return u.String(), nil
}

func DeploymentEndpoint(endpoint string, options DeploymentOpts) (string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}
	u.Path = fmt.Sprintf("/v13/deployments/%s", options.ID)
	return u.String(), nil
}

func CancelDeploymentEndpoint(endpoint string, options DeploymentOpts) (string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}
	u.Path = fmt.Sprintf("/v12/deployments/%s/cancel", options.ID)
	return u.String(), nil
}

func CreateDeploymentEndpoint(endpoint string, params RedeploymentParams) (string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}
	u.Path = "/v13/deployments"

	u.Query().Add("teamId", params.TeamID)
	if params.ForceNew {
		u.Query().Add("forceNew", "1")
	}

	return u.String(), nil
}
