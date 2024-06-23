package tag

import (
	"context"
	"fmt"
	"github.com/google/go-github/v62/github"
	"log"
	"strings"
)

type Tag struct {
	Version string       `json:"tag,omitempty"`
	Label   github.Label `json:"label,omitempty"`
}

func NewTag(label github.Label) Tag {
	tag := Tag{
		Version: "",
		Label:   label,
	}

	if tag.IsReleaseLabel() {
		tag.Version = tag.ExtractVersion()
	}

	return tag
}

func (t Tag) String() string {
	return t.Version
}

func (t Tag) IsReleaseLabel() bool {
	return strings.Contains(*t.Label.Name, "release:")
}

func (t Tag) ExtractVersion() string {
	version := strings.Replace(*t.Label.Name, "release:", "", 1)
	version = strings.TrimSpace(version)
	if version[0] == 'v' || version[0] == 'V' {
		version = version[1:]
	}
	if len(version) < 5 {
		version = fmt.Sprintf("%s.0", version)
	}

	return version
}

func GetReleaseTag(ctx context.Context, client *github.Client, repo *github.Repository, number int) Tag {
	var tag Tag
	pr, _, err := client.PullRequests.Get(ctx, *repo.Owner.Login, *repo.Name, number)
	if err != nil {
		log.Printf("Error fetching pull requests for %s/%s: %v", "opensourcepolitics", *repo.Name, err)
		return tag
	}

	for _, label := range pr.Labels {
		t := NewTag(*label)
		if t.IsReleaseLabel() {
			tag = t
			break
		}
	}

	return tag
}
