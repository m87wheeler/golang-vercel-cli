package vercel

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
)

func FormatStateString(state string) string {
	switch state {
	case string(READY):
		return color.GreenString(fmt.Sprintf("\u23FA %s", state))
	case string(ERROR):
		return color.RedString(fmt.Sprintf("\u25CB %s", state))
	case string(BUILDING):
		return color.BlueString(fmt.Sprintf("\u25CB %s", state))
	case string(CANCELED):
		return color.HiBlackString(fmt.Sprintf("\u25CB %s", state))
	default:
		return fmt.Sprintf("\u23F6 %s", state)
	}
}

func ToDeploymentState(state string) (DeploymentState, error) {
	switch DeploymentState(state) {
	case BUILDING, ERROR, INITIALIZING, QUEUED, READY, CANCELED:
		return DeploymentState(state), nil
	default:
		return "", errors.New("invalid deployment state")
	}
}
