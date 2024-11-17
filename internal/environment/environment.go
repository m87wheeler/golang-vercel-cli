package environment

import (
	"fmt"
	"log"
	"os"

	"github.com/m87wheeler/golang-vercel-cli/pkg/utils"
)

type Environment struct {
	// env file
	EnvFileName string
	EnvFound    bool
	// vercel values
	TeamID      string
	AuthKey     string
	ApiEndpoint string
	Projects    map[string]string
}

func NewEnvironment() *Environment {
	return &Environment{
		EnvFileName: ".env",
		EnvFound:    false,
		TeamID:      "",
		AuthKey:     "",
		ApiEndpoint: "https://api.vercel.com",
		Projects: map[string]string{ // TODO create in CLI
			"ved-front":            "prj_UgdzEpZRLCD7JwXVU4WsNFuFwKwi",
			"ved-front-split-test": "prj_FUa8su0HhrW31lJbNBExwFEFAiol",
		},
	}
}

func (e *Environment) Load() {
	root, err := utils.GetRootDir()
	if err != nil {
		panic(err)
	}
	dir, err := os.ReadDir(root)
	if err != nil {
		panic(err)
	}
	for _, f := range dir {
		if f.IsDir() {
			continue
		}
		if f.Name() == e.EnvFileName {
			e.EnvFound = true
		}
	}

}

func (e *Environment) Configure() {
	r := utils.Reader()
	teamId, err := utils.UserInput(r, "Enter your Vercel team ID")
	if err != nil {
		log.Fatal(err)
	}
	authKey, err := utils.UserInput(r, "Enter your Vercel auth token")
	if err != nil {
		log.Fatal(err)
	}
	e.TeamID = teamId
	e.AuthKey = authKey

	e.writeEnvFile(e.EnvFileName)
	fmt.Println("Configuration successful. Rerun the start command")
	os.Exit(0)
}

func (e *Environment) writeEnvFile(filename string) {
	temp := map[string]string{
		"VERCEL_ENDPOINT": e.ApiEndpoint,
		"VERCEL_AUTH_KEY": e.AuthKey,
		"VERCEL_TEAM_ID":  e.TeamID,
	}

	root, err := utils.GetRootDir()
	if err != nil {
		panic(err)
	}

	fp := root + "/" + filename
	fmt.Println("Creating file at " + fp)
	c := []byte{}
	for k, v := range temp {
		s := fmt.Sprintf("%s=%s\n", k, v)
		c = append(c, []byte(s)...)
	}

	err = os.WriteFile(fp, c, 0644)
	if err != nil {
		panic(err)
	}
}
