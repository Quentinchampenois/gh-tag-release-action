package main

import (
	"context"
	"fmt"
	"ghtagreleaseaction/internal/tag"
	"github.com/google/go-github/v62/github"
	"os"
	"strings"
)

const inputCrashOnError = "INPUT_CRASH_ON_ERROR"

func main() {
	ctx := context.Background()
	client := github.NewClient(nil)
	login, repoName := ParseRepositoryInformations()
	owner := github.User{Login: &login}

	pr, _, err := client.PullRequests.ListPullRequestsWithCommit(ctx, login, repoName, os.Getenv("GITHUB_SHA"), nil)
	if err != nil {
		fmt.Printf("Error fetching pull requests for %s/%s: %v", login, repoName, err)

		if os.Getenv(inputCrashOnError) == "true" {
			os.Exit(2)
		} else {
			return
		}
	}

	if len(pr) == 0 {
		fmt.Println("No pull request found")
		if os.Getenv(inputCrashOnError) == "true" {
			os.Exit(2)
		} else {
			return
		}
	}

	prNumber := pr[0].Number

	if prNumber == nil || *prNumber == 0 {
		fmt.Println("No PR found")
		if os.Getenv(inputCrashOnError) == "true" {
			os.Exit(2)
		} else {
			return
		}
	}

	var outputs []string

	releaseTag := tag.GetReleaseTag(ctx, client, &github.Repository{Owner: &owner, Name: &repoName}, *prNumber)
	if releaseTag.Version == "" {
		fmt.Println("No release tag found")
		if os.Getenv(inputCrashOnError) == "true" {
			os.Exit(2)
		} else {
			return
		}
	}

	outputs = append(outputs, fmt.Sprintf("tag=%s", releaseTag.Version))
	output := strings.Join(outputs, "\n")
	fmt.Println(output)
	printToStdout(output)
}

func ParseRepositoryInformations() (string, string) {
	repository := os.Getenv("GITHUB_REPOSITORY")
	if repository == "" {
		fmt.Println("GITHUB_REPOSITORY is not set")
		return "", ""
	}

	ownerRepo := strings.Split(repository, "/")
	return ownerRepo[0], ownerRepo[1]
}

func printToStdout(output string) {
	githubOutput := os.Getenv("GITHUB_OUTPUT")
	if githubOutput == "" {
		fmt.Println("GITHUB_OUTPUT is not set")
		return
	}

	file, err := os.OpenFile(githubOutput, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Printf("Error opening GITHUB_OUTPUT file: %v\n", err)
		return
	}
	defer file.Close()

	if _, err = file.WriteString(output); err != nil {
		fmt.Printf("Error writing to GITHUB_OUTPUT file: %v\n", err)
		return
	}
}
