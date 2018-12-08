package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/antihax/optional"

	vstsObj "github.com/samkreter/vsts-goclient/api/git"
	vsts "github.com/samkreter/vsts-goclient/client"
)

// Config holds the configuration from the config file
type Config struct {
	Token          string `json:"token"`
	Username       string `json:"username"`
	APIVersion     string `json:"apiVersion"`
	RepositoryName string `json:"repositoryName"`
	Project        string `json:"project"`
	Instance       string `json:"instance"`
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if len(os.Args) < 2 {
		log.Fatal("USAGE: vstsAlias <alias>")
	}

	alias := os.Args[1]

	configFilePath, ok := os.LookupEnv("CONFIG_PATH")
	if !ok {
		log.Fatal("CONFIG_PATH not set")
	}

	configFile, err := os.Open(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	defer configFile.Close()

	var config Config
	if err := json.NewDecoder(configFile).Decode(&config); err != nil {
		log.Fatal(err)
	}

	vstsConfig := &vsts.Config{
		Token:          config.Token,
		Username:       config.Username,
		APIVersion:     config.APIVersion,
		RepositoryName: config.RepositoryName,
		Project:        config.Project,
		Instance:       config.Instance,
	}

	vstsClient, err := vsts.NewClient(vstsConfig)
	if err != nil {
		log.Fatal(err)
	}

	getOpts := &vstsObj.GetPullRequestsOpts{
		SearchCriteriaStatus: optional.NewString("all"),
	}

	pullRequests, err := vstsClient.GetPullRequests(getOpts)
	if err != nil {
		log.Fatalf("get pull requests error: %v", err)
	}

	for _, pullRequest := range pullRequests {
		if getAliasFromEmail(pullRequest.CreatedBy.UniqueName) == alias {
			fmt.Printf("User %s has VSTS Id of %s", alias, pullRequests[0].CreatedBy.ID)
			return
		}
	}

	log.Fatalf("User %s has not made any pull requests, Could not retrieve their ID", alias)
}

func getAliasFromEmail(email string) string {
	splitEmail := strings.Split(email, "@")

	if len(splitEmail) == 2 {
		return splitEmail[0]
	}
	return ""
}
