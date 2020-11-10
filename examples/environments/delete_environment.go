package examples

import (
	"fmt"
	"net/url"

	"github.com/fqjony/go-octopusdeploy/octopusdeploy"
)

func DeleteEnvironmentExample() {
	var (
		apiKey     string = "API-YOUR_API_KEY"
		octopusURL string = "https://your_octopus_url"
		spaceID    string = "space-id"

		// environment values
		environmentID string = "environment-id"
	)

	apiURL, err := url.Parse(octopusURL)
	if err != nil {
		_ = fmt.Errorf("error parsing URL for Octopus API: %v", err)
		return
	}

	client, err := octopusdeploy.NewClient(nil, apiURL, apiKey, spaceID)
	if err != nil {
		_ = fmt.Errorf("error creating API client: %v", err)
		return
	}

	// delete environment
	err = client.Environments.DeleteByID(environmentID)
	if err != nil {
		_ = fmt.Errorf("error deleting environment: %v", err)
		return
	}

	fmt.Printf("environment deleted: (%s)\n", environmentID)
}
