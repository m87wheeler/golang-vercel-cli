package environment

import (
	"fmt"
	"log"
	"os"

	"github.com/m87wheeler/golang-vercel-cli/pkg/utils"
)

const (
	ENV_FILE_DIR  = "go_vercel_cli/"
	ENV_FILE_NAME = ".env"
)

type Environment struct {
	// env file
	EnvFileDir  string // path segment for where the .env file will be
	EnvFileName string // the name of the .env file
	EnvLoadFrom string // a complete filepath for the .env file for use with github.com/joho/godotenv
	EnvFound    bool
	// vercel values
	TeamID      string
	AuthKey     string
	ApiEndpoint string
	Projects    map[string]string
}

func NewEnvironment() *Environment {
	return &Environment{
		EnvFileDir:  ENV_FILE_DIR,
		EnvFileName: ENV_FILE_NAME,
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
	homeDir, err := utils.GetHomeDir()
	if err != nil {
		panic(err)
	}

	fp := homeDir + "/" + e.EnvFileDir

	dir, err := os.ReadDir(fp)
	if err != nil {
		if os.IsNotExist(err) {
			return
		} else {
			panic(err)
		}
	}

	for _, f := range dir {
		if f.IsDir() {
			continue
		}
		if f.Name() == e.EnvFileName {
			e.EnvLoadFrom = fp + e.EnvFileName
			e.EnvFound = true
		}
	}

}

func (e *Environment) Configure() {
	r := utils.Reader()
	teamId, err := utils.UserInput(r, "Enter your Vercel team ID", false)
	if err != nil {
		log.Fatal(err)
	}
	authKey, err := utils.UserInput(r, "Enter your Vercel auth token", true)
	if err != nil {
		log.Fatal(err)
	}
	e.TeamID = teamId
	e.AuthKey = authKey

	_, err = e.writeEnvFile()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Configuration successful")
	fmt.Println("Rerun the start command to continue")
	os.Exit(0)
}

func (e *Environment) writeEnvFile() (string, error) {
	temp := map[string]string{
		"VERCEL_ENDPOINT": e.ApiEndpoint,
		"VERCEL_AUTH_KEY": e.AuthKey,
		"VERCEL_TEAM_ID":  e.TeamID,
	}

	homeDir, err := utils.GetHomeDir()
	if err != nil {
		return "", err
	}

	// Create the necessary directories, if they don't exist
	dp := homeDir + "/" + e.EnvFileDir
	err = os.MkdirAll(dp, 0755)
	if err != nil {
		return "", err
	}

	fp := dp + e.EnvFileName
	fmt.Println("\nCreating environment file at " + fp)
	c := []byte{}
	for k, v := range temp {
		s := fmt.Sprintf("%s=%s\n", k, v)
		c = append(c, []byte(s)...)
	}

	err = os.WriteFile(fp, c, 0644)
	if err != nil {
		return "", err
	}

	return fp, nil
}
