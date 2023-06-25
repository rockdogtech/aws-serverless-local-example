package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/beatlabs/github-auth/app/inst"
	"github.com/beatlabs/github-auth/key"
	"github.com/google/go-github/v53/github"
	"os"
)

var (
	// DefaultHTTPGetAddress Default Address
	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"

	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("no IP in HTTP response")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("non 200 Response found")

	ErrNoGithubPemConfigured = errors.New("no Github Pem Configured")

	ErrNoGithubApplicationIdConfigured = errors.New("no Github Application Id Configured")

	ErrNoGithubInstallationIdConfigured = errors.New("no Github Installation Id Configured")

	ErrUnableToParsePem = errors.New("unable to parse pem")
)

type GithubProxyResponse struct {
	StatusCode        int                 `json:"statusCode"`
	Headers           map[string]string   `json:"headers"`
	MultiValueHeaders map[string][]string `json:"multiValueHeaders"`
	Body              string              `json:"body"`
	IsBase64Encoded   bool                `json:"isBase64Encoded,omitempty"`
}

func handler(request events.APIGatewayProxyRequest) (GithubProxyResponse, error) {
	//resp, err := http.Get(DefaultHTTPGetAddress)
	//if err != nil {
	//	return events.APIGatewayProxyResponse{}, err
	//}
	apiKey := os.Getenv("GITHUB_AUTH_PEM")
	if apiKey == "" {
		return GithubProxyResponse{}, ErrNoGithubPemConfigured
	}

	installIDEnvValue := os.Getenv("GITHUB_AUTH_INST_ID")
	if installIDEnvValue == "" {
		return GithubProxyResponse{}, ErrNoGithubInstallationIdConfigured
	}

	appIDEnvValue := os.Getenv("GITHUB_AUTH_APP_ID")
	if appIDEnvValue == "" {
		return GithubProxyResponse{}, ErrNoGithubApplicationIdConfigured
	}

	pemBytes, err := key.Parse([]byte(apiKey))
	if err != nil {
		return GithubProxyResponse{}, ErrUnableToParsePem
	}
	githubAppInstallation, err := inst.NewConfig(appIDEnvValue, installIDEnvValue, pemBytes)
	ctx := context.Background()

	client := github.NewClient(githubAppInstallation.Client(ctx))
	opt := &github.RepositoryListByOrgOptions{Type: "public"}
	repos, resp, err := client.Repositories.ListByOrg(context.Background(), "github", opt)

	if resp.StatusCode != 200 {
		return GithubProxyResponse{}, ErrNon200Response
	}

	for repo := range repos {
		fmt.Printf("repo %+v", repo)
	}

	return GithubProxyResponse{
		StatusCode: 200,
	}, nil

	//ip, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	return events.APIGatewayProxyResponse{}, err
	//}
	//
	//if len(ip) == 0 {
	//	return events.APIGatewayProxyResponse{}, ErrNoIP
	//}
}

func main() {
	lambda.Start(handler)
}
