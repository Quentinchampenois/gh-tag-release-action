package tag_test

import (
	"ghtagreleaseaction/internal/tag"
	"github.com/google/go-github/v62/github"
	"testing"
)

func TestNewTag(t *testing.T) {
	label := github.Label{Name: github.String("release: v0.0.1")}

	tag := tag.NewTag(label)
	if tag.Version != "0.0.1" {
		t.Errorf("Expected version to be 0.0.1, got %s", tag.Version)
	}
}

func TestNewTag_OtherLabel(t *testing.T) {
	label := github.Label{Name: github.String("good first issue")}

	tag := tag.NewTag(label)
	if tag.Version != "" {
		t.Errorf("Expected version to be empty, got %s", tag.Version)
	}
}

func TestTag_IsReleaseLabel(t *testing.T) {
	label := github.Label{Name: github.String("release: v0.0.1")}

	tag := tag.NewTag(label)
	if !tag.IsReleaseLabel() {
		t.Errorf("Expected IsReleaseLabel to be true, got false")
	}
}

func TestTag_IsReleaseLabel_OtherLabel(t *testing.T) {
	label := github.Label{Name: github.String("good first issue")}

	tag := tag.NewTag(label)
	if tag.IsReleaseLabel() {
		t.Errorf("Expected IsReleaseLabel to be false, got true")
	}
}

func TestTag_ExtractVersion(t *testing.T) {
	var tests = []struct {
		label string
		want  string
	}{
		{"0.0.1", "0.0.1"},
		{"v0.0.1", "0.0.1"},
		{"release:v0.0.1", "0.0.1"},
		{"release: v0.0.1", "0.0.1"},
		{"release:0.0.1", "0.0.1"},
		{"release: 0.0.1", "0.0.1"},
		{"release: 0.1", "0.1.0"},
		{"release:0.1", "0.1.0"},
		{"release:v0.1", "0.1.0"},
		{"release: v0.1", "0.1.0"},
	}

	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			label := github.Label{Name: github.String(tt.label)}
			tag := tag.NewTag(label)
			if tag.ExtractVersion() != tt.want {
				want := tt.want
				if len(tt.want) < 1 {
					want = "empty"
				}

				t.Errorf("got %s, want %s", tt.label, want)
			}
		})
	}
}

func TestTag_Version(t *testing.T) {
	var tests = []struct {
		label string
		want  string
	}{
		{"0.0.1", ""},
		{"v0.0.1", ""},
		{"release:v0.0.1", "0.0.1"},
		{"release: v0.0.1", "0.0.1"},
		{"release:0.0.1", "0.0.1"},
		{"release: 0.0.1", "0.0.1"},
		{"release: 0.1", "0.1.0"},
		{"release:0.1", "0.1.0"},
		{"release:v0.1", "0.1.0"},
		{"release: v0.1", "0.1.0"},
	}

	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			label := github.Label{Name: github.String(tt.label)}
			tag := tag.NewTag(label)
			if tag.Version != tt.want {
				want := tt.want
				if len(tt.want) < 1 {
					want = "empty"
				}

				t.Errorf("got %s, want %s", tt.label, want)
			}
		})
	}
}
